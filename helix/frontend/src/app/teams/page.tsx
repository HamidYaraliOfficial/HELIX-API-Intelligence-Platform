"use client";

import { useEffect, useState } from "react";
import { Users2, Plus } from "lucide-react";
import { useI18n } from "@/lib/i18n";
import { helixApi, type Team } from "@/lib/api-client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/field";

export default function TeamsPage() {
  const { t } = useI18n();
  const [teams, setTeams] = useState<Team[]>([]);
  const [name, setName] = useState("");
  const [submitting, setSubmitting] = useState(false);

  function load() {
    helixApi.listTeams().then(setTeams).catch(() => {});
  }

  useEffect(load, []);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!name.trim()) return;
    setSubmitting(true);
    try {
      await helixApi.createTeam(name.trim());
      setName("");
      load();
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold text-fg">{t("teams.title")}</h1>
        <p className="text-sm text-fg-muted">{t("teams.subtitle")}</p>
      </div>

      <form onSubmit={handleCreate} className="fluent-panel flex gap-2 p-3">
        <Input value={name} onChange={(e) => setName(e.target.value)} placeholder={t("teams.namePlaceholder")} />
        <Button type="submit" disabled={submitting}>
          <Plus className="h-4 w-4" />
          {t("teams.newTeam")}
        </Button>
      </form>

      <div className="fluent-panel divide-y divide-border">
        {teams.map((team) => (
          <div key={team.id} className="flex items-center gap-3 px-4 py-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-win-lg bg-[var(--color-accent-bg)] text-accent">
              <Users2 className="h-4 w-4" />
            </div>
            <div className="flex-1">
              <div className="text-sm font-medium text-fg">{team.name}</div>
              <div className="text-xs text-fg-muted">{team.slug}</div>
            </div>
            <div className="text-sm text-fg-muted">{t("teams.apiCount", { count: team.api_count })}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
