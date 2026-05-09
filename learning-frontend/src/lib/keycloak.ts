import { createKeycloak } from "@jeeb/react-shared/auth";

import { appConfig } from "./app-config";

const { keycloak, initKeycloak } = createKeycloak({
  url: appConfig.keycloakUrl,
  realm: appConfig.keycloakRealm,
  clientId: appConfig.keycloakClientId,
});

export { initKeycloak };
export default keycloak;
