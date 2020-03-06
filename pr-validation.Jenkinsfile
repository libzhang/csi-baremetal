// load pipelines libraries from https://eos2git.cec.lab.emc.com/ECS/pipelines
loader.loadFrom('pipelines': [common       : 'common',
                              pr_validation: 'infra/pr_validation',])

// run this job
this.runPullRequestValidationJob()

// define this job
void runPullRequestValidationJob() {
    Map<String, Object> args = [
        pullRequestNumber: params.PULL_REQUEST_NUMBER,
        repo             : common.BAREMETAL_CSI_PLUGIN_REPO_SSH,
    ]
    pr_validation.runPullRequestValidationJob(args) {
        String commit ->
            return this.validatePullRequest(commit)
    }
}

// define pr-validation logic
boolean validatePullRequest(String commit) {
    int lintExitCode = 0
    int testExitCode = 0
    int coverageExitCode = 0
    int buildExitCode = 0
    int imageExitCode = 0
    int depExitCode = 0
    int chartsLintExitCode = 0
    int chartsInstallExitCode = 0
    common.node(label: 'ubuntu_build_hosts', time: 180) {
        common.withInfraDevkitContainerKind() {
            try {
                stage('Git Clone') {
                    checkout scm
                }

                stage('Get dependencies') {
                    depExitCode = sh(script: '''
                                        make install-compile-proto
                                        make install-hal
                                        make install-controller-gen
                                        make generate-deepcopy
                                        make dependency
                                     ''', returnStatus: true)
                    if (depExitCode != 0) {
                        currentBuild.result = 'FAILURE'
                        throw new Exception("Get dependencies stage failed, check logs")
                    }
                }

                parallel(
                    'charts': {
                        stage('Create Kind Cluster') {
                            sh('kind create cluster --kubeconfig /root/.kube/config --config config.yaml')
                        }
                        stage('Lint') {
                            chartsLintExitCode = sh(script: 'make lint-charts', returnStatus: true)
                            if (chartsLintExitCode != 0) {
                                currentBuild.result = 'FAILURE'
                                throw new Exception("Helm lint stage failed, check logs")
                            }
                        }
                        stage('Install Charts') {
                            chartsInstallExitCode = sh(script: 'helm install csi ' +
                                '--kubeconfig=/root/.kube/config ./charts/baremetal-csi-plugin', returnStatus: true)
                            if (chartsInstallExitCode != 0) {
                                currentBuild.result = 'FAILURE'
                                throw new Exception("Install charts stage failed, check logs")
                            }
                        }
                    },
                    'plugin': {
                        stage('Lint') {
                            lintExitCode = sh(script: 'make lint', returnStatus: true)
                            if (lintExitCode != 0) {
                                currentBuild.result = 'FAILURE'
                                throw new Exception("Lint stage failed, check logs")
                            }
                        }

                        stage('Build') {
                            buildExitCode = sh(script: 'make build', returnStatus: true)
                            if (buildExitCode != 0) {
                                currentBuild.result = 'FAILURE'
                                throw new Exception("Build stage failed, check logs")
                            }
                        }

                        stage('Test and Coverage') {
                            testExitCode = sh(script: 'make test', returnStatus: true)
                            //split because our make test fails and make coverage isn't invoked during sh()
                            coverageExitCode = sh(script: 'make coverage', returnStatus: true)
                            if ((testExitCode != 0) || (coverageExitCode != 0)) {
                                currentBuild.result = 'FAILURE'
                                throw new Exception("Test and Coverage stage failed, check logs")
                            }
                        }

                    }
                )
                stage('Make image') {
                    imageExitCode = sh(script: 'make image', returnStatus: true)
                    if (imageExitCode != 0) {
                        currentBuild.result = 'FAILURE'
                        throw new Exception("Image stage failed, check logs")
                    }
                }

            }
            finally {
                sh('kind delete cluster')
                // publish in Jenkins test results
                archiveArtifacts('coverage.html')
            }
        }

    }
    // If we got here then nothing failed
    return true // as a mark of successful validation
}

this
