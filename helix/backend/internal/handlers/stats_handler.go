package handlers

import (
	"net/http"

	"github.com/helix-platform/helix-backend/internal/config"
	"github.com/helix-platform/helix-backend/internal/models"
	"github.com/helix-platform/helix-backend/internal/repository"
)

type StatsHandler struct {
	Repo *repository.StatsRepository
	Cfg  config.Config
}

func NewStatsHandler(repo *repository.StatsRepository, cfg config.Config) *StatsHandler {
	return &StatsHandler{Repo: repo, Cfg: cfg}
}

func (h *StatsHandler) Get(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Repo.Get(resolveOrgID(r, h.Cfg))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// Meta handles GET /api/v1/meta — returns every enum value the frontend
// needs for filter dropdowns and forms, kept in one place so backend and
// frontend never drift out of sync.
func Meta(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"protocols":        models.Protocols,
		"environments":     models.Environments,
		"visibilities":     models.Visibilities,
		"lifecycle_states": models.LifecycleStates,
		"statuses":         models.Statuses,
		"auth_types":       models.AuthTypes,
	})
}

// Health handles GET /health for container orchestrators (readiness/liveness).
func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "helix-backend"})
}
