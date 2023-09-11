# Backend

## Local Development Setup

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/) [No need to install separately on Windows]
- [Golang](https://golang.org/doc/install)
- [Postman](https://www.postman.com/downloads/)
- [MySQL Workbench](https://dev.mysql.com/downloads/workbench/) [Optional - For viewing database]
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) [Optional - For deploying to AWS]


### Setup

1. Clone the repository
2. Open a terminal in the root directory of the repository and run `go mod tidy` to install all dependencies
3. Run `make local` command to start building the docker image and run the docker container [This will take a while for the first time, also note that having Docker and SAM-CLI installed is a prerequisite for this step]
4. Once the container is up and running, you can access the API's with postman

## Deployment to AWS

1. Firstly, install AWS CLI and configure it with your keys, and make a profile named as `lambda-ci` which is used specifically for deploying the application to AWS
2. Run `make deploy` command to deploy the application to AWS [Having SAM-CLI installed is a prerequisite for this step]
