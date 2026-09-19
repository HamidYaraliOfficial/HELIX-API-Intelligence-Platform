"use client";

import { useState, useRef, useEffect } from "react";
import { Palette, Check } from "lucide-react";
import { useTheme, THEMES, type Theme } from "@/lib/theme";
import { useI18n } from "@/lib/i18n";
import { cn } from "@/lib/cn";

const swatch: Record<Theme, string> = {
  light: "#0067c0",
  dark: "#60cdff",
  amoled: "#4cc2ff",
  "windows-blue": "#3aa0ff",
  "windows-red": "#ff5f52",
};

export function ThemeSwitcher() {
  const { theme, setTheme } = useTheme();
  const { t } = useI18n();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function onClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, []);

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-label={t("theme.label")}
        className="flex h-9 w-9 items-center justify-center rounded-win border border-border bg-bg-elevated text-fg hover:bg-panel"
      >
        <Palette className="h-4 w-4" />
      </button>
      {open && (
        <div className="absolute end-0 z-50 mt-2 w-52 rounded-win-lg border border-border bg-panel p-1.5 shadow-win-lg">
          {THEMES.map((themeOption) => (
            <button
              key={themeOption}
              type="button"
              onClick={() => {
                setTheme(themeOption);
                setOpen(false);
              }}
              className={cn(
                "flex w-full items-center gap-2.5 rounded-win px-2.5 py-2 text-sm text-fg hover:bg-bg-elevated"
              )}
            >
              <span
                className="h-3.5 w-3.5 flex-shrink-0 rounded-full border border-border"
                style={{ backgroundColor: swatch[themeOption] }}
              />
              <span className="flex-1 text-start">{t(`theme.${themeOption}`)}</span>
              {theme === themeOption && <Check className="h-3.5 w-3.5 text-accent" />}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
