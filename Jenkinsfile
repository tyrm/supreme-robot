pipeline {
  environment {
    registry = "tyrm/supreme-robot-be"
    registryCredential = 'docker-io-tyrm'
    dockerImage = ''
  }

  agent any;

  stages {

    stage('Test') {
      agent {
        docker {
          image 'golang:1.17'
          args '-e GOCACHE=${WORKSPACE}'
        }
      }
      steps {
        script {
          sh "go get -t -v ./..."
          sh "go test -race -coverprofile=coverage.txt -covermode=atomic ./..."

          withCredentials([string(credentialsId: 'codecov-tyrm-supreme-robot', variable: 'CODECOV_TOKEN')]) {
            sh "bash <(curl -s https://codecov.io/bash)"
          }
        }
      }
    }

    stage('Building our image') {
      steps {
        script {
          dockerImage = docker.build registry + ":$BUILD_NUMBER"
        }
      }
    }

    stage('Deploy image') {
      steps {
        script {
          if (env.TAG_NAME) {
            docker.withRegistry('', registryCredential) {
              dockerImage.push(env.TAG_NAME)
            }
          } else {
            docker.withRegistry('', registryCredential) {
              dockerImage.push(env.BRANCH_NAME)
            }
          }
        }
      }
    }

    stage('Cleaning up') {
      steps {
        sh "docker rmi $registry:$BUILD_NUMBER"
      }
    }

  }
}
