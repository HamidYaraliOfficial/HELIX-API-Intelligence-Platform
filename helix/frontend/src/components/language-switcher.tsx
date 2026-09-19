"use client";

import { useState, useRef, useEffect } from "react";
import { Languages, Check } from "lucide-react";
import { useI18n, LOCALES } from "@/lib/i18n";
import { cn } from "@/lib/cn";

export function LanguageSwitcher() {
  const { locale, setLocale, t } = useI18n();
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
        aria-label={t("lang.label")}
        className="flex h-9 w-9 items-center justify-center rounded-win border border-border bg-bg-elevated text-fg hover:bg-panel"
      >
        <Languages className="h-4 w-4" />
      </button>
      {open && (
        <div className="absolute end-0 z-50 mt-2 w-40 rounded-win-lg border border-border bg-panel p-1.5 shadow-win-lg">
          {LOCALES.map(({ code }) => (
            <button
              key={code}
              type="button"
              onClick={() => {
                setLocale(code);
                setOpen(false);
              }}
              className={cn("flex w-full items-center gap-2.5 rounded-win px-2.5 py-2 text-sm text-fg hover:bg-bg-elevated")}
            >
              <span className="flex-1 text-start">{t(`lang.${code}`)}</span>
              {locale === code && <Check className="h-3.5 w-3.5 text-accent" />}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
