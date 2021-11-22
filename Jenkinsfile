pipeline {
  environment {
    networkName = "network-${env.BUILD_TAG}"
    pgContainerName = "postgres-${env.BUILD_TAG}"
    registry = 'tyrm/supreme-robot-be'
    registryCredential = 'docker-io-tyrm'
    dockerImage = ''
    gitDescribe = ''
  }

  agent any

  stages {

    stage('Setup') {
      steps {
        script {
          echo "creating network ${networkName}"
          sh "docker network create ${networkName}"
          gitDescribe = sh(returnStdout: true, script: 'git describe --tag').trim()
          writeFile file: "./version/version.go", text: """package version

// Version of the application
const Version = "${gitDescribe}"

          """
        }
      }
    }

    stage('Start Postgres'){
      steps{
        script{
          retry(4) {
            echo 'trying to start postgres'
            withCredentials([usernamePassword(credentialsId: 'integration-postgres-test', usernameVariable: 'POSTGRES_USER', passwordVariable: 'POSTGRES_PASSWORD')]) {
              sh """docker run -d \
                      --name ${pgContainerName} \
                      --network ${networkName} \
                      --env POSTGRES_DB=supremerobot \
                      --env POSTGRES_USER=${POSTGRES_USER} \
                      --env POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
                      --pull always \
                      postgres:14"""
            }
          }
        }
      }
    }

    stage('Test') {
      agent {
        docker {
          image 'golang:1.17'
          args '--network ${networkName} -e GOCACHE=/gocache -e HOME=${WORKSPACE} -v /var/lib/jenkins/gocache:/gocache'
        }
      }
      steps {
        script {
          withCredentials(
            [string(credentialsId: 'codecov-tyrm-supreme-robot', variable: 'CODECOV_TOKEN')],
            [usernamePassword(credentialsId: 'integration-postgres-test', usernameVariable: 'POSTGRES_USER', passwordVariable: 'POSTGRES_PASSWORD')]
          ) {
            pgConnectionDSN = "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${pgContainerName}:5432/supremerobot?sslmode=disable"

            sh """#!/bin/bash
            go get -t -v ./...
            TEST_DSN="${pgContainerName}" go test --tags=integration -race -coverprofile=coverage.txt -covermode=atomic ./...
            bash <(curl -s https://codecov.io/bash)
            """
          }
        }
      }
    }

    stage('Upload image') {
      steps {
        script {
          retry(3) {
            timeout(time: 15, unit: 'MINUTES') {
              if (env.TAG_NAME) {
                sh "DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/arm64,linux/amd64 -t ${registry}:${env.TAG_NAME} . --push"
              } else {
                sh "DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/arm64,linux/amd64 -t ${registry}:${env.BRANCH_NAME} . --push"
              }
            }
          }
        }
      }
    }

    stage('Deploy develop') {
      when{
        branch 'develop'
      }
      steps {
        script {
          build job: 'deploy-supreme-robot_develop', wait: false
        }
      }
    }

  }

  post {
    always {
      sh "docker rm --force ${pgContainerName}"
      sh "docker network rm ${env.networkName}"
    }
  }

}
