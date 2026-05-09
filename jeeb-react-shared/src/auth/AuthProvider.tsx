import { type ReactNode, useEffect, useMemo, useState } from "react";

import type Keycloak from "keycloak-js";

import { AuthContext, type AuthUser } from "./auth-context";

interface AuthProviderProps {
  children: ReactNode;
  keycloak: Keycloak;
  initKeycloak: () => Promise<boolean>;
  onTokenChange?: (token: string) => void;
}

export function AuthProvider({ children, keycloak, initKeycloak, onTokenChange }: AuthProviderProps) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [ready, setReady] = useState(false);
  const [initError, setInitError] = useState<string | null>(null);

  const authValue = useMemo(
    () => ({ user, ready, logout: () => keycloak.logout() }),
    [user, ready, keycloak],
  );

  useEffect(() => {
    let active = true;

    initKeycloak()
      .then((authenticated) => {
        if (!active) return;
        if (authenticated) {
          onTokenChange?.(keycloak.token ?? "");
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
        .then((refreshed) => {
          if (refreshed) onTokenChange?.(keycloak.token ?? "");
        })
        .catch(() => keycloak.login());
    };

    return () => {
      active = false;
    };
  }, []);

  if (ready && initError) {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-4 bg-slate-50 dark:bg-slate-950">
        <div className="font-medium text-red-500">Authentication unavailable</div>
        <div className="max-w-xs text-center text-sm text-slate-400">{initError}</div>
        <button
          onClick={() => globalThis.location.reload()}
          className="rounded-lg bg-blue-600 px-4 py-2 text-sm text-white transition-colors hover:bg-blue-700"
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <AuthContext.Provider value={authValue}>
      {ready ? (
        children
      ) : (
        <div className="flex h-screen items-center justify-center bg-slate-50 dark:bg-slate-950">
          <div className="text-sm text-slate-400">Loading…</div>
        </div>
      )}
    </AuthContext.Provider>
  );
}
