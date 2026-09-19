package handlers

import (
	"database/sql"
	"net/http"

	"github.com/helix-platform/helix-backend/internal/config"
	"github.com/helix-platform/helix-backend/internal/middleware"
	"github.com/helix-platform/helix-backend/internal/repository"
)

// NewRouter builds the full HELIX HTTP API surface. Uses Go 1.22's
// standard-library ServeMux, which natively supports method-based routing
// and {param} path segments — no external router dependency required.
func NewRouter(db *sql.DB, cfg config.Config) http.Handler {
	apiRepo := repository.NewAPIRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	tagRepo := repository.NewTagRepository(db)
	statsRepo := repository.NewStatsRepository(db)
	discoveryRepo := repository.NewDiscoveryRepository(db)

	apiHandler := NewAPIHandler(apiRepo, cfg)
	teamHandler := NewTeamHandler(teamRepo, cfg)
	tagHandler := NewTagHandler(tagRepo, cfg)
	statsHandler := NewStatsHandler(statsRepo, cfg)
	discoveryHandler := NewDiscoveryHandler(apiRepo, discoveryRepo, cfg)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", Health)
	mux.HandleFunc("GET /api/v1/meta", Meta)

	mux.HandleFunc("GET /api/v1/apis", apiHandler.List)
	mux.HandleFunc("POST /api/v1/apis", apiHandler.Create)
	mux.HandleFunc("GET /api/v1/apis/{id}", apiHandler.Get)
	mux.HandleFunc("PUT /api/v1/apis/{id}", apiHandler.Update)
	mux.HandleFunc("DELETE /api/v1/apis/{id}", apiHandler.Delete)

	mux.HandleFunc("POST /api/v1/discovery/upload", discoveryHandler.Upload)
	mux.HandleFunc("POST /api/v1/discovery/url", discoveryHandler.FromURL)
	mux.HandleFunc("GET /api/v1/discovery/jobs", discoveryHandler.ListJobs)

	mux.HandleFunc("GET /api/v1/teams", teamHandler.List)
	mux.HandleFunc("POST /api/v1/teams", teamHandler.Create)

	mux.HandleFunc("GET /api/v1/tags", tagHandler.List)

	mux.HandleFunc("GET /api/v1/stats", statsHandler.Get)

	return middleware.Chain(mux,
		middleware.Recover,
		middleware.Logging,
		middleware.CORS(cfg.CORSOrigin),
	)
}
