{{- define "jeeb-learning.vaultAgentHCL" -}}
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

{{- define "jeeb-learning.vaultAgentContainer" -}}
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

{{- define "jeeb-learning.vaultKubectlInitContainer" -}}
- name: kubectl-installer
  image: bitnami/kubectl:latest
  command: ["cp", "/opt/bitnami/kubectl/bin/kubectl", "/tools/kubectl"]
  volumeMounts:
    - name: tools
      mountPath: /tools
  securityContext:
    allowPrivilegeEscalation: false
{{- end -}}

{{- define "jeeb-learning.vaultVolumes" -}}
- name: vault-config
  configMap:
    name: {{ .configMapName }}
- name: vault-secrets
  emptyDir: {}
- name: tools
  emptyDir: {}
{{- end -}}
