# Simple WikiMedia API endpoint to search for people and output short description as JSON

## Features

- Docker file that hosts on PORT 8080
- All REST APIs (GET)

## Getting started: Start server

```bash or zsh
go get
go build
./wiki-names &
```

Run tests

```bash or zsh
go test -v ./...
```

### While running locally, you can browse to search end point:

http://localhost:8080/search/Yoshua_Bengio

**Response:**
FYI, the URL names are case sensitive, so be careful when searching to use the correct uppercase lowercase letters. Maybe it would be nice to add a API endpoint that returns suggested list of search terms based on what the user types in. This would allow UI developers to add a auto-suggest box to improve usability. We could create it with a regex patterns, soundex, double metaphone or n-gram matching algorithm.

I did not add any "fuzzy" matching to the API request string, because I'd like to keep the inputs and outputs deterministic over a long period of time. If fuzzy matching was part of the solution, I'd push that to the client-side team :)

```json
{
  "short_description": "....."
}
```

## Run as Docker Container:

```bash or zsh
docker build -t wiki_names .
docker run -p 8080:8081 -d --name wiki_service wiki_names
#Open in Browser window
open http://localhost:8081/search/Yoshua_Bengio

#You may need the container IP address, in some scenarios on Windows desktops
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wiki_service

#Clean up image and container from your local Docker Desktop cache
docker rmi wiki_service
docker rm wiki_names
```

## API

There are the 4 end points accessible with the API. Not the Swagger is a WIP, not enough time to build it out
[GIN-debug] GET /search/:name --> wiki-names/controllers.GetContentSummary (4 handlers)
[GIN-debug] GET /extract/:name --> wiki-names/controllers.GetExtract (4 handlers)
[GIN-debug] GET /extract/:name/:locale --> wiki-names/controllers.GetExtract (4 handlers)
[GIN-debug] GET /swagger/\*any --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (4 handlers)

## Deployment to Production and CI/CD pipeline

Jenkins build script. Leaving out the final deployment because there are too many dependencies on how DEV, UAT and PROD manages their infrastructures.

```
pipeline {
    agent any
    tools {
        go 'go1.17'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }
    stages {
        stage('Pre Test') {
            steps {
                echo 'Installing dependencies'
                sh 'go version'
                sh 'go get -u golang.org/x/lint/golint'
            }
        }

        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'go build'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet .'
                    echo 'Running linting'
                    sh 'golint .'
                    echo 'Running test'
                    sh 'go test -v ./...'
                }
            }
        }
    }
    post {
        always {
            emailext body: "${currentBuild.currentResult}: Job ${env.JOB_NAME} build ${env.BUILD_NUMBER}\n More info at: ${env.BUILD_URL}",
                recipientProviders: [[$class: 'DevelopersRecipientProvider'], [$class: 'RequesterRecipientProvider']],
                to: "${params.RECIPIENTS}",
                subject: "Jenkins Build ${currentBuild.currentResult}: Job ${env.JOB_NAME}"

        }
    }
}
```

## Scaling, fallback, observability and performance

My preferred would be to have local Helm files in the current Repo, and run a Helm deploy CLI command to deploy to EKS on AWS.
The Helm files would contain auto-scaling based on CPU load on the PODs, an example of the Helm script would be, below. The way the PODs scale is based on CPU load, but we can also use network traffic and maybe look at a per-customer rate limit to stop some users from overloading and exhausting the Kubernetes cluster (like BOT attacks).

```
{{- if .Values.hpa.enabled -}}
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: {{ template "customerapi.fullname" . }}
spec:
  maxReplicas: {{ .Values.hpa.maxReplicas }}
  minReplicas: {{ .Values.hpa.minReplicas }}
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ template "customerapi.fullname" . }}
  targetCPUUtilizationPercentage: {{ .Values.hpa.averageCpuUtilization }}
{{- end }}
```

To that end, I'd configure Nginx to be the Load Balancer to proxy the inbound HTTP request to the PODS, and keep metrics on which custom makes each request. This would allow us to create a rate limit or fast-track configuration, where higher paying customers get less restrictions.
Depending on the number of users and network traffic, the number of PODs could be small, starting out with 3 and going up to having individual namespaces or clusters for larger clients.

The logs from each service would be directed to AWS CloudWatch, so all observability is within one main index. And add SRE alerts to trigger depending on some log filters and metrics.

If security and over use uis a worry, I'd look at adding a WAF firewall and maybe putting a commercial CDN in front of our Nginx load balancer. Previously, I used CloudFlare but now AWS CloudFront. My knowledge on that part is sketchy, it was mostly phone calls to CloudFlare support when they needed to change filtering rules for BOT attacks in Adidas.

## Network reliability and availability

We are using the Wikimedia Endpoints here and there is a concern that we may overload their network with requests, that is why I added a in memory (or Redis caching) on our API results. THis is a question that should be brought to the client and out PO to see if timeliness of data feeds is a issue, current code can cache up to two minutes of data per unique URL. Some customers may want more timely information. This could be configured on a per-customer basis, if we add API keys and track key usage with the API code.

One option is to add better fallback logic if the network between this service and the WikiMedia APIs is unreliable. We could add a Retry with exponential time outs, if the Wikimedia endpoints are not responding. We'd need to monitor this change in logic, because there is a danger our API would accelerate the Wikimedia API delays, by flooding their service with multiple requests, while there service is not responding or enduring high demand. Again, this needs to be discussed with Tech Leads and POs to decide the probability of slow (or poor) network responses from the main Wikimedia APIs. My assumption is that this API and network infrastructure is stable and scalable for our needs with the simple API.

## Learning Outcomes

I switched to Golang for this project, because I know this is the direction your team wants to go towards. My preference would have been NodeJS with Typescript, because it would have been faster to develop the Unit tests and the JSON parsing would have had a richer ecosystem of 3rd party libraries.

I didn't have enough time to complete:

1. Add full unit tests, closer to 100% coverage
2. Add the Helm files for Kubernetes deployment
3. Test out the Jenkins build pipeline script that I show above. I don't have a local Jenkins sandbox, to play with.
4. From a real-world perspective the code needs a bit more time and a few discussions with PO and TechLead on best Business strategy for the more non-functional aspects of the project.

## Multi-lingual Markup and WikiText Parsing Complexity

Overall it was a fun exercise, it allowed me to get a glimpse at the complexity of the WikiText markdown and the variety of formatting you have across the different languages. After chatting with Stefania, my first thought would be to get a better visibility of the types of markup tags are use and frequency of tag patterns, to see which ones are highest priority to parse. IT's probably too difficult to ask editors to agree on a single formatting vocabulary for the tagging, so it's not simple. My solution would be to experiment with ML parsing models to find the distribution of tagging patterns, and tune a ML model to accurately parse the tagging correctly. I'm excited to delve into the Table parsing question to come up with an AI model that extract table data in a structured and flexible way, training a model to adapt to formatting styles and choose the best table parser it can produce.

As an interim solution, I added a /extract/ API endpoint to the solution code. This returns the first two sentences of the Extract text for that query. This is a compromise, and not in the requirements. IT would need sign-off by the PO and TechLead in the team :)
