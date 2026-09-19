"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { FieldGroup, Input, Label, Select, Textarea } from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import { useI18n } from "@/lib/i18n";
import { helixApi, type API, type Meta, type Team } from "@/lib/api-client";

const DAY_CODES = [1, 2, 3, 4, 5, 6, 7];
const DAY_LABEL_KEYS = ["mon", "tue", "wed", "thu", "fri", "sat", "sun"];

const COMMON_TIMEZONES = [
  "UTC",
  "Asia/Tehran",
  "Asia/Dubai",
  "Asia/Shanghai",
  "Asia/Tokyo",
  "Europe/London",
  "Europe/Berlin",
  "America/New_York",
  "America/Los_Angeles",
];

export interface ApiFormValues {
  team_id: string;
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
  auth_type: string;
  operating_hours_start: string;
  operating_hours_end: string;
  operating_days: string;
  timezone: string;
  tags: string;
}

function fromApi(api: API): ApiFormValues {
  return {
    team_id: api.team_id || "",
    name: api.name,
    description: api.description,
    owner: api.owner,
    protocol: api.protocol,
    version: api.version,
    base_url: api.base_url,
    environment: api.environment,
    visibility: api.visibility,
    status: api.status,
    lifecycle_state: api.lifecycle_state,
    auth_type: api.auth_type,
    operating_hours_start: api.operating_hours_start,
    operating_hours_end: api.operating_hours_end,
    operating_days: api.operating_days,
    timezone: api.timezone,
    tags: api.tags?.join(", ") || "",
  };
}

const defaults: ApiFormValues = {
  team_id: "",
  name: "",
  description: "",
  owner: "",
  protocol: "REST",
  version: "1.0.0",
  base_url: "",
  environment: "development",
  visibility: "internal",
  status: "unknown",
  lifecycle_state: "draft",
  auth_type: "none",
  operating_hours_start: "00:00",
  operating_hours_end: "23:59",
  operating_days: "1,2,3,4,5,6,7",
  timezone: "UTC",
  tags: "",
};

export function ApiForm({ meta, teams, existing }: { meta: Meta | null; teams: Team[]; existing?: API }) {
  const { t } = useI18n();
  const router = useRouter();
  const [values, setValues] = useState<ApiFormValues>(existing ? fromApi(existing) : defaults);
  const [days, setDays] = useState<number[]>(
    (existing?.operating_days || defaults.operating_days).split(",").map((d) => parseInt(d, 10))
  );
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  function set<K extends keyof ApiFormValues>(key: K, value: ApiFormValues[K]) {
    setValues((v) => ({ ...v, [key]: value }));
  }

  function toggleDay(day: number) {
    setDays((prev) => (prev.includes(day) ? prev.filter((d) => d !== day) : [...prev, day].sort()));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const payload = {
        ...values,
        team_id: values.team_id || null,
        operating_days: days.join(","),
        tags: values.tags
          .split(",")
          .map((tag) => tag.trim())
          .filter(Boolean),
      };

      const result = existing ? await helixApi.updateAPI(existing.id, payload) : await helixApi.createAPI(payload);
      router.push(`/apis/${result.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-6">
      {error && <div className="fluent-panel border-[var(--color-danger)] p-3 text-sm text-danger">{error}</div>}

      <div className="fluent-panel grid grid-cols-1 gap-4 p-5 sm:grid-cols-2">
        <FieldGroup>
          <Label>{t("form.name")}</Label>
          <Input required value={values.name} onChange={(e) => set("name", e.target.value)} />
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.owner")}</Label>
          <Input value={values.owner} onChange={(e) => set("owner", e.target.value)} />
        </FieldGroup>

        <div className="sm:col-span-2">
          <FieldGroup>
            <Label>{t("form.description")}</Label>
            <Textarea rows={3} value={values.description} onChange={(e) => set("description", e.target.value)} />
          </FieldGroup>
        </div>

        <FieldGroup>
          <Label>{t("form.team")}</Label>
          <Select value={values.team_id} onChange={(e) => set("team_id", e.target.value)}>
            <option value="">{t("common.none")}</option>
            {teams.map((team) => (
              <option key={team.id} value={team.id}>
                {team.name}
              </option>
            ))}
          </Select>
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.protocol")}</Label>
          <Select value={values.protocol} onChange={(e) => set("protocol", e.target.value)}>
            {(meta?.protocols || ["REST"]).map((p) => (
              <option key={p} value={p}>
                {p}
              </option>
            ))}
          </Select>
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.version")}</Label>
          <Input value={values.version} onChange={(e) => set("version", e.target.value)} />
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.baseUrl")}</Label>
          <Input value={values.base_url} onChange={(e) => set("base_url", e.target.value)} placeholder="https://" />
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.environment")}</Label>
          <Select value={values.environment} onChange={(e) => set("environment", e.target.value)}>
            {(meta?.environments || ["development"]).map((env) => (
              <option key={env} value={env}>
                {t(`environment.${env}`)}
              </option>
            ))}
          </Select>
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.visibility")}</Label>
          <Select value={values.visibility} onChange={(e) => set("visibility", e.target.value)}>
            {(meta?.visibilities || ["internal"]).map((v) => (
              <option key={v} value={v}>
                {t(`visibility.${v}`)}
              </option>
            ))}
          </Select>
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.status")}</Label>
          <Select value={values.status} onChange={(e) => set("status", e.target.value)}>
            {(meta?.statuses || ["unknown"]).map((s) => (
              <option key={s} value={s}>
                {t(`status.${s}`)}
              </option>
            ))}
          </Select>
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.lifecycle")}</Label>
          <Select value={values.lifecycle_state} onChange={(e) => set("lifecycle_state", e.target.value)}>
            {(meta?.lifecycle_states || ["draft"]).map((ls) => (
              <option key={ls} value={ls}>
                {t(`lifecycle.${ls}`)}
              </option>
            ))}
          </Select>
        </FieldGroup>

        <FieldGroup>
          <Label>{t("form.authType")}</Label>
          <Select value={values.auth_type} onChange={(e) => set("auth_type", e.target.value)}>
            {(meta?.auth_types || ["none"]).map((a) => (
              <option key={a} value={a}>
                {a}
              </option>
            ))}
          </Select>
        </FieldGroup>

        <div className="sm:col-span-2">
          <FieldGroup>
            <Label>{t("form.tags")}</Label>
            <Input value={values.tags} onChange={(e) => set("tags", e.target.value)} placeholder={t("form.tagsHint")} />
          </FieldGroup>
        </div>
      </div>

      <div className="fluent-panel flex flex-col gap-4 p-5">
        <div>
          <h3 className="text-sm font-semibold text-fg">{t("form.availabilityTitle")}</h3>
          <p className="text-xs text-fg-muted">{t("form.availabilityHint")}</p>
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <FieldGroup>
            <Label>{t("form.hoursStart")}</Label>
            <Input
              type="time"
              value={values.operating_hours_start}
              onChange={(e) => set("operating_hours_start", e.target.value)}
            />
          </FieldGroup>
          <FieldGroup>
            <Label>{t("form.hoursEnd")}</Label>
            <Input
              type="time"
              value={values.operating_hours_end}
              onChange={(e) => set("operating_hours_end", e.target.value)}
            />
          </FieldGroup>
          <FieldGroup>
            <Label>{t("form.timezone")}</Label>
            <Select value={values.timezone} onChange={(e) => set("timezone", e.target.value)}>
              {COMMON_TIMEZONES.map((tz) => (
                <option key={tz} value={tz}>
                  {tz}
                </option>
              ))}
            </Select>
          </FieldGroup>
        </div>

        <FieldGroup>
          <Label>{t("form.days")}</Label>
          <div className="flex flex-wrap gap-2">
            {DAY_CODES.map((day, i) => (
              <button
                key={day}
                type="button"
                onClick={() => toggleDay(day)}
                className={
                  days.includes(day)
                    ? "rounded-win border border-accent bg-[var(--color-accent-bg)] px-3 py-1.5 text-xs font-medium text-accent"
                    : "rounded-win border border-border bg-bg-elevated px-3 py-1.5 text-xs font-medium text-fg-muted"
                }
              >
                {t(`form.day.${DAY_LABEL_KEYS[i]}`)}
              </button>
            ))}
          </div>
        </FieldGroup>
      </div>

      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={() => router.back()}>
          {t("common.cancel")}
        </Button>
        <Button type="submit" disabled={submitting}>
          {existing ? t("form.submitEdit") : t("form.submit")}
        </Button>
      </div>
    </form>
  );
}
