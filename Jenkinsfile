pipeline {
  environment {
    registry = "tyrm/supreme-robot-be"
    registryCredential = 'docker-io-tyrm'
    dockerImage = ''
    gitDescribe = ''
  }

  agent any;

  stages {

    stage('Setup') {
      steps {
        script {
          gitDescribe = sh(returnStdout: true, script: 'git describe --tag').trim()
          writeFile file: "./version/version.go", text: """package version

// Version of the application
const Version = "${gitDescribe}"

          """
          sh "mkdir -p ${WORKSPACE}/embedded-postgres-go"
          sh "chmod 777 ${WORKSPACE}/embedded-postgres-go"
        }
      }
    }
    lock('port_21542') {

      stage('Setup Test'){
        steps{
          script{
            sh """docker run -d \
                    --name postgres-${BUILD_TAG} \
                    -e POSTGRES_PASSWORD=mysecretpassword \
                    postgres"""
          }
        }
      }

      stage('Test') {
        agent {
          docker {
            image 'golang:1.17'
            args '-e GOCACHE=/gocache -e HOME=${WORKSPACE} -v /var/lib/jenkins/gocache:/gocache '
          }
        }
        steps {
          script {
            sh "go get -t -v ./..."
            sh "go test -race -coverprofile=coverage.txt -covermode=atomic ./..."

            withCredentials([string(credentialsId: 'codecov-tyrm-supreme-robot', variable: 'CODECOV_TOKEN')]) {
              sh """#!/bin/bash
              bash <(curl -s https://codecov.io/bash)
              """
            }
          }
        }
      }

      stage('Teardown Test'){
        steps{
          script{
            sh docker rm --force postgres-${BUILD_TAG}"
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
}
