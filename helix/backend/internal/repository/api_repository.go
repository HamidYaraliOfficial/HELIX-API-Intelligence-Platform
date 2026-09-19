package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/helix-platform/helix-backend/internal/models"
	"github.com/helix-platform/helix-backend/internal/openapi"
)

type APIRepository struct {
	DB  *sql.DB
	Tag *TagRepository
}

func NewAPIRepository(db *sql.DB) *APIRepository {
	return &APIRepository{DB: db, Tag: NewTagRepository(db)}
}

// ListFilter carries every optional filter the Catalog UI can apply.
type ListFilter struct {
	OrgID          string
	Query          string
	Protocol       string
	Status         string
	Environment    string
	Visibility     string
	LifecycleState string
	TeamID         string
	Tag            string
	Page           int
	PageSize       int
}

func (r *APIRepository) List(f ListFilter) (models.PageResult[models.API], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 200 {
		f.PageSize = 25
	}

	where := []string{"a.org_id = $1"}
	args := []interface{}{f.OrgID}
	arg := func(v interface{}) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	if f.Query != "" {
		where = append(where, fmt.Sprintf(
			"(a.name ILIKE %s OR a.description ILIKE %s OR a.owner ILIKE %s)",
			arg("%"+f.Query+"%"), arg("%"+f.Query+"%"), arg("%"+f.Query+"%"),
		))
	}
	if f.Protocol != "" {
		where = append(where, "a.protocol = "+arg(f.Protocol))
	}
	if f.Status != "" {
		where = append(where, "a.status = "+arg(f.Status))
	}
	if f.Environment != "" {
		where = append(where, "a.environment = "+arg(f.Environment))
	}
	if f.Visibility != "" {
		where = append(where, "a.visibility = "+arg(f.Visibility))
	}
	if f.LifecycleState != "" {
		where = append(where, "a.lifecycle_state = "+arg(f.LifecycleState))
	}
	if f.TeamID != "" {
		where = append(where, "a.team_id = "+arg(f.TeamID))
	}
	if f.Tag != "" {
		where = append(where, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM api_tags at2 JOIN tags tg ON tg.id = at2.tag_id WHERE at2.api_id = a.id AND tg.name = %s)",
			arg(f.Tag),
		))
	}

	whereSQL := strings.Join(where, " AND ")

	var total int
	countQuery := "SELECT COUNT(*) FROM apis a WHERE " + whereSQL
	if err := r.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return models.PageResult[models.API]{}, fmt.Errorf("count apis: %w", err)
	}

	offset := (f.Page - 1) * f.PageSize
	limitArg := arg(f.PageSize)
	offsetArg := arg(offset)

	query := fmt.Sprintf(`
		SELECT a.id, a.org_id, a.team_id, t.name AS team_name, a.name, a.description, a.owner,
		       a.protocol, a.version, a.base_url, a.environment, a.visibility, a.status,
		       a.lifecycle_state, a.source_type, a.auth_type,
		       a.operating_hours_start, a.operating_hours_end, a.operating_days, a.timezone,
		       a.created_at, a.updated_at,
		       COALESCE((SELECT COUNT(*) FROM endpoints e WHERE e.api_id = a.id), 0) AS endpoint_count,
		       COALESCE(ARRAY_AGG(tg.name) FILTER (WHERE tg.name IS NOT NULL), ARRAY[]::text[]) AS tags
		FROM apis a
		LEFT JOIN teams t ON t.id = a.team_id
		LEFT JOIN api_tags at2 ON at2.api_id = a.id
		LEFT JOIN tags tg ON tg.id = at2.tag_id
		WHERE %s
		GROUP BY a.id, t.name
		ORDER BY a.updated_at DESC
		LIMIT %s OFFSET %s
	`, whereSQL, limitArg, offsetArg)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return models.PageResult[models.API]{}, fmt.Errorf("list apis: %w", err)
	}
	defer rows.Close()

	var items []models.API
	for rows.Next() {
		a, err := scanAPI(rows)
		if err != nil {
			return models.PageResult[models.API]{}, err
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return models.PageResult[models.API]{}, err
	}
	if items == nil {
		items = []models.API{}
	}

	totalPages := (total + f.PageSize - 1) / f.PageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return models.PageResult[models.API]{
		Items:      items,
		Total:      total,
		Page:       f.Page,
		PageSize:   f.PageSize,
		TotalPages: totalPages,
	}, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanAPI(row rowScanner) (models.API, error) {
	var a models.API
	var tags []string
	err := row.Scan(
		&a.ID, &a.OrgID, &a.TeamID, &a.TeamName, &a.Name, &a.Description, &a.Owner,
		&a.Protocol, &a.Version, &a.BaseURL, &a.Environment, &a.Visibility, &a.Status,
		&a.LifecycleState, &a.SourceType, &a.AuthType,
		&a.OperatingStart, &a.OperatingEnd, &a.OperatingDays, &a.Timezone,
		&a.CreatedAt, &a.UpdatedAt, &a.EndpointCount, pq.Array(&tags),
	)
	a.Tags = tags
	return a, err
}

func (r *APIRepository) GetByID(orgID, id string) (*models.API, error) {
	query := `
		SELECT a.id, a.org_id, a.team_id, t.name AS team_name, a.name, a.description, a.owner,
		       a.protocol, a.version, a.base_url, a.environment, a.visibility, a.status,
		       a.lifecycle_state, a.source_type, a.auth_type,
		       a.operating_hours_start, a.operating_hours_end, a.operating_days, a.timezone,
		       a.created_at, a.updated_at,
		       COALESCE((SELECT COUNT(*) FROM endpoints e WHERE e.api_id = a.id), 0) AS endpoint_count,
		       COALESCE(ARRAY_AGG(tg.name) FILTER (WHERE tg.name IS NOT NULL), ARRAY[]::text[]) AS tags
		FROM apis a
		LEFT JOIN teams t ON t.id = a.team_id
		LEFT JOIN api_tags at2 ON at2.api_id = a.id
		LEFT JOIN tags tg ON tg.id = at2.tag_id
		WHERE a.org_id = $1 AND a.id = $2
		GROUP BY a.id, t.name
	`
	a, err := scanAPI(r.DB.QueryRow(query, orgID, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get api: %w", err)
	}
	return &a, nil
}

func (r *APIRepository) GetEndpoints(apiID string) ([]models.Endpoint, error) {
	rows, err := r.DB.Query(`
		SELECT id, api_id, path, method, operation_id, summary, description, auth_type, deprecated
		FROM endpoints WHERE api_id = $1 ORDER BY path ASC, method ASC
	`, apiID)
	if err != nil {
		return nil, fmt.Errorf("list endpoints: %w", err)
	}
	defer rows.Close()

	var endpoints []models.Endpoint
	for rows.Next() {
		var e models.Endpoint
		if err := rows.Scan(&e.ID, &e.APIID, &e.Path, &e.Method, &e.OperationID, &e.Summary, &e.Description, &e.AuthType, &e.Deprecated); err != nil {
			return nil, fmt.Errorf("scan endpoint: %w", err)
		}
		endpoints = append(endpoints, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range endpoints {
		params, err := r.getParameters(endpoints[i].ID)
		if err != nil {
			return nil, err
		}
		endpoints[i].Parameters = params

		responses, err := r.getResponses(endpoints[i].ID)
		if err != nil {
			return nil, err
		}
		endpoints[i].Responses = responses
	}

	if endpoints == nil {
		endpoints = []models.Endpoint{}
	}
	return endpoints, nil
}

func (r *APIRepository) getParameters(endpointID string) ([]models.EndpointParameter, error) {
	rows, err := r.DB.Query(`
		SELECT id, endpoint_id, name, location, type, required, description
		FROM endpoint_parameters WHERE endpoint_id = $1 ORDER BY name ASC
	`, endpointID)
	if err != nil {
		return nil, fmt.Errorf("list parameters: %w", err)
	}
	defer rows.Close()

	var params []models.EndpointParameter
	for rows.Next() {
		var p models.EndpointParameter
		if err := rows.Scan(&p.ID, &p.EndpointID, &p.Name, &p.In, &p.Type, &p.Required, &p.Description); err != nil {
			return nil, err
		}
		params = append(params, p)
	}
	return params, rows.Err()
}

func (r *APIRepository) getResponses(endpointID string) ([]models.EndpointResponse, error) {
	rows, err := r.DB.Query(`
		SELECT id, endpoint_id, status_code, description, content_type
		FROM endpoint_responses WHERE endpoint_id = $1 ORDER BY status_code ASC
	`, endpointID)
	if err != nil {
		return nil, fmt.Errorf("list responses: %w", err)
	}
	defer rows.Close()

	var responses []models.EndpointResponse
	for rows.Next() {
		var resp models.EndpointResponse
		if err := rows.Scan(&resp.ID, &resp.EndpointID, &resp.StatusCode, &resp.Description, &resp.ContentType); err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, rows.Err()
}

// CreateInput carries the fields accepted from manual registration or discovery.
type CreateInput struct {
	OrgID          string
	TeamID         *string
	Name           string
	Description    string
	Owner          string
	Protocol       string
	Version        string
	BaseURL        string
	Environment    string
	Visibility     string
	Status         string
	LifecycleState string
	SourceType     string
	AuthType       string
	OperatingStart string
	OperatingEnd   string
	OperatingDays  string
	Timezone       string
	Tags           []string
	RawSpec        []byte // optional, stored for audit / future diffing
}

func (r *APIRepository) Create(in CreateInput) (*models.API, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var id string
	err = tx.QueryRow(`
		INSERT INTO apis (
			org_id, team_id, name, description, owner, protocol, version, base_url,
			environment, visibility, status, lifecycle_state, source_type, auth_type,
			operating_hours_start, operating_hours_end, operating_days, timezone, raw_spec
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18, NULLIF($19, '')::jsonb)
		RETURNING id
	`, in.OrgID, in.TeamID, in.Name, in.Description, in.Owner, in.Protocol, in.Version, in.BaseURL,
		in.Environment, in.Visibility, in.Status, in.LifecycleState, in.SourceType, in.AuthType,
		in.OperatingStart, in.OperatingEnd, in.OperatingDays, in.Timezone, safeJSON(in.RawSpec),
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert api: %w", err)
	}

	if len(in.Tags) > 0 {
		if err := r.Tag.EnsureAndAttach(tx, in.OrgID, id, in.Tags); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(in.OrgID, id)
}

func (r *APIRepository) Update(orgID, id string, in CreateInput) (*models.API, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		UPDATE apis SET
			team_id = $1, name = $2, description = $3, owner = $4, protocol = $5, version = $6,
			base_url = $7, environment = $8, visibility = $9, status = $10, lifecycle_state = $11,
			auth_type = $12, operating_hours_start = $13, operating_hours_end = $14,
			operating_days = $15, timezone = $16, updated_at = now()
		WHERE org_id = $17 AND id = $18
	`, in.TeamID, in.Name, in.Description, in.Owner, in.Protocol, in.Version,
		in.BaseURL, in.Environment, in.Visibility, in.Status, in.LifecycleState,
		in.AuthType, in.OperatingStart, in.OperatingEnd, in.OperatingDays, in.Timezone,
		orgID, id)
	if err != nil {
		return nil, fmt.Errorf("update api: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, nil
	}

	if err := r.Tag.DetachAll(tx, id); err != nil {
		return nil, err
	}
	if len(in.Tags) > 0 {
		if err := r.Tag.EnsureAndAttach(tx, orgID, id, in.Tags); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(orgID, id)
}

func (r *APIRepository) Delete(orgID, id string) (bool, error) {
	res, err := r.DB.Exec(`DELETE FROM apis WHERE org_id = $1 AND id = $2`, orgID, id)
	if err != nil {
		return false, fmt.Errorf("delete api: %w", err)
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// ReplaceEndpointsFromContract stores the endpoints discovered from a parsed
// OpenAPI contract, replacing any endpoints previously imported for this API.
func (r *APIRepository) ReplaceEndpointsFromContract(apiID string, contract *openapi.Contract) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM endpoints WHERE api_id = $1`, apiID); err != nil {
		return 0, fmt.Errorf("clear endpoints: %w", err)
	}

	for _, ep := range contract.Endpoints {
		var endpointID string
		err := tx.QueryRow(`
			INSERT INTO endpoints (api_id, path, method, operation_id, summary, description, auth_type, deprecated)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			RETURNING id
		`, apiID, ep.Path, ep.Method, ep.OperationID, ep.Summary, ep.Description, ep.AuthType, ep.Deprecated).Scan(&endpointID)
		if err != nil {
			return 0, fmt.Errorf("insert endpoint %s %s: %w", ep.Method, ep.Path, err)
		}

		for _, p := range ep.Parameters {
			if _, err := tx.Exec(`
				INSERT INTO endpoint_parameters (endpoint_id, name, location, type, required, description)
				VALUES ($1,$2,$3,$4,$5,$6)
			`, endpointID, p.Name, p.In, p.Type, p.Required, p.Description); err != nil {
				return 0, fmt.Errorf("insert parameter %s: %w", p.Name, err)
			}
		}
		for _, resp := range ep.Responses {
			if _, err := tx.Exec(`
				INSERT INTO endpoint_responses (endpoint_id, status_code, description, content_type)
				VALUES ($1,$2,$3,$4)
			`, endpointID, resp.StatusCode, resp.Description, resp.ContentType); err != nil {
				return 0, fmt.Errorf("insert response %s: %w", resp.StatusCode, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(contract.Endpoints), nil
}

func safeJSON(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	return string(raw)
}
