"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { UploadCloud, FileJson, CheckCircle2, XCircle, Loader2 } from "lucide-react";
import { useI18n } from "@/lib/i18n";
import { helixApi, type DiscoveryJob, type Meta, type Team } from "@/lib/api-client";
import { Button } from "@/components/ui/button";
import { FieldGroup, Input, Label, Select } from "@/components/ui/field";
import { cn } from "@/lib/cn";

type Tab = "upload" | "url";

export default function DiscoveryPage() {
  const { t } = useI18n();
  const [tab, setTab] = useState<Tab>("upload");
  const [meta, setMeta] = useState<Meta | null>(null);
  const [teams, setTeams] = useState<Team[]>([]);
  const [jobs, setJobs] = useState<DiscoveryJob[]>([]);

  const [file, setFile] = useState<File | null>(null);
  const [url, setUrl] = useState("");
  const [environment, setEnvironment] = useState("development");
  const [teamId, setTeamId] = useState("");
  const [dragActive, setDragActive] = useState(false);

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<{ id: string; name: string; count: number } | null>(null);

  const inputRef = useRef<HTMLInputElement>(null);

  function refreshJobs() {
    helixApi.listDiscoveryJobs().then(setJobs).catch(() => {});
  }

  useEffect(() => {
    helixApi.getMeta().then(setMeta).catch(() => {});
    helixApi.listTeams().then(setTeams).catch(() => {});
    refreshJobs();
  }, []);

  async function handleSubmit() {
    setSubmitting(true);
    setError(null);
    setResult(null);
    try {
      if (tab === "upload") {
        if (!file) throw new Error("Please choose a file");
        const res = await helixApi.uploadDiscovery(file, environment, teamId);
        setResult({ id: res.api.id, name: res.api.name, count: res.endpoint_count });
      } else {
        if (!url) throw new Error("Please enter a URL");
        const res = await helixApi.discoverFromUrl(url, environment, teamId);
        setResult({ id: res.api.id, name: res.api.name, count: res.endpoint_count });
      }
      refreshJobs();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto flex max-w-3xl flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold text-fg">{t("discovery.title")}</h1>
        <p className="text-sm text-fg-muted">{t("discovery.subtitle")}</p>
      </div>

      <div className="flex gap-2">
        <TabButton active={tab === "upload"} onClick={() => setTab("upload")}>
          {t("discovery.uploadTab")}
        </TabButton>
        <TabButton active={tab === "url"} onClick={() => setTab("url")}>
          {t("discovery.urlTab")}
        </TabButton>
      </div>

      <div className="fluent-panel flex flex-col gap-4 p-5">
        {tab === "upload" ? (
          <div
            onDragOver={(e) => {
              e.preventDefault();
              setDragActive(true);
            }}
            onDragLeave={() => setDragActive(false)}
            onDrop={(e) => {
              e.preventDefault();
              setDragActive(false);
              const dropped = e.dataTransfer.files?.[0];
              if (dropped) setFile(dropped);
            }}
            onClick={() => inputRef.current?.click()}
            className={cn(
              "flex cursor-pointer flex-col items-center justify-center gap-2 rounded-win-lg border-2 border-dashed p-10 text-center transition-colors",
              dragActive ? "border-accent bg-[var(--color-accent-bg)]" : "border-border"
            )}
          >
            <input
              ref={inputRef}
              type="file"
              accept=".json,.yaml,.yml"
              className="hidden"
              onChange={(e) => setFile(e.target.files?.[0] || null)}
            />
            {file ? (
              <>
                <FileJson className="h-8 w-8 text-accent" />
                <div className="text-sm font-medium text-fg">{file.name}</div>
                <div className="text-xs text-fg-muted">{t("discovery.selectedFile")}</div>
              </>
            ) : (
              <>
                <UploadCloud className="h-8 w-8 text-fg-muted" />
                <div className="text-sm text-fg">{t("discovery.dropHint")}</div>
                <Button type="button" variant="secondary" size="sm" onClick={() => inputRef.current?.click()}>
                  {t("discovery.chooseFile")}
                </Button>
              </>
            )}
          </div>
        ) : (
          <FieldGroup>
            <Label>{t("discovery.urlLabel")}</Label>
            <Input value={url} onChange={(e) => setUrl(e.target.value)} placeholder={t("discovery.urlPlaceholder")} />
          </FieldGroup>
        )}

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <FieldGroup>
            <Label>{t("discovery.environmentLabel")}</Label>
            <Select value={environment} onChange={(e) => setEnvironment(e.target.value)}>
              {(meta?.environments || ["development"]).map((env) => (
                <option key={env} value={env}>
                  {t(`environment.${env}`)}
                </option>
              ))}
            </Select>
          </FieldGroup>
          <FieldGroup>
            <Label>{t("discovery.teamLabel")}</Label>
            <Select value={teamId} onChange={(e) => setTeamId(e.target.value)}>
              <option value="">{t("common.none")}</option>
              {teams.map((team) => (
                <option key={team.id} value={team.id}>
                  {team.name}
                </option>
              ))}
            </Select>
          </FieldGroup>
        </div>

        {error && <div className="text-sm text-danger">{error}</div>}

        {result && (
          <div className="flex items-center justify-between rounded-win border border-[var(--color-success)] bg-[var(--color-success-bg)] px-4 py-3 text-sm">
            <span className="flex items-center gap-2 text-success">
              <CheckCircle2 className="h-4 w-4" />
              {t("discovery.success", { name: result.name, count: result.count })}
            </span>
            <Link href={`/apis/${result.id}`} className="font-medium text-accent hover:underline">
              {t("discovery.viewApi")}
            </Link>
          </div>
        )}

        <Button onClick={handleSubmit} disabled={submitting} className="self-start">
          {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
          {submitting ? t("discovery.running") : t("discovery.submit")}
        </Button>
      </div>

      {jobs.length > 0 && (
        <div>
          <h2 className="mb-2 text-sm font-semibold text-fg">{t("discovery.recentJobs")}</h2>
          <div className="fluent-panel divide-y divide-border">
            {jobs.map((job) => (
              <div key={job.id} className="flex items-center gap-3 px-4 py-2.5 text-sm">
                {job.status === "completed" && <CheckCircle2 className="h-4 w-4 flex-shrink-0 text-success" />}
                {job.status === "failed" && <XCircle className="h-4 w-4 flex-shrink-0 text-danger" />}
                {job.status === "running" && <Loader2 className="h-4 w-4 flex-shrink-0 animate-spin text-accent" />}
                <span className="flex-1 truncate text-fg">{job.filename}</span>
                <span className="text-xs text-fg-muted">
                  {job.apis_found} / {job.endpoints_found} eps
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function TabButton({ active, onClick, children }: { active: boolean; onClick: () => void; children: React.ReactNode }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "rounded-win px-3 py-1.5 text-sm font-medium",
        active ? "bg-accent text-accent-fg" : "bg-bg-elevated text-fg-muted hover:text-fg"
      )}
    >
      {children}
    </button>
  );
}
