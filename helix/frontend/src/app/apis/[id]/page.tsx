"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowLeft, Pencil, Trash2, ExternalLink } from "lucide-react";
import { useI18n } from "@/lib/i18n";
import { helixApi, type API, type Endpoint } from "@/lib/api-client";
import { Badge, statusTone, lifecycleTone } from "@/components/ui/badge";
import { AvailabilityBadge } from "@/components/availability-badge";
import { Button } from "@/components/ui/button";

const methodTone: Record<string, "accent" | "success" | "warning" | "danger" | "neutral"> = {
  GET: "accent",
  POST: "success",
  PUT: "warning",
  PATCH: "warning",
  DELETE: "danger",
};

export default function ApiDetailPage({ params }: { params: { id: string } }) {
  const { id } = params;
  const { t } = useI18n();
  const router = useRouter();

  const [api, setApi] = useState<API | null>(null);
  const [endpoints, setEndpoints] = useState<Endpoint[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [expanded, setExpanded] = useState<string | null>(null);

  useEffect(() => {
    helixApi
      .getAPI(id)
      .then((res) => {
        setApi(res.api);
        setEndpoints(res.endpoints);
      })
      .catch((err) => setError(err instanceof Error ? err.message : String(err)));
  }, [id]);

  async function handleDelete() {
    if (!api) return;
    if (!window.confirm(t("common.confirmDeleteBody"))) return;
    await helixApi.deleteAPI(api.id);
    router.push("/");
  }

  if (error) return <div className="fluent-panel p-4 text-sm text-danger">{error}</div>;
  if (!api) return <div className="p-8 text-center text-sm text-fg-muted">{t("common.loading")}</div>;

  return (
    <div className="mx-auto flex max-w-5xl flex-col gap-6">
      <Link href="/" className="flex w-fit items-center gap-1.5 text-sm text-fg-muted hover:text-fg">
        <ArrowLeft className="h-3.5 w-3.5" />
        {t("common.back")}
      </Link>

      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-semibold text-fg">{api.name}</h1>
            <Badge tone="accent">{api.protocol}</Badge>
            <Badge tone={statusTone(api.status)}>{t(`status.${api.status}`)}</Badge>
            <Badge tone={lifecycleTone(api.lifecycle_state)}>{t(`lifecycle.${api.lifecycle_state}`)}</Badge>
          </div>
          <p className="mt-1 max-w-2xl text-sm text-fg-muted">{api.description || "—"}</p>
        </div>
        <div className="flex gap-2">
          <Link href={`/apis/${api.id}/edit`}>
            <Button variant="secondary" size="sm">
              <Pencil className="h-3.5 w-3.5" />
              {t("common.edit")}
            </Button>
          </Link>
          <Button variant="danger" size="sm" onClick={handleDelete}>
            <Trash2 className="h-3.5 w-3.5" />
            {t("common.delete")}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <InfoBlock label={t("table.owner")} value={api.owner || "—"} />
        <InfoBlock label={t("table.team")} value={api.team_name || "—"} />
        <InfoBlock label={t("table.version")} value={api.version} />
        <InfoBlock label={t("table.environment")} value={t(`environment.${api.environment}`)} />
        <InfoBlock label={t("detail.baseUrl")}>
          {api.base_url ? (
            <a
              href={api.base_url}
              target="_blank"
              rel="noreferrer"
              className="flex items-center gap-1 text-accent hover:underline"
            >
              {api.base_url}
              <ExternalLink className="h-3 w-3" />
            </a>
          ) : (
            "—"
          )}
        </InfoBlock>
        <InfoBlock label={t("detail.authType")} value={api.auth_type} />
        <InfoBlock label={t("table.availability")}>
          <AvailabilityBadge api={api} />
        </InfoBlock>
        <InfoBlock label={t("detail.tags")}>
          {api.tags && api.tags.length > 0 ? (
            <div className="flex flex-wrap gap-1">
              {api.tags.map((tag) => (
                <Badge key={tag}>{tag}</Badge>
              ))}
            </div>
          ) : (
            t("detail.noTags")
          )}
        </InfoBlock>
      </div>

      <div>
        <h2 className="mb-3 text-lg font-semibold text-fg">
          {t("detail.endpoints")} ({endpoints.length})
        </h2>

        {endpoints.length === 0 ? (
          <div className="fluent-panel p-6 text-center text-sm text-fg-muted">{t("detail.noEndpoints")}</div>
        ) : (
          <div className="fluent-panel divide-y divide-border">
            {endpoints.map((ep) => (
              <div key={ep.id}>
                <button
                  type="button"
                  onClick={() => setExpanded(expanded === ep.id ? null : ep.id)}
                  className="flex w-full items-center gap-3 px-4 py-3 text-start hover:bg-bg-elevated"
                >
                  <Badge tone={methodTone[ep.method] || "neutral"} className="w-16 justify-center">
                    {ep.method}
                  </Badge>
                  <span className="flex-1 font-mono text-sm text-fg">{ep.path}</span>
                  {ep.deprecated && <Badge tone="warning">{t("detail.deprecated")}</Badge>}
                  <span className="hidden text-xs text-fg-muted sm:block">{ep.summary}</span>
                </button>

                {expanded === ep.id && (
                  <div className="grid grid-cols-1 gap-4 border-t border-border bg-bg px-4 py-3 sm:grid-cols-2">
                    <div>
                      <div className="mb-1.5 text-xs font-semibold uppercase text-fg-muted">
                        {t("detail.parameters")}
                      </div>
                      {ep.parameters.length === 0 ? (
                        <div className="text-xs text-fg-muted">—</div>
                      ) : (
                        <ul className="space-y-1 text-xs">
                          {ep.parameters.map((p) => (
                            <li key={p.id} className="flex items-center gap-1.5">
                              <span className="font-mono text-fg">{p.name}</span>
                              <span className="text-fg-muted">({p.in}, {p.type || "string"})</span>
                              {p.required && <Badge tone="danger">required</Badge>}
                            </li>
                          ))}
                        </ul>
                      )}
                    </div>
                    <div>
                      <div className="mb-1.5 text-xs font-semibold uppercase text-fg-muted">
                        {t("detail.responses")}
                      </div>
                      {ep.responses.length === 0 ? (
                        <div className="text-xs text-fg-muted">—</div>
                      ) : (
                        <ul className="space-y-1 text-xs">
                          {ep.responses.map((r) => (
                            <li key={r.id} className="flex items-center gap-1.5">
                              <span className="font-mono text-fg">{r.status_code}</span>
                              <span className="text-fg-muted">{r.description}</span>
                            </li>
                          ))}
                        </ul>
                      )}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function InfoBlock({ label, value, children }: { label: string; value?: string; children?: React.ReactNode }) {
  return (
    <div className="fluent-panel p-3">
      <div className="text-xs text-fg-muted">{label}</div>
      <div className="mt-1 text-sm text-fg">{children ?? value}</div>
    </div>
  );
}
