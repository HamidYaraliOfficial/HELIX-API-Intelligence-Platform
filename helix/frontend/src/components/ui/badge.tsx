import { cn } from "@/lib/cn";

type BadgeTone = "neutral" | "success" | "warning" | "danger" | "accent";

const toneClasses: Record<BadgeTone, string> = {
  neutral: "bg-bg-elevated text-fg-muted border-border",
  success: "bg-[var(--color-success-bg)] text-success border-[var(--color-success)]",
  warning: "bg-[var(--color-warning-bg)] text-warning border-[var(--color-warning)]",
  danger: "bg-[var(--color-danger-bg)] text-danger border-[var(--color-danger)]",
  accent: "bg-[var(--color-accent-bg)] text-accent border-[var(--color-accent)]",
};

export function Badge({
  children,
  tone = "neutral",
  className,
}: {
  children: React.ReactNode;
  tone?: BadgeTone;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-win border px-2 py-0.5 text-xs font-medium whitespace-nowrap",
        toneClasses[tone],
        className
      )}
    >
      {children}
    </span>
  );
}

export function statusTone(status: string): BadgeTone {
  switch (status) {
    case "healthy":
      return "success";
    case "degraded":
      return "warning";
    case "down":
      return "danger";
    default:
      return "neutral";
  }
}

export function lifecycleTone(state: string): BadgeTone {
  switch (state) {
    case "active":
      return "success";
    case "deprecated":
    case "sunset":
      return "warning";
    case "retired":
      return "danger";
    default:
      return "neutral";
  }
}
