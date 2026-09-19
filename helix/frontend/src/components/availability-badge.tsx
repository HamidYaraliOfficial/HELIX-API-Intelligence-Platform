"use client";

import { useEffect, useState } from "react";
import { Clock } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { getAvailabilityStatus, formatMinutes, type AvailabilityInput } from "@/lib/availability";
import { useI18n } from "@/lib/i18n";

export function AvailabilityBadge({ api }: { api: AvailabilityInput }) {
  const { t, locale } = useI18n();
  const [tick, setTick] = useState(0);

  // Recompute once a minute so the countdown stays accurate without a
  // full page reload.
  useEffect(() => {
    const id = setInterval(() => setTick((n) => n + 1), 60_000);
    return () => clearInterval(id);
  }, []);

  const status = getAvailabilityStatus(api);
  void tick;

  if (status.isAlwaysOpen) {
    return (
      <Badge tone="success" className="items-center">
        <Clock className="h-3 w-3" />
        {t("availability.alwaysOpen")}
      </Badge>
    );
  }

  const timeLabel = status.minutesUntilChange !== null ? formatMinutes(status.minutesUntilChange, locale) : "";

  return (
    <Badge tone={status.isOpen ? "success" : "neutral"} className="items-center">
      <Clock className="h-3 w-3" />
      {status.isOpen ? t("availability.closesIn", { time: timeLabel }) : t("availability.opensIn", { time: timeLabel })}
    </Badge>
  );
}
