-- HELIX API Intelligence Platform — initial schema
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS organizations (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    slug       TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS teams (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (org_id, slug)
);

CREATE TABLE IF NOT EXISTS apis (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id                 UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    team_id                UUID REFERENCES teams(id) ON DELETE SET NULL,
    name                   TEXT NOT NULL,
    description            TEXT NOT NULL DEFAULT '',
    owner                  TEXT NOT NULL DEFAULT '',
    protocol               TEXT NOT NULL DEFAULT 'REST',
    version                TEXT NOT NULL DEFAULT '1.0.0',
    base_url               TEXT NOT NULL DEFAULT '',
    environment            TEXT NOT NULL DEFAULT 'development',
    visibility             TEXT NOT NULL DEFAULT 'internal',
    status                 TEXT NOT NULL DEFAULT 'unknown',
    lifecycle_state        TEXT NOT NULL DEFAULT 'draft',
    source_type            TEXT NOT NULL DEFAULT 'manual',
    auth_type              TEXT NOT NULL DEFAULT 'none',
    operating_hours_start  TEXT NOT NULL DEFAULT '00:00',
    operating_hours_end    TEXT NOT NULL DEFAULT '23:59',
    operating_days         TEXT NOT NULL DEFAULT '1,2,3,4,5,6,7',
    timezone               TEXT NOT NULL DEFAULT 'UTC',
    raw_spec               JSONB,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_apis_org ON apis (org_id);
CREATE INDEX IF NOT EXISTS idx_apis_team ON apis (team_id);
CREATE INDEX IF NOT EXISTS idx_apis_protocol ON apis (protocol);
CREATE INDEX IF NOT EXISTS idx_apis_status ON apis (status);
CREATE INDEX IF NOT EXISTS idx_apis_environment ON apis (environment);
CREATE INDEX IF NOT EXISTS idx_apis_lifecycle ON apis (lifecycle_state);
CREATE INDEX IF NOT EXISTS idx_apis_search ON apis USING gin (
    to_tsvector('simple', coalesce(name, '') || ' ' || coalesce(description, '') || ' ' || coalesce(owner, ''))
);

CREATE TABLE IF NOT EXISTS endpoints (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_id       UUID NOT NULL REFERENCES apis(id) ON DELETE CASCADE,
    path         TEXT NOT NULL,
    method       TEXT NOT NULL,
    operation_id TEXT NOT NULL DEFAULT '',
    summary      TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    auth_type    TEXT NOT NULL DEFAULT '',
    deprecated   BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_endpoints_api ON endpoints (api_id);

CREATE TABLE IF NOT EXISTS endpoint_parameters (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint_id UUID NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    location    TEXT NOT NULL DEFAULT 'query',
    type        TEXT NOT NULL DEFAULT 'string',
    required    BOOLEAN NOT NULL DEFAULT false,
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_params_endpoint ON endpoint_parameters (endpoint_id);

CREATE TABLE IF NOT EXISTS endpoint_responses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint_id  UUID NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    status_code  TEXT NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    content_type TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_responses_endpoint ON endpoint_responses (endpoint_id);

CREATE TABLE IF NOT EXISTS tags (
    id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name   TEXT NOT NULL,
    UNIQUE (org_id, name)
);

CREATE TABLE IF NOT EXISTS api_tags (
    api_id UUID NOT NULL REFERENCES apis(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (api_id, tag_id)
);

CREATE TABLE IF NOT EXISTS discovery_jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id          UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    filename        TEXT NOT NULL DEFAULT '',
    source_type     TEXT NOT NULL DEFAULT 'openapi',
    status          TEXT NOT NULL DEFAULT 'pending',
    apis_found      INT NOT NULL DEFAULT 0,
    endpoints_found INT NOT NULL DEFAULT 0,
    error_message   TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_jobs_org ON discovery_jobs (org_id);

-- Seed a default organization and team so the platform is usable immediately
-- after `docker compose up`, without a separate onboarding step.
INSERT INTO organizations (id, name, slug)
VALUES ('00000000-0000-0000-0000-000000000001', 'Default Organization', 'default')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO teams (id, org_id, name, slug)
VALUES ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Platform Team', 'platform-team')
ON CONFLICT (org_id, slug) DO NOTHING;
