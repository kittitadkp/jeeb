# Plan: Remove GitHub pull from jeeb-seed + clean up stale `seed` job

## Context

Jenkins currently has **two** seed jobs (`seed` + `jeeb-seed`) — only `jeeb-seed` should exist.
`jeeb-seed` is defined in CasC, so it's the authoritative one.
The stale `seed` job is a leftover (unknown origin) that should be removed.

Additionally, `jeeb-seed` currently pulls `seed.groovy` from GitHub (`jeeb-jenkins.git`) via SCM
and runs it as a JobDSL script. This GitHub dependency is unnecessary — the script can be inlined
directly in the CasC job definition.

**Goal:**
- Keep `jeeb-seed` but change it to run an **inline** JobDSL script (no GitHub pull)
- Have that inline script clean up the stale `seed` job on first run
- Add app repo URLs to `values.yaml` so the inline script can reference them via Helm

---

## Files to change

| File | Change |
|------|--------|
| `k8s/charts/jeeb-infra/values.yaml` | Add per-repo git URLs under `jenkins:` |
| `k8s/charts/jeeb-infra/templates/jenkins/configmap-casc.yaml` | Replace jeeb-seed SCM+JobDSL step with inline script; remove `scm` + `triggers` blocks |

---

## Step 1 — Add repo URLs to values.yaml

Under `jenkins:` add:

```yaml
jenkins:
  ...
  repos:
    backend: ""
    frontend: ""
    learningBackend: ""
    learningFrontend: ""
```

---

## Step 2 — Rewrite jeeb-seed in configmap-casc.yaml

Replace the current `jobs:` section (lines 68–92) with a `freeStyleJob('jeeb-seed')` that:
- Has **no `scm` block** (no GitHub pull)
- Has **no `triggers` block** (no polling)
- Has a `steps { jobDsl { ... } }` block with the pipeline definitions **inline** using `scriptText`

```yaml
    jobs:
      - script: |
          freeStyleJob('jeeb-seed') {
            description('Seed job - creates all jeeb pipelines (inline, no SCM pull)')
            steps {
              jobDsl {
                scriptText('''
                  import jenkins.model.Jenkins

                  // Remove stale seed job if it exists
                  def stale = Jenkins.instance.getItem('seed')
                  if (stale) { stale.delete() }

                  folder('dev') {
                    description('DEV environment pipelines')
                  }

                  def pipelines = [
                    [name: 'jeeb-backend',           path: 'pipelines/backend/Jenkinsfile',           gitUrl: '{{ .Values.jenkins.repos.backend }}'],
                    [name: 'jeeb-frontend',          path: 'pipelines/frontend/Jenkinsfile',          gitUrl: '{{ .Values.jenkins.repos.frontend }}'],
                    [name: 'jeeb-learning-backend',  path: 'pipelines/learning-backend/Jenkinsfile',  gitUrl: '{{ .Values.jenkins.repos.learningBackend }}'],
                    [name: 'jeeb-learning-frontend', path: 'pipelines/learning-frontend/Jenkinsfile', gitUrl: '{{ .Values.jenkins.repos.learningFrontend }}'],
                  ]

                  pipelines.each { p ->
                    pipelineJob("dev/${p.name}") {
                      description("${p.name} pipeline - dev")
                      parameters {
                        stringParam('BRANCH', 'main', 'Branch to build')
                        stringParam('GIT_URL', p.gitUrl, 'Repository URL')
                        stringParam('HELM_RELEASE', 'jeeb-dev', 'Helm release name')
                        stringParam('DEPLOY_NAMESPACE', 'jeeb-dev', 'Deploy namespace')
                        booleanParam('SKIP_SONAR', true, 'Skip Test & SonarQube stage')
                      }
                      definition {
                        cpsScm {
                          scm {
                            git {
                              remote {
                                url('{{ .Values.jenkins.jenkinsRepo }}')
                                credentials('github-creds')
                              }
                              branches('*/main')
                              extensions { cleanBeforeCheckout() }
                            }
                          }
                          scriptPath(p.path)
                          lightweight(true)
                        }
                      }
                    }
                  }
                ''')
                sandbox(false)
                ignoreExisting(false)
              }
            }
          }
```

> **Helm note:** `{{ .Values.jenkins.repos.backend }}` etc. inside `scriptText` are Helm template expressions — they get substituted at `helm template` time, before the ConfigMap reaches Jenkins. Groovy's own `${p.name}` interpolation is safe since Helm only processes `{{ }}`.

---

## Steps

- [x] Add `repos:` block to `k8s/charts/jeeb-infra/values.yaml` under `jenkins:`
- [x] Replace `jobs:` section in `configmap-casc.yaml` (lines 68–92) with the new inline version
- [x] Verify Helm rendering: `helm template jeeb-infra k8s/charts/jeeb-infra` produces valid YAML

---

## Verification

```bash
# Render and inspect
helm template jeeb-infra k8s/charts/jeeb-infra | grep -A 120 "jenkins-casc"

# Apply + restart Jenkins
kubectl rollout restart deployment/jenkins -n jeeb-infra

# In Jenkins UI (localhost:30082):
# 1. Run jeeb-seed manually (Build Now)
# 2. Confirm `seed` job is GONE
# 3. Confirm dev/jeeb-backend, dev/jeeb-frontend, dev/jeeb-learning-backend,
#    dev/jeeb-learning-frontend exist
```
