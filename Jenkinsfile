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
          gitDescribe = sh(script: "git describe", returnStdout: true).trim()
        }
      }
    }

    stage('Test') {
      agent {
        docker {
          image 'golang:1.17'
        }
      }
      steps {
        script {
          sh "go test ./..."
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

    stage('Deploy our image') {
      steps {
        script {
          docker.withRegistry( '', registryCredential ) {
            dockerImage.push(gitDescribe)
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
