package repository

import (
	"database/sql"
	"fmt"

	"github.com/helix-platform/helix-backend/internal/models"
)

type TagRepository struct {
	DB *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{DB: db}
}

func (r *TagRepository) List(orgID string) ([]models.Tag, error) {
	rows, err := r.DB.Query(`SELECT id, name FROM tags WHERE org_id = $1 ORDER BY name ASC`, orgID)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// EnsureAndAttach creates any tags that don't already exist for the
// organization and attaches all of them (existing + new) to the given API.
func (r *TagRepository) EnsureAndAttach(tx *sql.Tx, orgID, apiID string, names []string) error {
	for _, name := range names {
		if name == "" {
			continue
		}
		var tagID string
		err := tx.QueryRow(`
			INSERT INTO tags (org_id, name) VALUES ($1, $2)
			ON CONFLICT (org_id, name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, orgID, name).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("ensure tag %q: %w", name, err)
		}
		if _, err := tx.Exec(`
			INSERT INTO api_tags (api_id, tag_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, apiID, tagID); err != nil {
			return fmt.Errorf("attach tag %q: %w", name, err)
		}
	}
	return nil
}

func (r *TagRepository) DetachAll(tx *sql.Tx, apiID string) error {
	_, err := tx.Exec(`DELETE FROM api_tags WHERE api_id = $1`, apiID)
	return err
}
