package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/helix-platform/helix-backend/internal/config"
	"github.com/helix-platform/helix-backend/internal/repository"
)

type TeamHandler struct {
	Repo *repository.TeamRepository
	Cfg  config.Config
}

func NewTeamHandler(repo *repository.TeamRepository, cfg config.Config) *TeamHandler {
	return &TeamHandler{Repo: repo, Cfg: cfg}
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	teams, err := h.Repo.List(resolveOrgID(r, h.Cfg))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, teams)
}

type teamPayload struct {
	Name string `json:"name"`
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p teamPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil || strings.TrimSpace(p.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	slug := slugify(p.Name)
	team, err := h.Repo.Create(resolveOrgID(r, h.Cfg), p.Name, slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, team)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

type TagHandler struct {
	Repo *repository.TagRepository
	Cfg  config.Config
}

func NewTagHandler(repo *repository.TagRepository, cfg config.Config) *TagHandler {
	return &TagHandler{Repo: repo, Cfg: cfg}
}

func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	tags, err := h.Repo.List(resolveOrgID(r, h.Cfg))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tags)
}
