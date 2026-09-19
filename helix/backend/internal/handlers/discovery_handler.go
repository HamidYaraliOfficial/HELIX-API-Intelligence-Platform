package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/helix-platform/helix-backend/internal/config"
	"github.com/helix-platform/helix-backend/internal/openapi"
	"github.com/helix-platform/helix-backend/internal/repository"
)

type DiscoveryHandler struct {
	APIRepo       *repository.APIRepository
	DiscoveryRepo *repository.DiscoveryRepository
	Cfg           config.Config
}

func NewDiscoveryHandler(apiRepo *repository.APIRepository, discoveryRepo *repository.DiscoveryRepository, cfg config.Config) *DiscoveryHandler {
	return &DiscoveryHandler{APIRepo: apiRepo, DiscoveryRepo: discoveryRepo, Cfg: cfg}
}

// Upload handles POST /api/v1/discovery/upload — a multipart form upload of
// an OpenAPI/Swagger document (JSON or YAML), which is parsed and registered
// as a new cataloged API in one step.
func (h *DiscoveryHandler) Upload(w http.ResponseWriter, r *http.Request) {
	orgID := resolveOrgID(r, h.Cfg)

	if err := r.ParseMultipartForm(h.Cfg.MaxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "could not parse upload (file too large or malformed): "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing 'file' field in multipart upload")
		return
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read uploaded file")
		return
	}

	job, err := h.DiscoveryRepo.Create(orgID, header.Filename, "openapi")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	environment := r.FormValue("environment")
	teamID := r.FormValue("team_id")

	api, endpointCount, err := h.ingest(orgID, header.Filename, environment, teamID, raw)
	if err != nil {
		_ = h.DiscoveryRepo.Complete(job.ID, 0, 0, err.Error())
		writeError(w, http.StatusUnprocessableEntity, "discovery failed: "+err.Error())
		return
	}

	_ = h.DiscoveryRepo.Complete(job.ID, 1, endpointCount, "")

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"job":            job,
		"api":            api,
		"endpoint_count": endpointCount,
	})
}

type discoverURLRequest struct {
	URL         string `json:"url"`
	Environment string `json:"environment"`
	TeamID      string `json:"team_id"`
}

// FromURL handles POST /api/v1/discovery/url — fetches an OpenAPI document
// from a reachable URL (e.g. a service's own /openapi.json) and imports it.
func (h *DiscoveryHandler) FromURL(w http.ResponseWriter, r *http.Request) {
	var req discoverURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		writeError(w, http.StatusBadRequest, "request body must include a non-empty 'url'")
		return
	}

	orgID := resolveOrgID(r, h.Cfg)
	job, err := h.DiscoveryRepo.Create(orgID, req.URL, "openapi-url")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(req.URL)
	if err != nil {
		_ = h.DiscoveryRepo.Complete(job.ID, 0, 0, err.Error())
		writeError(w, http.StatusBadGateway, "could not fetch URL: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		msg := "remote server returned status " + resp.Status
		_ = h.DiscoveryRepo.Complete(job.ID, 0, 0, msg)
		writeError(w, http.StatusBadGateway, msg)
		return
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, h.Cfg.MaxUploadSize))
	if err != nil {
		_ = h.DiscoveryRepo.Complete(job.ID, 0, 0, err.Error())
		writeError(w, http.StatusBadGateway, "failed reading remote response")
		return
	}

	api, endpointCount, err := h.ingest(orgID, req.URL, req.Environment, req.TeamID, raw)
	if err != nil {
		_ = h.DiscoveryRepo.Complete(job.ID, 0, 0, err.Error())
		writeError(w, http.StatusUnprocessableEntity, "discovery failed: "+err.Error())
		return
	}

	_ = h.DiscoveryRepo.Complete(job.ID, 1, endpointCount, "")

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"job":            job,
		"api":            api,
		"endpoint_count": endpointCount,
	})
}

// ListJobs handles GET /api/v1/discovery/jobs
func (h *DiscoveryHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	orgID := resolveOrgID(r, h.Cfg)
	jobs, err := h.DiscoveryRepo.List(orgID, 50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, jobs)
}

func (h *DiscoveryHandler) ingest(orgID, sourceName, environment, teamID string, raw []byte) (interface{}, int, error) {
	contract, err := openapi.Parse(raw)
	if err != nil {
		return nil, 0, err
	}

	name := contract.Title
	if name == "" {
		name = sourceName
	}
	version := contract.Version
	if version == "" {
		version = "1.0.0"
	}
	if environment == "" {
		environment = "development"
	}

	var teamIDPtr *string
	if teamID != "" {
		teamIDPtr = &teamID
	}

	api, err := h.APIRepo.Create(repository.CreateInput{
		OrgID:          orgID,
		TeamID:         teamIDPtr,
		Name:           name,
		Description:    contract.Description,
		Protocol:       "REST",
		Version:        version,
		BaseURL:        contract.BaseURL,
		Environment:    environment,
		Visibility:     "internal",
		Status:         "unknown",
		LifecycleState: "active",
		SourceType:     "openapi",
		AuthType:       contract.AuthType,
		OperatingStart: "00:00",
		OperatingEnd:   "23:59",
		OperatingDays:  "1,2,3,4,5,6,7",
		Timezone:       "UTC",
		Tags:           []string{"discovered"},
		RawSpec:        toStorableJSON(raw),
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := h.APIRepo.ReplaceEndpointsFromContract(api.ID, contract)
	if err != nil {
		return nil, 0, err
	}

	full, err := h.APIRepo.GetByID(orgID, api.ID)
	if err != nil {
		return nil, 0, err
	}

	return full, count, nil
}

// toStorableJSON converts an uploaded spec (which may be JSON or YAML) into
// well-formed JSON bytes so it can be stored in a JSONB column, regardless
// of the format the user originally provided. Returns nil if the content
// cannot be represented as JSON, in which case the raw spec is simply not
// persisted (parsing into a Contract has already succeeded by this point).
func toStorableJSON(raw []byte) []byte {
	var generic interface{}
	if err := yaml.Unmarshal(raw, &generic); err != nil {
		return nil
	}
	out, err := json.Marshal(generic)
	if err != nil {
		return nil
	}
	return out
}
