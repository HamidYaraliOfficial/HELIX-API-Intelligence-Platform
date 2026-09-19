package repository

import (
	"database/sql"
	"fmt"

	"github.com/helix-platform/helix-backend/internal/models"
)

type DiscoveryRepository struct {
	DB *sql.DB
}

func NewDiscoveryRepository(db *sql.DB) *DiscoveryRepository {
	return &DiscoveryRepository{DB: db}
}

func (r *DiscoveryRepository) Create(orgID, filename, sourceType string) (*models.DiscoveryJob, error) {
	var j models.DiscoveryJob
	err := r.DB.QueryRow(`
		INSERT INTO discovery_jobs (org_id, filename, source_type, status)
		VALUES ($1, $2, $3, 'running')
		RETURNING id, org_id, filename, source_type, status, apis_found, endpoints_found, created_at
	`, orgID, filename, sourceType).Scan(
		&j.ID, &j.OrgID, &j.Filename, &j.SourceType, &j.Status, &j.APIsFound, &j.EndpointsFound, &j.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create discovery job: %w", err)
	}
	return &j, nil
}

func (r *DiscoveryRepository) Complete(id string, apisFound, endpointsFound int, errMsg string) error {
	status := "completed"
	if errMsg != "" {
		status = "failed"
	}
	_, err := r.DB.Exec(`
		UPDATE discovery_jobs
		SET status = $1, apis_found = $2, endpoints_found = $3, error_message = $4, completed_at = now()
		WHERE id = $5
	`, status, apisFound, endpointsFound, errMsg, id)
	return err
}

func (r *DiscoveryRepository) List(orgID string, limit int) ([]models.DiscoveryJob, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.DB.Query(`
		SELECT id, org_id, filename, source_type, status, apis_found, endpoints_found,
		       COALESCE(error_message, ''), created_at, completed_at
		FROM discovery_jobs
		WHERE org_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, orgID, limit)
	if err != nil {
		return nil, fmt.Errorf("list discovery jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.DiscoveryJob
	for rows.Next() {
		var j models.DiscoveryJob
		var completedAt sql.NullTime
		if err := rows.Scan(&j.ID, &j.OrgID, &j.Filename, &j.SourceType, &j.Status,
			&j.APIsFound, &j.EndpointsFound, &j.ErrorMessage, &j.CreatedAt, &completedAt); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			j.CompletedAt = &completedAt.Time
		}
		jobs = append(jobs, j)
	}
	if jobs == nil {
		jobs = []models.DiscoveryJob{}
	}
	return jobs, rows.Err()
}
