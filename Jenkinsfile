pipeline {
    agent { label 'go' }
    stages {
        stage('build') {
            steps {
                sh 'go build wiki.go'
            }
        }
    }
}
