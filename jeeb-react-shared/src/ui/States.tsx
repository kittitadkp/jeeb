import { C, S, T, W } from "../utils/design";

interface EmptyProps {
  emoji: string;
  title: string;
  description?: string;
}

interface ErrorProps {
  message?: string;
  onRetry?: () => void;
}

export function EmptyState({ emoji, title, description }: EmptyProps) {
  return (
    <div style={{ textAlign: "center", padding: `${S[10]}px ${S[6]}px` }}>
      <div style={{ fontSize: 36, marginBottom: S[3] }}>{emoji}</div>
      <div style={{ fontSize: T.md, fontWeight: W.semi, color: C.text, marginBottom: S[1] + 2 }}>{title}</div>
      {description && <div style={{ fontSize: T.base, color: C.text2 }}>{description}</div>}
    </div>
  );
}

export function LoadingState() {
  return (
    <div style={{ padding: `${S[10]}px 0`, textAlign: "center", fontSize: T.base, color: C.text2 }}>
      Loading…
    </div>
  );
}

export function ErrorState({ message, onRetry }: ErrorProps) {
  return (
    <div style={{ padding: `${S[10]}px 0`, textAlign: "center" }}>
      <div style={{ fontSize: T.base, color: C.danger, marginBottom: S[2] }}>
        {message ?? "Something went wrong"}
      </div>
      {onRetry && (
        <button
          onClick={onRetry}
          style={{ fontSize: T.sm, color: C.primary, background: "none", border: "none", cursor: "pointer" }}
        >
          Retry
        </button>
      )}
    </div>
  );
}
