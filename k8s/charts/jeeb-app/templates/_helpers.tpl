{{/*
  jeeb-app.vaultAgentHCL
  ----------------------
  Renders the shared Vault Agent agent.hcl content:
    - vault address block
    - kubernetes auto_auth + file sink
    - template block pointing at the service-specific env template

  Call with a dict:
    "vaultAddr"  — .Values.infra.vaultAddr
    "role"       — .Values.<service>.vault.role
    "envFile"    — .Values.<service>.vault.envFile  (e.g. .env.develop)
    "tplFile"    — computed tpl filename            (e.g. env-develop.tpl)
*/}}
{{- define "jeeb-app.vaultAgentHCL" -}}
vault {
  address = "{{ .vaultAddr }}"
}

auto_auth {
  method "kubernetes" {
    mount_path = "auth/kubernetes"
    config = {
      role = "{{ .role }}"
    }
  }

  sink "file" {
    config = {
      path = "/app/env/.vault-token"
    }
  }
}

template {
  source      = "/vault/config/{{ .tplFile }}"
  destination = "/app/env/{{ .envFile }}"
  perms       = "0640"
  {{- if .restartCmd }}
  exec {
    command = ["/bin/sh", "-c", "{{ .restartCmd }}"]
    timeout = "30s"
  }
  {{- end }}
}
{{- end -}}

{{/*
  jeeb-app.vaultAgentContainer
  ----------------------------
  Renders the vault-agent sidecar container spec.

  Call with a dict:
    "vaultImage" — .Values.<service>.vaultImage
    "vaultAddr"  — .Values.infra.vaultAddr
*/}}
{{- define "jeeb-app.vaultAgentContainer" -}}
- name: vault-agent
  image: {{ .vaultImage }}
  args:
    - agent
    - -config=/vault/config/agent.hcl
  env:
    - name: VAULT_ADDR
      value: {{ .vaultAddr | quote }}
  volumeMounts:
    - name: vault-config
      mountPath: /vault/config
    - name: vault-secrets
      mountPath: /app/env
    - name: tools
      mountPath: /tools
      readOnly: true
  resources:
    requests:
      memory: 64Mi
      cpu: 50m
    limits:
      memory: 128Mi
      cpu: 100m
  securityContext:
    allowPrivilegeEscalation: false
    capabilities:
      add:
        - IPC_LOCK
{{- end -}}

{{/*
  jeeb-app.vaultKubectlInitContainer
  -----------------------------------
  Init container that copies kubectl into the shared /tools volume so
  vault-agent can call "kubectl rollout restart" when secrets change.
*/}}
{{- define "jeeb-app.vaultKubectlInitContainer" -}}
- name: kubectl-installer
  image: bitnami/kubectl:latest
  command: ["cp", "/opt/bitnami/kubectl/bin/kubectl", "/tools/kubectl"]
  volumeMounts:
    - name: tools
      mountPath: /tools
  securityContext:
    allowPrivilegeEscalation: false
{{- end -}}

{{/*
  jeeb-app.vaultVolumes
  ---------------------
  Renders the two volumes needed by the vault-agent sidecar.
  Call with a dict:
    "configMapName" — name of the vault-agent ConfigMap for this service
*/}}
{{- define "jeeb-app.vaultVolumes" -}}
- name: vault-config
  configMap:
    name: {{ .configMapName }}
- name: vault-secrets
  emptyDir: {}
- name: tools
  emptyDir: {}
{{- end -}}
