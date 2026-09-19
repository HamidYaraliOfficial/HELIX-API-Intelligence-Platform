import { cn } from "@/lib/cn";

export function Input(props: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className={cn(
        "h-10 w-full rounded-win border border-border bg-bg-elevated px-3 text-sm text-fg placeholder:text-fg-muted",
        "focus:border-accent focus:outline-none",
        props.className
      )}
    />
  );
}

export function Textarea(props: React.TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      {...props}
      className={cn(
        "w-full rounded-win border border-border bg-bg-elevated px-3 py-2 text-sm text-fg placeholder:text-fg-muted",
        "focus:border-accent focus:outline-none",
        props.className
      )}
    />
  );
}

export function Select({
  children,
  className,
  ...props
}: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      {...props}
      className={cn(
        "h-10 w-full rounded-win border border-border bg-bg-elevated px-3 text-sm text-fg",
        "focus:border-accent focus:outline-none",
        className
      )}
    >
      {children}
    </select>
  );
}

export function Label({ children, htmlFor }: { children: React.ReactNode; htmlFor?: string }) {
  return (
    <label htmlFor={htmlFor} className="mb-1.5 block text-sm font-medium text-fg">
      {children}
    </label>
  );
}

export function FieldGroup({ children }: { children: React.ReactNode }) {
  return <div className="space-y-1.5">{children}</div>;
}
