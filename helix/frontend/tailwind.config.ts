import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: ["class"],
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        bg: "var(--color-bg)",
        "bg-elevated": "var(--color-bg-elevated)",
        "bg-mica": "var(--color-bg-mica)",
        panel: "var(--color-panel)",
        border: "var(--color-border)",
        fg: "var(--color-fg)",
        "fg-muted": "var(--color-fg-muted)",
        accent: "var(--color-accent)",
        "accent-hover": "var(--color-accent-hover)",
        "accent-fg": "var(--color-accent-fg)",
        success: "var(--color-success)",
        warning: "var(--color-warning)",
        danger: "var(--color-danger)",
      },
      borderRadius: {
        win: "8px",
        "win-lg": "12px",
      },
      fontFamily: {
        sans: [
          "Segoe UI Variable",
          "Segoe UI",
          "Vazirmatn",
          "Noto Sans SC",
          "system-ui",
          "sans-serif",
        ],
      },
      boxShadow: {
        win: "0 2px 8px rgba(0,0,0,0.12), 0 0 1px rgba(0,0,0,0.2)",
        "win-lg": "0 8px 28px rgba(0,0,0,0.18), 0 0 1px rgba(0,0,0,0.24)",
      },
    },
  },
  plugins: [],
};

export default config;
