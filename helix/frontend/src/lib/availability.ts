// Computes whether an API is currently "open" per its configured operating
// hours, and how long until the next state change — entirely client-side,
// timezone-aware, and recomputed on every render tick.

export interface AvailabilityInput {
  operating_hours_start: string; // "HH:MM"
  operating_hours_end: string; // "HH:MM"
  operating_days: string; // "1,2,3,4,5,6,7" (1 = Monday ... 7 = Sunday)
  timezone: string; // IANA timezone, e.g. "Asia/Tehran"
}

export interface AvailabilityStatus {
  isOpen: boolean;
  isAlwaysOpen: boolean;
  minutesUntilChange: number | null;
}

function nowPartsInTimezone(timezone: string): { weekday: number; minutesOfDay: number } {
  const formatter = new Intl.DateTimeFormat("en-US", {
    timeZone: timezone || "UTC",
    weekday: "short",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });

  const parts = formatter.formatToParts(new Date());
  const weekdayStr = parts.find((p) => p.type === "weekday")?.value ?? "Mon";
  const hourStr = parts.find((p) => p.type === "hour")?.value ?? "00";
  const minuteStr = parts.find((p) => p.type === "minute")?.value ?? "00";

  const weekdayMap: Record<string, number> = {
    Mon: 1,
    Tue: 2,
    Wed: 3,
    Thu: 4,
    Fri: 5,
    Sat: 6,
    Sun: 7,
  };

  const hour = parseInt(hourStr, 10) % 24;
  const minute = parseInt(minuteStr, 10);

  return {
    weekday: weekdayMap[weekdayStr] ?? 1,
    minutesOfDay: hour * 60 + minute,
  };
}

function toMinutes(hhmm: string): number {
  const [h, m] = hhmm.split(":").map((v) => parseInt(v, 10) || 0);
  return h * 60 + m;
}

export function getAvailabilityStatus(input: AvailabilityInput): AvailabilityStatus {
  const days = (input.operating_days || "1,2,3,4,5,6,7")
    .split(",")
    .map((d) => parseInt(d.trim(), 10))
    .filter((d) => d >= 1 && d <= 7);

  const startMin = toMinutes(input.operating_hours_start || "00:00");
  const endMin = toMinutes(input.operating_hours_end || "23:59");

  if (days.length === 7 && startMin === 0 && endMin >= 1439) {
    return { isOpen: true, isAlwaysOpen: true, minutesUntilChange: null };
  }

  const { weekday, minutesOfDay } = nowPartsInTimezone(input.timezone);
  const isActiveDay = days.includes(weekday);
  const withinHours = startMin <= endMin ? minutesOfDay >= startMin && minutesOfDay < endMin : minutesOfDay >= startMin || minutesOfDay < endMin;

  const isOpen = isActiveDay && withinHours;

  let minutesUntilChange: number;
  if (isOpen) {
    minutesUntilChange = endMin > minutesOfDay ? endMin - minutesOfDay : 1440 - minutesOfDay + endMin;
  } else {
    // find minutes until the next start boundary, scanning forward up to 7 days
    minutesUntilChange = 0;
    for (let offset = 0; offset <= 7; offset++) {
      const dayToCheck = ((weekday - 1 + offset) % 7) + 1;
      if (!days.includes(dayToCheck)) continue;
      const candidateStart = offset === 0 ? startMin : startMin;
      const minutesFromNow = offset * 1440 + candidateStart - minutesOfDay;
      if (minutesFromNow > 0) {
        minutesUntilChange = minutesFromNow;
        break;
      }
    }
  }

  return { isOpen, isAlwaysOpen: false, minutesUntilChange };
}

export function formatMinutes(totalMinutes: number, locale: string): string {
  const hours = Math.floor(totalMinutes / 60);
  const minutes = totalMinutes % 60;
  if (locale === "fa") {
    return hours > 0 ? `${hours} ساعت و ${minutes} دقیقه` : `${minutes} دقیقه`;
  }
  if (locale === "zh") {
    return hours > 0 ? `${hours} 小时 ${minutes} 分钟` : `${minutes} 分钟`;
  }
  return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
}
