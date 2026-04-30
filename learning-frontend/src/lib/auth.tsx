import { type ReactNode, useEffect, useMemo, useState } from "react";

import { setAccessToken } from "./api";
import { AuthContext, type AuthUser } from "./auth-context";
import keycloak, { initKeycloak } from "./keycloak";

export function AuthProvider({ children }: Readonly<{ children: ReactNode }>) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [ready, setReady] = useState(false);
  const [initError, setInitError] = useState<string | null>(null);
  const authValue = useMemo(
    () => ({ user, ready, logout: () => keycloak.logout() }),
    [user, ready],
  );

  useEffect(() => {
    let active = true;

    initKeycloak()
      .then((authenticated) => {
        if (!active) return;
        if (authenticated) {
          setAccessToken(keycloak.token ?? "");
          const parsed = keycloak.tokenParsed as Record<string, string> | undefined;
          setUser({
            id: parsed?.sub ?? "",
            email: parsed?.email ?? "",
            name: parsed?.name ?? parsed?.preferred_username ?? "User",
          });
        } else {
          keycloak.login();
          return;
        }
        setReady(true);
      })
      .catch((err: unknown) => {
        if (!active) return;
        const msg = err instanceof Error ? err.message : "Unable to connect to auth server";
        setInitError(msg);
        setReady(true);
      });

    keycloak.onTokenExpired = () => {
      keycloak.updateToken(60)
        .then((refreshed) => { if (refreshed) setAccessToken(keycloak.token ?? ""); })
        .catch(() => keycloak.login());
    };

    return () => { active = false; };
  }, []);

  if (ready && initError) {
    return (
      <div style={{ display: "flex", height: "100vh", flexDirection: "column", alignItems: "center", justifyContent: "center", gap: 16, background: "var(--c-bg)" }}>
        <div style={{ color: "#DC2626", fontWeight: 500 }}>Authentication unavailable</div>
        <div style={{ color: "var(--c-text2)", fontSize: 13, textAlign: "center", maxWidth: 280 }}>{initError}</div>
        <button onClick={() => globalThis.location.reload()} style={{ background: "#0ea5e9", color: "#fff", border: "none", borderRadius: 8, padding: "8px 16px", fontSize: 13, cursor: "pointer" }}>
          Retry
        </button>
      </div>
    );
  }

  return (
    <AuthContext.Provider value={authValue}>
      {ready ? children : (
        <div style={{ display: "flex", height: "100vh", alignItems: "center", justifyContent: "center", background: "var(--c-bg)" }}>
          <div style={{ color: "var(--c-text2)", fontSize: 13 }}>Loading…</div>
        </div>
      )}
    </AuthContext.Provider>
  );
}
