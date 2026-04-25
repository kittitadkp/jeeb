def REPO_URL = 'https://github.com/kittitadkp/jeeb.git'
def CREDS_ID  = 'github-creds'

def pipelines = [
    [name: 'jeeb-backend',  path: 'pipelines/backend/Jenkinsfile'],
    [name: 'jeeb-frontend', path: 'pipelines/frontend/Jenkinsfile'],
]

pipelines.each { p ->
    pipelineJob(p.name) {
        description("Pipeline for ${p.name}")

        triggers {
            scm('H/5 * * * *')   // Poll SCM every 5 minutes
        }

        definition {
            cpsScm {
                scm {
                    git {
                        remote {
                            url(REPO_URL)
                            credentials(CREDS_ID)
                        }
                        branches('*/main')
                        extensions {
                            cleanBeforeCheckout()
                        }
                    }
                }
                scriptPath(p.path)
                lightweight(true)
            }
        }
    }
}
