"use client";

import { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import { Plus, RefreshCw } from "lucide-react";
import { useI18n } from "@/lib/i18n";
import { helixApi, type API, type Stats, type Meta, type Team } from "@/lib/api-client";
import { StatsCards } from "@/components/stats-cards";
import { FilterBar, type FilterState } from "@/components/filter-bar";
import { ApiTable } from "@/components/api-table";
import { Button } from "@/components/ui/button";

const emptyFilters: FilterState = { q: "", protocol: "", status: "", environment: "", lifecycle_state: "", team_id: "" };

export default function DashboardPage() {
  const { t } = useI18n();
  const [apis, setApis] = useState<API[] | null>(null);
  const [stats, setStats] = useState<Stats | null>(null);
  const [meta, setMeta] = useState<Meta | null>(null);
  const [teams, setTeams] = useState<Team[]>([]);
  const [filters, setFilters] = useState<FilterState>(emptyFilters);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [apiPage, statsRes] = await Promise.all([
        helixApi.listAPIs({ ...filters, page_size: 100 }),
        helixApi.getStats(),
      ]);
      setApis(apiPage.items);
      setStats(statsRes);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    helixApi.getMeta().then(setMeta).catch(() => {});
    helixApi.listTeams().then(setTeams).catch(() => {});
  }, []);

  useEffect(() => {
    const id = setTimeout(load, 200); // debounce search typing
    return () => clearTimeout(id);
  }, [load]);

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-fg">{t("dashboard.title")}</h1>
          <p className="text-sm text-fg-muted">{t("dashboard.subtitle")}</p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" size="sm" onClick={load} disabled={loading}>
            <RefreshCw className={loading ? "h-4 w-4 animate-spin" : "h-4 w-4"} />
          </Button>
          <Link href="/apis/new">
            <Button size="sm">
              <Plus className="h-4 w-4" />
              {t("common.create")}
            </Button>
          </Link>
        </div>
      </div>

      {stats && <StatsCards stats={stats} t={t} />}

      <FilterBar filters={filters} onChange={setFilters} meta={meta} teams={teams} />

      {error && (
        <div className="fluent-panel border-[var(--color-danger)] p-4 text-sm text-danger">
          {t("common.error")}: {error}
        </div>
      )}

      {apis ? <ApiTable apis={apis} /> : <div className="p-8 text-center text-sm text-fg-muted">{t("common.loading")}</div>}
    </div>
  );
}
