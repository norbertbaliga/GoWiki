pipeline {
    agent { label 'agent1' }
    stages {
        stage('build') {
            steps {
                sh 'go build wiki.go'
            }
        }
    }
}