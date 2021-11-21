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
          sh "mkdir ${WORKSPACE}/embedded-postgres-go"
        }
      }
    }

    stage('Test') {
      agent {
        docker {
          image 'golang:1.17'
          args '-e GOCACHE=${WORKSPACE} -v ${WORKSPACE}/embedded-postgres-go:/.embedded-postgres-go'
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
