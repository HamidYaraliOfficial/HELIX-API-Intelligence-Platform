const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

export interface API {
  id: string;
  org_id: string;
  team_id?: string | null;
  team_name?: string | null;
  name: string;
  description: string;
  owner: string;
  protocol: string;
  version: string;
  base_url: string;
  environment: string;
  visibility: string;
  status: string;
  lifecycle_state: string;
  source_type: string;
  auth_type: string;
  operating_hours_start: string;
  operating_hours_end: string;
  operating_days: string;
  timezone: string;
  endpoint_count: number;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface EndpointParameter {
  id: string;
  name: string;
  in: string;
  type: string;
  required: boolean;
  description: string;
}

export interface EndpointResponse {
  id: string;
  status_code: string;
  description: string;
  content_type: string;
}

export interface Endpoint {
  id: string;
  api_id: string;
  path: string;
  method: string;
  operation_id: string;
  summary: string;
  description: string;
  auth_type: string;
  deprecated: boolean;
  parameters: EndpointParameter[];
  responses: EndpointResponse[];
}

export interface PageResult<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface Team {
  id: string;
  org_id: string;
  name: string;
  slug: string;
  api_count: number;
  created_at: string;
}

export interface Tag {
  id: string;
  name: string;
}

export interface Stats {
  total_apis: number;
  total_endpoints: number;
  total_teams: number;
  by_protocol: Record<string, number>;
  by_status: Record<string, number>;
  by_environment: Record<string, number>;
  by_lifecycle_state: Record<string, number>;
}

export interface Meta {
  protocols: string[];
  environments: string[];
  visibilities: string[];
  lifecycle_states: string[];
  statuses: string[];
  auth_types: string[];
}

export interface DiscoveryJob {
  id: string;
  filename: string;
  source_type: string;
  status: string;
  apis_found: number;
  endpoints_found: number;
  error_message?: string;
  created_at: string;
  completed_at?: string | null;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers || {}),
    },
    cache: "no-store",
  });

  if (!res.ok) {
    let message = `Request failed with status ${res.status}`;
    try {
      const body = await res.json();
      if (body?.error) message = body.error;
    } catch {
      // response had no JSON body
    }
    throw new Error(message);
  }

  if (res.status === 204) return undefined as unknown as T;
  return res.json() as Promise<T>;
}

export interface ListAPIsParams {
  q?: string;
  protocol?: string;
  status?: string;
  environment?: string;
  visibility?: string;
  lifecycle_state?: string;
  team_id?: string;
  tag?: string;
  page?: number;
  page_size?: number;
}

function toQueryString(params: Record<string, string | number | undefined>): string {
  const usp = new URLSearchParams();
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== "") usp.set(k, String(v));
  }
  const s = usp.toString();
  return s ? `?${s}` : "";
}

export const helixApi = {
  listAPIs: (params: ListAPIsParams = {}) => request<PageResult<API>>(`/api/v1/apis${toQueryString(params)}`),

  getAPI: (id: string) => request<{ api: API; endpoints: Endpoint[] }>(`/api/v1/apis/${id}`),

  createAPI: (payload: Record<string, unknown>) =>
    request<API>(`/api/v1/apis`, { method: "POST", body: JSON.stringify(payload) }),

  updateAPI: (id: string, payload: Record<string, unknown>) =>
    request<API>(`/api/v1/apis/${id}`, { method: "PUT", body: JSON.stringify(payload) }),

  deleteAPI: (id: string) => request<void>(`/api/v1/apis/${id}`, { method: "DELETE" }),

  listTeams: () => request<Team[]>(`/api/v1/teams`),

  createTeam: (name: string) => request<Team>(`/api/v1/teams`, { method: "POST", body: JSON.stringify({ name }) }),

  listTags: () => request<Tag[]>(`/api/v1/tags`),

  getStats: () => request<Stats>(`/api/v1/stats`),

  getMeta: () => request<Meta>(`/api/v1/meta`),

  listDiscoveryJobs: () => request<DiscoveryJob[]>(`/api/v1/discovery/jobs`),

  uploadDiscovery: async (file: File, environment: string, teamId: string) => {
    const form = new FormData();
    form.append("file", file);
    if (environment) form.append("environment", environment);
    if (teamId) form.append("team_id", teamId);

    const res = await fetch(`${API_BASE_URL}/api/v1/discovery/upload`, { method: "POST", body: form });
    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new Error(body?.error || `Upload failed with status ${res.status}`);
    }
    return res.json() as Promise<{ api: API; endpoint_count: number }>;
  },

  discoverFromUrl: (url: string, environment: string, teamId: string) =>
    request<{ api: API; endpoint_count: number }>(`/api/v1/discovery/url`, {
      method: "POST",
      body: JSON.stringify({ url, environment, team_id: teamId }),
    }),
};
