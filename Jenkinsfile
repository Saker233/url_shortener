pipeline {
    agent any

    tools {
        go 'go1.24.3'
    }
    triggers {
        pollSCM('H/1 * * * *')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o url-shortener .'
            }
        }

        stage('Docker build') {
            steps {
                sh 'Docker build -t saker233/url-shortener:latest .'
            }
        }
        stage('Docker Push') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'docker-cred',
                        usernameVariable: 'DOCKER_USERNAME',
                        passwordVariable: 'DOCKER_PASSWORD'
                    )
                ]) {
                    sh '''
                        echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
                        docker push saker233/url-shortener:latest
                        docker logout
                    '''
                }
            }
        }
    }
}