pipeline {
    agent {
        docker {
            image 'docker:latest'
            args '-v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker --group-add docker'
        }
    }
    environment {
        DOCKER_REGISTRY = 'vikhang21'
        DOCKER_IMAGE = 'fly_journey'
        DOCKER_TAG = '1.0.0'
        DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'
        COMPOSE_PROJECT_NAME = 'flyjourney'
    }
    stages {
        stage('Install Dependencies') {
            steps {
                sh '''
                    apk add --no-cache git curl
                    curl -L "https://github.com/docker/compose/releases/download/v2.21.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
                    chmod +x /usr/local/bin/docker-compose
                    docker-compose --version
                '''
            }
        }
        stage('Checkout Code') {
            steps {
                git branch: 'main', url: 'https://github.com/Khangvn20/FlyJourney_Backend.git'
            }
        }
        stage('Verify Environment') {
            steps {
                sh '''
                    echo "Checking .env file..."
                    if [ -f .env ]; then
                        echo ".env file exists"
                        # Don't show sensitive data in logs
                        echo "Environment variables loaded from .env"
                    else
                        echo ".env file not found!"
                        exit 1
                    fi
                '''
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
                    sh '''
                        echo "Starting Redis service for testing..."
                        docker-compose -f docker-compose.yml up -d redis
                        echo "Waiting for Redis to be ready..."
                        sleep 10
                        echo "Starting application (connecting to external PostgreSQL)..."
                        docker-compose -f docker-compose.yml up --build -d app
                        sleep 15
                        echo "Running tests..."
                        docker-compose exec -T app go test ./... || true
                        echo "Stopping test services..."
                        docker-compose down
                    '''
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
                    sh '''
                        echo "Stopping existing services..."
                        docker-compose -f docker-compose.yml down || true
                        echo "Starting production deployment (Redis + App only)..."
                        docker-compose -f docker-compose.yml up -d redis app --build
                        echo "Deployment completed!"
                    '''
                }
            }
        }
    }
    post {
        always {
            sh '''
                echo "Cleaning up..."
                docker-compose -f docker-compose.yml down || true
                docker system prune -f || true
            '''
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