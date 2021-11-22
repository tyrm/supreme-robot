Random rnd = new Random()

pipeline {
  environment {
    networkName = 'network-${env.BUILD_TAG}'
    pgContainerName = 'postgres-${env.BUILD_TAG}'
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
          gitDescribe = sh(returnStdout: true, script: 'git describe --tag').trim()
          writeFile file: "./version/version.go", text: """package version

// Version of the application
const Version = "${gitDescribe}"

          """
          sh "docker network create ${networkName}"
        }
      }
    }

    stage('Start Postgres'){
      steps{
        script{
          retry(4) {
            newPort = rnd.nextInt(9999) + 30000
            echo 'Trying to start postgres on port ${newPort}'
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
          args '-e GOCACHE=/gocache -e HOME=${WORKSPACE} -v /var/lib/jenkins/gocache:/gocache '
        }
      }
      steps {
        script {
          withCredentials([string(credentialsId: 'codecov-tyrm-supreme-robot', variable: 'CODECOV_TOKEN')]) {
            sh """#!/bin/bash
            go get -t -v ./...
            go test -race -coverprofile=coverage.txt -covermode=atomic ./...
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
      sh "docker rm --force postgres-${BUILD_TAG}"
      sh "docker network ${networkName}"
    }
  }

}
