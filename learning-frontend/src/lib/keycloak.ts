import Keycloak from "keycloak-js";

const keycloak = new Keycloak({
  url: import.meta.env.VITE_KEYCLOAK_URL ?? "http://localhost:30081",
  realm: import.meta.env.VITE_KEYCLOAK_REALM ?? "jeeb",
  clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID ?? "jeeb-app",
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
