import Keycloak from "keycloak-js";

interface KeycloakConfig {
  url: string;
  realm: string;
  clientId: string;
}

export function createKeycloak(config: KeycloakConfig) {
  const keycloak = new Keycloak(config);
  let initPromise: Promise<boolean> | null = null;

  function initKeycloak(): Promise<boolean> {
    if (initPromise !== null) return initPromise;
    initPromise = keycloak
      .init({ onLoad: "login-required", checkLoginIframe: false, pkceMethod: "S256" })
      .catch((err) => {
        initPromise = null;
        throw err;
      });
    return initPromise;
  }

  return { keycloak, initKeycloak };
}
