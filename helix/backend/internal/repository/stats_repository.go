package repository

import (
	"database/sql"
	"fmt"

	"github.com/helix-platform/helix-backend/internal/models"
)

type StatsRepository struct {
	DB *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{DB: db}
}

func (r *StatsRepository) Get(orgID string) (*models.Stats, error) {
	s := &models.Stats{
		ByProtocol:       map[string]int{},
		ByStatus:         map[string]int{},
		ByEnvironment:    map[string]int{},
		ByLifecycleState: map[string]int{},
	}

	if err := r.DB.QueryRow(`SELECT COUNT(*) FROM apis WHERE org_id = $1`, orgID).Scan(&s.TotalAPIs); err != nil {
		return nil, fmt.Errorf("count apis: %w", err)
	}

	if err := r.DB.QueryRow(`
		SELECT COUNT(*) FROM endpoints e JOIN apis a ON a.id = e.api_id WHERE a.org_id = $1
	`, orgID).Scan(&s.TotalEndpoints); err != nil {
		return nil, fmt.Errorf("count endpoints: %w", err)
	}

	if err := r.DB.QueryRow(`SELECT COUNT(*) FROM teams WHERE org_id = $1`, orgID).Scan(&s.TotalTeams); err != nil {
		return nil, fmt.Errorf("count teams: %w", err)
	}

	if err := fillGroupCount(r.DB, orgID, "protocol", s.ByProtocol); err != nil {
		return nil, err
	}
	if err := fillGroupCount(r.DB, orgID, "status", s.ByStatus); err != nil {
		return nil, err
	}
	if err := fillGroupCount(r.DB, orgID, "environment", s.ByEnvironment); err != nil {
		return nil, err
	}
	if err := fillGroupCount(r.DB, orgID, "lifecycle_state", s.ByLifecycleState); err != nil {
		return nil, err
	}

	return s, nil
}

func fillGroupCount(db *sql.DB, orgID, column string, into map[string]int) error {
	// column comes only from hardcoded call sites above, never user input.
	query := fmt.Sprintf(`SELECT %s, COUNT(*) FROM apis WHERE org_id = $1 GROUP BY %s`, column, column)
	rows, err := db.Query(query, orgID)
	if err != nil {
		return fmt.Errorf("group by %s: %w", column, err)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return err
		}
		into[key] = count
	}
	return rows.Err()
}
