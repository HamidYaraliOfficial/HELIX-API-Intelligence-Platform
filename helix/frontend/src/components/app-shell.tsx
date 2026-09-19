"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutGrid, Radar, Users2, Hexagon } from "lucide-react";
import { useI18n } from "@/lib/i18n";
import { ThemeSwitcher } from "@/components/theme-switcher";
import { LanguageSwitcher } from "@/components/language-switcher";
import { cn } from "@/lib/cn";

function NavLink({ href, icon: Icon, label, active }: { href: string; icon: any; label: string; active: boolean }) {
  return (
    <Link
      href={href}
      className={cn(
        "flex items-center gap-3 rounded-win px-3 py-2.5 text-sm font-medium transition-colors",
        active ? "bg-[var(--color-accent-bg)] text-accent" : "text-fg-muted hover:bg-bg-elevated hover:text-fg"
      )}
    >
      <Icon className="h-4 w-4 flex-shrink-0" />
      <span>{label}</span>
    </Link>
  );
}

export function AppShell({ children }: { children: React.ReactNode }) {
  const { t } = useI18n();
  const pathname = usePathname();

  const nav = [
    { href: "/", icon: LayoutGrid, label: t("nav.dashboard") },
    { href: "/discovery", icon: Radar, label: t("nav.discovery") },
    { href: "/teams", icon: Users2, label: t("nav.teams") },
  ];

  return (
    <div className="flex h-screen w-full overflow-hidden">
      <aside className="fluent-mica flex w-64 flex-shrink-0 flex-col border-e border-border p-4">
        <div className="mb-6 flex items-center gap-2.5 px-1">
          <div className="flex h-9 w-9 items-center justify-center rounded-win-lg bg-accent text-accent-fg">
            <Hexagon className="h-5 w-5" />
          </div>
          <div>
            <div className="text-base font-semibold leading-tight text-fg">{t("app.name")}</div>
            <div className="text-xs leading-tight text-fg-muted">{t("app.tagline")}</div>
          </div>
        </div>

        <nav className="flex flex-col gap-1">
          {nav.map((item) => (
            <NavLink key={item.href} {...item} active={pathname === item.href} />
          ))}
        </nav>
      </aside>

      <div className="flex flex-1 flex-col overflow-hidden">
        <header className="flex h-14 flex-shrink-0 items-center justify-end gap-2 border-b border-border bg-bg-elevated px-4">
          <ThemeSwitcher />
          <LanguageSwitcher />
        </header>
        <main className="flex-1 overflow-y-auto bg-bg p-6">{children}</main>
      </div>
    </div>
  );
}
