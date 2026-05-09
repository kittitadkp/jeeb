import Keycloak from "keycloak-js";

import { appConfig } from "./app-config";

const keycloak = new Keycloak({
  url: appConfig.keycloakUrl,
  realm: appConfig.keycloakRealm,
  clientId: appConfig.keycloakClientId,
});

let initPromise: Promise<boolean> | null = null;

export function initKeycloak(): Promise<boolean> {
  if (initPromise !== null) return initPromise;
  initPromise = keycloak
    .init({ onLoad: "login-required", checkLoginIframe: false, pkceMethod: "S256" })
    .catch((err) => {
      initPromise = null;
      throw err;
    });
  return initPromise;
}

export default keycloak;
