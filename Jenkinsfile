pipeline {
    agent any
    environment {
        DOCKER_REGISTRY = 'vikhang21'
        DOCKER_IMAGE = 'fly_journey'
        DOCKER_TAG = '1.0.0'
        DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'
        COMPOSE_PROJECT_NAME = 'flyjourney'
    }
    stages {
        stage('Setup Docker Permissions') {
            steps {
                script {
                    sh '''
                        sudo usermod -aG docker jenkins || true
                        sudo chmod 666 /var/run/docker.sock || true
                    '''
                }
            }
        }
        stage('Checkout Code') {
            steps {
                git branch: 'main', url: 'https://github.com/Khangvn20/FlyJourney_Backend.git'
            }
        }
        stage('Build Docker Image') {
            steps {
                script {
                    sh "docker build -t ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG} -f Dockerfile ."
                }
            }
        }
        stage('Run Tests') {
            steps {
                script {
                    sh "docker-compose -f docker-compose.yml up -d"
                    sh "docker-compose exec -T app go test ./... || true" 
                    sh "docker-compose down"
                }
            }
        }
        stage('Push Docker Image') {
            steps {
                script {
                    docker.withRegistry('https://index.docker.io/v1/', DOCKER_CREDENTIALS_ID) {
                        sh "docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:${DOCKER_TAG}"
                    }
                }
            }
        }
        stage('Deploy') {
            steps {
                script {
                    sh "docker-compose -f docker-compose.yml down || true"
                    sh "docker-compose -f docker-compose.yml up -d --build"
                }
            }
        }
    }
    post {
        always {
            sh "docker-compose -f docker-compose.yml down || true"
            cleanWs()
        }
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}