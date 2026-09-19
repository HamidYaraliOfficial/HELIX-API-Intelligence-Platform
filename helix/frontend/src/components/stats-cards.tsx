import { Boxes, Waypoints, Users2, ShieldCheck } from "lucide-react";
import type { Stats } from "@/lib/api-client";

function StatCard({ icon: Icon, label, value }: { icon: any; label: string; value: number }) {
  return (
    <div className="fluent-panel flex items-center gap-4 p-4">
      <div className="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-win-lg bg-[var(--color-accent-bg)] text-accent">
        <Icon className="h-5 w-5" />
      </div>
      <div>
        <div className="text-2xl font-semibold leading-tight text-fg">{value.toLocaleString()}</div>
        <div className="text-xs text-fg-muted">{label}</div>
      </div>
    </div>
  );
}

export function StatsCards({ stats, t }: { stats: Stats; t: (key: string) => string }) {
  const active = stats.by_lifecycle_state?.active ?? 0;
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <StatCard icon={Boxes} label={t("dashboard.stat.totalApis")} value={stats.total_apis} />
      <StatCard icon={Waypoints} label={t("dashboard.stat.totalEndpoints")} value={stats.total_endpoints} />
      <StatCard icon={Users2} label={t("dashboard.stat.totalTeams")} value={stats.total_teams} />
      <StatCard icon={ShieldCheck} label={t("dashboard.stat.active")} value={active} />
    </div>
  );
}
