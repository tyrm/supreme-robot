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
          image 'gradle:6.7-jdk11'
          // Run the container on the node specified at the top-level of the Pipeline, in the same workspace, rather than on a new node entirely:
          reuseNode true
        }
      }
      steps {
        script {
          sh "docker run --rm -v $WORKSPACE:/app --workdir="/app" golang:1.17 go test ./..."
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
