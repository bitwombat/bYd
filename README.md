![Build Status](https://github.com/eddie023/bYd/actions/workflows/main.yml/badge.svg?branch=main)

# Bootstrap Your Dream (bYd)
A serverless starter kit built with Golang, AWS, Postgres, and Terraform. This project provides a cost-effective, enterprise-grade foundation to launch a new project into production. It draws heavily from real-world experience working on production systems.

## High Level System Design Considerations:
1. **Clean Code Architecture**: The codebase follows [Clean Code Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) promoting separation of concerns which makes codebase more maintainable, testable and scalable.

1. **OpenAPI 3.1 standard**: All API endpoints are defined in the 'openapi.yaml'. API types are autogenerated using [oapi-codegen](github.com/deepmap/oapi-codegen/cmd/oapi-codegen) package.

1. **High test coverage and simplified testing**:  KWe prioritize simplicity by minimizing reliance on third-party tools for generating mocks or stubs. Tests focus on the public methods of a package by placing them in a `_test` package, ensuring emphasis on the external behavior and API rather than internal implementation details.

1. **Dockerized Environment**: This codebase is fully Dockerized for seamless local development and effortless production deployment. An AWS Lambda deployment example is provided as Serverless Application so that we can keep the infrastructure costs low until we see some customers. However, the codebase can easily be updated to be deployed to AWS Fargate, AWS ECS or any other Kubernetes services.

1. **CI/CD via GitHub Actions**:  This repository integrates a fully automated Continuous Integration/Continuous Deployment (CI/CD) pipeline using GitHub Actions. It ensures that code is automatically tested, linted and built upon every commit or pull request, improving development speed and code quality. 

1. **Infrastructure as Code(IaC) with Terraform**: The repository leverages Terraform to define and manage cloud infrastructure as code, ensuring a reproducible and scalable setup for AWS resources.

## Running The Project Locally

### Initial Database Setup

Start the API server along with the Postgres database by running: 
```
make run 
```

Apply database migrations using: 
``` 
run `make migrate-up DB_CONNECTION_URI="postgres://root:postgres@localhost:5432/postgres?sslmode=disable"
```

Connect to Postgres via `psql` 
```
psql --host localhost --port 5432 --user root --db postgres` 
```
Then, run the seed script by copying and executing: 
and run your seed script by copy pasting `./migrations/seeds/insert_fakes.sql`

Verify that the service running as expected by using 
```
make get-posts
``` 

## Deploying to AWS Using Terraform 

### Prerequisites:
1. **AWS Account**: Ensure you have an AWS account with appropriate access permissions to provision resources (EC2, RDS, etc.).
2. **Terraform**: Install [Terraform](https://developer.hashicorp.com/terraform/install?product_intent=terraform) on your local machine.
3. **IAM Role/Access**: Ensure you have sufficient AWS IAM permission to create resources such as Cognito, API Gateway, S3, Lambda etc.     

### Basic Steps to Deploy
1. **Set Up S3 and DynamoDB for Terraform State Management**: To enable remote state management and state locking, we will configure Terraform to use AWS S3 and DynamoDB as the backend.

Move to `./infra/terraform/remote-state` directory and run `terraform init` to initialize a separete terraform state for remote s3 backend. Run `terraform plan` verify the details and apply using `terraform apply` 

1. **Initialize Terraform**: 

    - In `./infra/terraform` directory, run: 
    ```
    terraform init 
    ```
    
    - To preview the changes that Terraform will make to your AWS environment, run:
    ```
    terraform plan
    ``` 
    This step allows you to review the resources that will be created, updated, or destroyed before making any actual changes. 

    - To apply the changes and deploy the infrastructure to AWS, run:
    ```
    terraform apply
    ```

1. **Verify API Deployment**:
    - **Get the API Gateway Endpoint**:
        Use the following command to retrieve the API Gateway endpoint from Terraform outputs 
        ```
        terraform output -raw api_gateway_endpoint`
        ```
        Alternatively, you can run a curl request to check if the endpoint is up and running. The following command will send a request to your API Gateway:
        ```
        curl $(terraform output -raw api_gateway_endpoint)/v1/posts
        ```
        If you see a response with {"message": "Unauthorized"}, this is normal since AWS Cognito is being used for authentication and we haven't passed a Authorization header. 

    - **Create New User in AWS Cognito and Authenticate**: 
        - Create a new user in your AWS Cognito User Pool (via the AWS Console or CLI). 
        - After creating the user, obtain their ID token for authentication.
        - Use the following curl command to make an authorized request to the API with the correct ID token:
        ```
        curl -X GET $(terraform output -raw api_gateway_endopint)/v1/posts -H "Authorization: Bearer <YOUR_ID_TOKEN>"
        ```
        You should receive a valid response from the API when the correct ID token is provided.
