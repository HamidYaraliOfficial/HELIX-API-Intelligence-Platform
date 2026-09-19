package models

import "time"

// Organization is the top-level multi-tenant boundary in HELIX.
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// Team owns one or more APIs inside an organization.
type Team struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	APICount  int       `json:"api_count"`
	CreatedAt time.Time `json:"created_at"`
}

// API is a single cataloged API (REST, GraphQL, gRPC, WebSocket, SOAP, ...).
type API struct {
	ID              string    `json:"id"`
	OrgID           string    `json:"org_id"`
	TeamID          *string   `json:"team_id,omitempty"`
	TeamName        *string   `json:"team_name,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Owner           string    `json:"owner"`
	Protocol        string    `json:"protocol"`
	Version         string    `json:"version"`
	BaseURL         string    `json:"base_url"`
	Environment     string    `json:"environment"`
	Visibility      string    `json:"visibility"`
	Status          string    `json:"status"`
	LifecycleState  string    `json:"lifecycle_state"`
	SourceType      string    `json:"source_type"`
	AuthType        string    `json:"auth_type"`
	OperatingStart  string    `json:"operating_hours_start"`
	OperatingEnd    string    `json:"operating_hours_end"`
	OperatingDays   string    `json:"operating_days"`
	Timezone        string    `json:"timezone"`
	EndpointCount   int       `json:"endpoint_count"`
	Tags            []string  `json:"tags"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Endpoint is a single method+path operation inside an API.
type Endpoint struct {
	ID          string               `json:"id"`
	APIID       string               `json:"api_id"`
	Path        string               `json:"path"`
	Method      string               `json:"method"`
	OperationID string               `json:"operation_id"`
	Summary     string               `json:"summary"`
	Description string               `json:"description"`
	AuthType    string               `json:"auth_type"`
	Deprecated  bool                 `json:"deprecated"`
	Parameters  []EndpointParameter  `json:"parameters"`
	Responses   []EndpointResponse   `json:"responses"`
}

// EndpointParameter describes one request parameter of an endpoint.
type EndpointParameter struct {
	ID          string `json:"id"`
	EndpointID  string `json:"endpoint_id"`
	Name        string `json:"name"`
	In          string `json:"in"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// EndpointResponse describes one possible response of an endpoint.
type EndpointResponse struct {
	ID          string `json:"id"`
	EndpointID  string `json:"endpoint_id"`
	StatusCode  string `json:"status_code"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
}

// Tag is a free-form label attached to APIs for filtering and grouping.
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DiscoveryJob tracks the result of one OpenAPI/GraphQL discovery import.
type DiscoveryJob struct {
	ID             string     `json:"id"`
	OrgID          string     `json:"org_id"`
	Filename       string     `json:"filename"`
	SourceType     string     `json:"source_type"`
	Status         string     `json:"status"`
	APIsFound      int        `json:"apis_found"`
	EndpointsFound int        `json:"endpoints_found"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// Stats aggregates catalog-wide numbers for the dashboard.
type Stats struct {
	TotalAPIs        int            `json:"total_apis"`
	TotalEndpoints   int            `json:"total_endpoints"`
	TotalTeams       int            `json:"total_teams"`
	ByProtocol       map[string]int `json:"by_protocol"`
	ByStatus         map[string]int `json:"by_status"`
	ByEnvironment    map[string]int `json:"by_environment"`
	ByLifecycleState map[string]int `json:"by_lifecycle_state"`
}

// PageResult wraps any paginated list response.
type PageResult[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// Allowed enum values, used for both backend validation and to drive the
// frontend's filter dropdowns via the /api/v1/meta endpoint.
var (
	Protocols       = []string{"REST", "GraphQL", "gRPC", "WebSocket", "SOAP"}
	Environments    = []string{"development", "test", "staging", "production"}
	Visibilities    = []string{"public", "internal", "partner", "private"}
	LifecycleStates = []string{"draft", "experimental", "active", "deprecated", "sunset", "retired"}
	Statuses        = []string{"healthy", "degraded", "down", "unknown"}
	AuthTypes       = []string{"none", "apiKey", "basic", "bearer", "oauth2", "oidc", "mtls"}
)
