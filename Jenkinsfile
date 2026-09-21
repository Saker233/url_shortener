pipeline {
    agent any

    tools {
        go 'golang'
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
                sh 'docker build -t saker233/url-shortener:latest .'
            }
        }
        stage('Docker Push') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: '1e647c10-3771-419e-894c-5679a19f5f59',
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