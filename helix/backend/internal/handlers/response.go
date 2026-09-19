package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/helix-platform/helix-backend/internal/config"
)

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// resolveOrgID reads the tenant from the X-Org-Id header, falling back to
// the configured default organization. This keeps the MVP usable without a
// full auth flow while still keeping every query tenant-scoped.
func resolveOrgID(r *http.Request, cfg config.Config) string {
	if id := r.Header.Get("X-Org-Id"); id != "" {
		return id
	}
	return cfg.DefaultOrgID
}
