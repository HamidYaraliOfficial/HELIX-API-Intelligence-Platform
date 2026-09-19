"use client";

import { useEffect, useState } from "react";
import { useI18n } from "@/lib/i18n";
import { helixApi, type API, type Meta, type Team } from "@/lib/api-client";
import { ApiForm } from "@/components/api-form";

export default function EditApiPage({ params }: { params: { id: string } }) {
  const { id } = params;
  const { t } = useI18n();
  const [api, setApi] = useState<API | null>(null);
  const [meta, setMeta] = useState<Meta | null>(null);
  const [teams, setTeams] = useState<Team[]>([]);

  useEffect(() => {
    helixApi.getAPI(id).then((res) => setApi(res.api));
    helixApi.getMeta().then(setMeta).catch(() => {});
    helixApi.listTeams().then(setTeams).catch(() => {});
  }, [id]);

  if (!api) return <div className="p-8 text-center text-sm text-fg-muted">{t("common.loading")}</div>;

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold text-fg">{t("common.edit")}</h1>
      </div>
      <ApiForm meta={meta} teams={teams} existing={api} />
    </div>
  );
}
