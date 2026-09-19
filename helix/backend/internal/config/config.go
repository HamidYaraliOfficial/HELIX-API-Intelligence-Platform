package config

import (
	"os"
)

// Config holds all runtime configuration for the HELIX backend.
// Every value is sourced from the environment so the service behaves
// identically whether it is run locally, in Docker Compose or in Kubernetes.
type Config struct {
	Port          string
	DatabaseURL   string
	CORSOrigin    string
	Environment   string
	DefaultOrgID  string
	MaxUploadSize int64 // bytes
}

// Load reads configuration from environment variables, falling back to
// sane local-development defaults so the server can boot without any
// manual setup.
func Load() Config {
	return Config{
		Port:          getEnv("HELIX_PORT", "8080"),
		DatabaseURL:   getEnv("HELIX_DATABASE_URL", "postgres://helix:helix@localhost:5432/helix?sslmode=disable"),
		CORSOrigin:    getEnv("HELIX_CORS_ORIGIN", "*"),
		Environment:   getEnv("HELIX_ENV", "development"),
		DefaultOrgID:  getEnv("HELIX_DEFAULT_ORG_ID", "00000000-0000-0000-0000-000000000001"),
		MaxUploadSize: 10 << 20, // 10 MB, generous enough for large OpenAPI documents
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
