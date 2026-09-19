package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/helix-platform/helix-backend/internal/config"
	"github.com/helix-platform/helix-backend/internal/models"
	"github.com/helix-platform/helix-backend/internal/repository"
)

type APIHandler struct {
	Repo *repository.APIRepository
	Cfg  config.Config
}

func NewAPIHandler(repo *repository.APIRepository, cfg config.Config) *APIHandler {
	return &APIHandler{Repo: repo, Cfg: cfg}
}

// List handles GET /api/v1/apis
func (h *APIHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	result, err := h.Repo.List(repository.ListFilter{
		OrgID:          resolveOrgID(r, h.Cfg),
		Query:          q.Get("q"),
		Protocol:       q.Get("protocol"),
		Status:         q.Get("status"),
		Environment:    q.Get("environment"),
		Visibility:     q.Get("visibility"),
		LifecycleState: q.Get("lifecycle_state"),
		TeamID:         q.Get("team_id"),
		Tag:            q.Get("tag"),
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// Get handles GET /api/v1/apis/{id}
func (h *APIHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := resolveOrgID(r, h.Cfg)

	api, err := h.Repo.GetByID(orgID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if api == nil {
		writeError(w, http.StatusNotFound, "API not found")
		return
	}

	endpoints, err := h.Repo.GetEndpoints(api.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"api":       api,
		"endpoints": endpoints,
	})
}

type apiPayload struct {
	TeamID         *string  `json:"team_id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Owner          string   `json:"owner"`
	Protocol       string   `json:"protocol"`
	Version        string   `json:"version"`
	BaseURL        string   `json:"base_url"`
	Environment    string   `json:"environment"`
	Visibility     string   `json:"visibility"`
	Status         string   `json:"status"`
	LifecycleState string   `json:"lifecycle_state"`
	AuthType       string   `json:"auth_type"`
	OperatingStart string   `json:"operating_hours_start"`
	OperatingEnd   string   `json:"operating_hours_end"`
	OperatingDays  string   `json:"operating_days"`
	Timezone       string   `json:"timezone"`
	Tags           []string `json:"tags"`
}

func (p apiPayload) validate() string {
	if p.Name == "" {
		return "name is required"
	}
	if !contains(models.Protocols, p.Protocol) {
		return "protocol must be one of " + joinList(models.Protocols)
	}
	if !contains(models.Environments, p.Environment) {
		return "environment must be one of " + joinList(models.Environments)
	}
	if !contains(models.Visibilities, p.Visibility) {
		return "visibility must be one of " + joinList(models.Visibilities)
	}
	if !contains(models.LifecycleStates, p.LifecycleState) {
		return "lifecycle_state must be one of " + joinList(models.LifecycleStates)
	}
	return ""
}

func (p apiPayload) toCreateInput(orgID string) repository.CreateInput {
	start, end, days, tz := p.OperatingStart, p.OperatingEnd, p.OperatingDays, p.Timezone
	if start == "" {
		start = "00:00"
	}
	if end == "" {
		end = "23:59"
	}
	if days == "" {
		days = "1,2,3,4,5,6,7"
	}
	if tz == "" {
		tz = "UTC"
	}
	status := p.Status
	if status == "" {
		status = "unknown"
	}
	authType := p.AuthType
	if authType == "" {
		authType = "none"
	}

	return repository.CreateInput{
		OrgID:          orgID,
		TeamID:         p.TeamID,
		Name:           p.Name,
		Description:    p.Description,
		Owner:          p.Owner,
		Protocol:       p.Protocol,
		Version:        p.Version,
		BaseURL:        p.BaseURL,
		Environment:    p.Environment,
		Visibility:     p.Visibility,
		Status:         status,
		LifecycleState: p.LifecycleState,
		SourceType:     "manual",
		AuthType:       authType,
		OperatingStart: start,
		OperatingEnd:   end,
		OperatingDays:  days,
		Timezone:       tz,
		Tags:           p.Tags,
	}
}

// Create handles POST /api/v1/apis (manual registration)
func (h *APIHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload apiPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusUnprocessableEntity, msg)
		return
	}

	orgID := resolveOrgID(r, h.Cfg)
	api, err := h.Repo.Create(payload.toCreateInput(orgID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, api)
}

// Update handles PUT /api/v1/apis/{id}
func (h *APIHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload apiPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if msg := payload.validate(); msg != "" {
		writeError(w, http.StatusUnprocessableEntity, msg)
		return
	}

	orgID := resolveOrgID(r, h.Cfg)
	api, err := h.Repo.Update(orgID, id, payload.toCreateInput(orgID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if api == nil {
		writeError(w, http.StatusNotFound, "API not found")
		return
	}
	writeJSON(w, http.StatusOK, api)
}

// Delete handles DELETE /api/v1/apis/{id}
func (h *APIHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	orgID := resolveOrgID(r, h.Cfg)

	deleted, err := h.Repo.Delete(orgID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !deleted {
		writeError(w, http.StatusNotFound, "API not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

func joinList(list []string) string {
	out := ""
	for i, v := range list {
		if i > 0 {
			out += ", "
		}
		out += v
	}
	return out
}
