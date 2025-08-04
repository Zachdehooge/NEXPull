pipeline {
    agent any

    stages {
        stage('Run Go Tests') {
            steps {
                dir('.') {
                    sh 'go test ./nexpull_test.go'
                }
            }
        }

        stage('Push to GitHub') {
            steps {
                dir('.') {
                    // Push changes from local repo to GitHub
                    sh 'git push origin main'
                }
            }
        }
    }
}

