import { type ReactNode } from "react";

import { AuthProvider as SharedAuthProvider } from "@jeeb/react-shared/auth";

import { setAccessToken } from "./api";
import keycloak, { initKeycloak } from "./keycloak";

export function AuthProvider({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <SharedAuthProvider
      keycloak={keycloak}
      initKeycloak={initKeycloak}
      onTokenChange={setAccessToken}
    >
      {children}
    </SharedAuthProvider>
  );
}
