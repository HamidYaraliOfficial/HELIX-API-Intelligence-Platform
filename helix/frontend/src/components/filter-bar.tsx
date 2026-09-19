"use client";

import { Search, X } from "lucide-react";
import { Input, Select } from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import { useI18n } from "@/lib/i18n";
import type { Meta, Team } from "@/lib/api-client";

export interface FilterState {
  q: string;
  protocol: string;
  status: string;
  environment: string;
  lifecycle_state: string;
  team_id: string;
}

export function FilterBar({
  filters,
  onChange,
  meta,
  teams,
}: {
  filters: FilterState;
  onChange: (next: FilterState) => void;
  meta: Meta | null;
  teams: Team[];
}) {
  const { t } = useI18n();

  function set<K extends keyof FilterState>(key: K, value: FilterState[K]) {
    onChange({ ...filters, [key]: value });
  }

  const hasActiveFilters =
    filters.protocol || filters.status || filters.environment || filters.lifecycle_state || filters.team_id;

  return (
    <div className="fluent-panel flex flex-col gap-3 p-3">
      <div className="relative">
        <Search className="absolute start-3 top-1/2 h-4 w-4 -translate-y-1/2 text-fg-muted" />
        <Input
          value={filters.q}
          onChange={(e) => set("q", e.target.value)}
          placeholder={t("common.search")}
          className="ps-9"
        />
      </div>

      <div className="flex flex-wrap gap-2">
        <Select value={filters.protocol} onChange={(e) => set("protocol", e.target.value)} className="w-auto min-w-[9rem]">
          <option value="">{t("filter.protocol")}</option>
          {meta?.protocols.map((p) => (
            <option key={p} value={p}>
              {p}
            </option>
          ))}
        </Select>

        <Select value={filters.status} onChange={(e) => set("status", e.target.value)} className="w-auto min-w-[9rem]">
          <option value="">{t("filter.status")}</option>
          {meta?.statuses.map((s) => (
            <option key={s} value={s}>
              {t(`status.${s}`)}
            </option>
          ))}
        </Select>

        <Select
          value={filters.environment}
          onChange={(e) => set("environment", e.target.value)}
          className="w-auto min-w-[9rem]"
        >
          <option value="">{t("filter.environment")}</option>
          {meta?.environments.map((env) => (
            <option key={env} value={env}>
              {t(`environment.${env}`)}
            </option>
          ))}
        </Select>

        <Select
          value={filters.lifecycle_state}
          onChange={(e) => set("lifecycle_state", e.target.value)}
          className="w-auto min-w-[9rem]"
        >
          <option value="">{t("filter.lifecycle")}</option>
          {meta?.lifecycle_states.map((ls) => (
            <option key={ls} value={ls}>
              {t(`lifecycle.${ls}`)}
            </option>
          ))}
        </Select>

        <Select value={filters.team_id} onChange={(e) => set("team_id", e.target.value)} className="w-auto min-w-[9rem]">
          <option value="">{t("filter.team")}</option>
          {teams.map((team) => (
            <option key={team.id} value={team.id}>
              {team.name}
            </option>
          ))}
        </Select>

        {hasActiveFilters && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onChange({ q: filters.q, protocol: "", status: "", environment: "", lifecycle_state: "", team_id: "" })}
          >
            <X className="h-3.5 w-3.5" />
            {t("common.clear")}
          </Button>
        )}
      </div>
    </div>
  );
}
