"use client";

import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

import en from "@/locales/en.json";
import fa from "@/locales/fa.json";
import zh from "@/locales/zh.json";

export type Locale = "en" | "fa" | "zh";

const dictionaries: Record<Locale, Record<string, string>> = { en, fa, zh };

export const LOCALES: { code: Locale; dir: "ltr" | "rtl" }[] = [
  { code: "en", dir: "ltr" },
  { code: "fa", dir: "rtl" },
  { code: "zh", dir: "ltr" },
];

function directionFor(locale: Locale): "ltr" | "rtl" {
  return locale === "fa" ? "rtl" : "ltr";
}

interface I18nContextValue {
  locale: Locale;
  dir: "ltr" | "rtl";
  setLocale: (locale: Locale) => void;
  t: (key: string, vars?: Record<string, string | number>) => string;
}

const I18nContext = createContext<I18nContextValue | null>(null);

const STORAGE_KEY = "helix.locale";

export function I18nProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>("en");
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    const stored = typeof window !== "undefined" ? window.localStorage.getItem(STORAGE_KEY) : null;
    if (stored === "en" || stored === "fa" || stored === "zh") {
      setLocaleState(stored);
    }
    setHydrated(true);
  }, []);

  useEffect(() => {
    if (!hydrated) return;
    document.documentElement.lang = locale;
    document.documentElement.dir = directionFor(locale);
    window.localStorage.setItem(STORAGE_KEY, locale);
  }, [locale, hydrated]);

  const setLocale = useCallback((next: Locale) => setLocaleState(next), []);

  const t = useCallback(
    (key: string, vars?: Record<string, string | number>) => {
      const dict = dictionaries[locale] ?? dictionaries.en;
      let str = dict[key] ?? dictionaries.en[key] ?? key;
      if (vars) {
        for (const [k, v] of Object.entries(vars)) {
          str = str.replace(`{${k}}`, String(v));
        }
      }
      return str;
    },
    [locale]
  );

  const value = useMemo(
    () => ({ locale, dir: directionFor(locale), setLocale, t }),
    [locale, setLocale, t]
  );

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n() {
  const ctx = useContext(I18nContext);
  if (!ctx) throw new Error("useI18n must be used within an I18nProvider");
  return ctx;
}
