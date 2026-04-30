import { C, T } from "@/lib/design";

export function LoadingState({ message = "Loading…" }: { message?: string }) {
  return (
    <div style={{ display: "flex", alignItems: "center", justifyContent: "center", padding: 64, color: C.text2, fontSize: T.sm }}>
      {message}
    </div>
  );
}

export function ErrorState({ message = "Something went wrong" }: { message?: string }) {
  return (
    <div style={{ display: "flex", alignItems: "center", justifyContent: "center", padding: 64, color: C.danger, fontSize: T.sm }}>
      {message}
    </div>
  );
}

export function EmptyState({ message = "Nothing here yet" }: { message?: string }) {
  return (
    <div style={{ display: "flex", alignItems: "center", justifyContent: "center", padding: 64, color: C.text3, fontSize: T.sm }}>
      {message}
    </div>
  );
}
