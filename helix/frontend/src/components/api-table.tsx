"use client";

import Link from "next/link";
import { Inbox } from "lucide-react";
import { Badge, statusTone, lifecycleTone } from "@/components/ui/badge";
import { AvailabilityBadge } from "@/components/availability-badge";
import { useI18n } from "@/lib/i18n";
import type { API } from "@/lib/api-client";

export function ApiTable({ apis }: { apis: API[] }) {
  const { t, locale } = useI18n();

  if (apis.length === 0) {
    return (
      <div className="fluent-panel flex flex-col items-center justify-center gap-3 py-16 text-center">
        <Inbox className="h-10 w-10 text-fg-muted" />
        <div className="text-sm font-medium text-fg">{t("common.noResults")}</div>
        <div className="text-xs text-fg-muted">{t("common.noResultsHint")}</div>
      </div>
    );
  }

  return (
    <div className="fluent-panel overflow-x-auto">
      <table className="w-full text-start text-sm">
        <thead>
          <tr className="border-b border-border text-xs text-fg-muted">
            <th className="px-4 py-3 text-start font-medium">{t("table.name")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.team")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.protocol")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.version")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.environment")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.status")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.lifecycle")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.endpoints")}</th>
            <th className="px-4 py-3 text-start font-medium">{t("table.availability")}</th>
          </tr>
        </thead>
        <tbody>
          {apis.map((api) => (
            <tr key={api.id} className="border-b border-border last:border-0 hover:bg-bg-elevated">
              <td className="px-4 py-3">
                <Link href={`/apis/${api.id}`} className="font-medium text-fg hover:text-accent">
                  {api.name}
                </Link>
                <div className="text-xs text-fg-muted">{api.owner}</div>
              </td>
              <td className="px-4 py-3 text-fg-muted">{api.team_name || "—"}</td>
              <td className="px-4 py-3">
                <Badge tone="accent">{api.protocol}</Badge>
              </td>
              <td className="px-4 py-3 text-fg-muted">{api.version}</td>
              <td className="px-4 py-3 text-fg-muted">{t(`environment.${api.environment}`)}</td>
              <td className="px-4 py-3">
                <Badge tone={statusTone(api.status)}>{t(`status.${api.status}`)}</Badge>
              </td>
              <td className="px-4 py-3">
                <Badge tone={lifecycleTone(api.lifecycle_state)}>{t(`lifecycle.${api.lifecycle_state}`)}</Badge>
              </td>
              <td className="px-4 py-3 text-fg-muted">{api.endpoint_count}</td>
              <td className="px-4 py-3">
                <AvailabilityBadge api={api} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className="sr-only">{locale}</div>
    </div>
  );
}
