package repository

import (
	"database/sql"
	"fmt"

	"github.com/helix-platform/helix-backend/internal/models"
)

// TeamRepository provides data access for teams.
type TeamRepository struct {
	DB *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{DB: db}
}

func (r *TeamRepository) List(orgID string) ([]models.Team, error) {
	rows, err := r.DB.Query(`
		SELECT t.id, t.org_id, t.name, t.slug, t.created_at,
		       COUNT(a.id) AS api_count
		FROM teams t
		LEFT JOIN apis a ON a.team_id = t.id
		WHERE t.org_id = $1
		GROUP BY t.id
		ORDER BY t.name ASC
	`, orgID)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.Slug, &t.CreatedAt, &t.APICount); err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *TeamRepository) Create(orgID, name, slug string) (*models.Team, error) {
	var t models.Team
	err := r.DB.QueryRow(`
		INSERT INTO teams (org_id, name, slug)
		VALUES ($1, $2, $3)
		ON CONFLICT (org_id, slug) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, org_id, name, slug, created_at
	`, orgID, name, slug).Scan(&t.ID, &t.OrgID, &t.Name, &t.Slug, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}
	return &t, nil
}

// EnsureDefaultOrganization guarantees at least one organization exists,
// returning its ID. Used so a brand-new HELIX instance is immediately usable.
func EnsureDefaultOrganization(db *sql.DB, id, name, slug string) error {
	_, err := db.Exec(`
		INSERT INTO organizations (id, name, slug)
		VALUES ($1, $2, $3)
		ON CONFLICT (slug) DO NOTHING
	`, id, name, slug)
	return err
}
