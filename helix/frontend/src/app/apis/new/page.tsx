"use client";

import { useEffect, useState } from "react";
import { useI18n } from "@/lib/i18n";
import { helixApi, type Meta, type Team } from "@/lib/api-client";
import { ApiForm } from "@/components/api-form";

export default function NewApiPage() {
  const { t } = useI18n();
  const [meta, setMeta] = useState<Meta | null>(null);
  const [teams, setTeams] = useState<Team[]>([]);

  useEffect(() => {
    helixApi.getMeta().then(setMeta).catch(() => {});
    helixApi.listTeams().then(setTeams).catch(() => {});
  }, []);

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold text-fg">{t("common.create")}</h1>
      </div>
      <ApiForm meta={meta} teams={teams} />
    </div>
  );
}
