# Plato API

Plato API is a prompt management platform developed by PayloadOps. This API streamlines the management of AI prompts, projects, organizations, teams, and users through conventional HTTP requests. The platform provides robust tools for developers to manage settings, memberships, and activities seamlessly.

## Table of Contents

- [Plato API](#plato-api)
    - [Table of Contents](#table-of-contents)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Usage](#usage)
        - [Running Locally](#running-locally)
        - [Running with Docker Compose](#running-with-docker-compose)
    - [Configuration](#configuration)
    - [API Documentation](#api-documentation)
    - [Development](#development)
    - [Running Tests](#running-tests)
        - [Running Unit Tests](#running-unit-tests)
        - [Running e2e Tests](#running-e2e-tests)

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go 1.22.2 or later
- Docker
- Docker Compose
- AWS CLI (for setting up AWS services)
- LocalStack (for local AWS service emulation)

## Installation

1. Clone the repository:

```sh
git clone https://github.com/payloadops/plato-api.git
cd plato-api
```

2. Install the Go dependencies:

```sh
go mod download
```

## Usage

### Running Locally

1. Set up your environment variables:

```sh
export AWS_DEFAULT_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your-access-key-id
export AWS_SECRET_ACCESS_KEY=your-secret-access-key
export JWT_SECRET=your-jwt-secret
export BIND_ADDRESS=:8080
export ENVIRONMENT=local
export DYNAMODB_ENDPOINT=http://localhost:4566
export S3_ENDPOINT=http://localhost:4566
```

2. Run the application:

```sh
go run cmd/platoapi/main.go
```

### Running with Docker Compose

1. Ensure you have Docker and Docker Compose installed.

2. Build and start the services using Docker Compose:

```sh
docker-compose up --build
```

This command will build the Docker image and start both the Plato API and LocalStack services.

## Configuration

Configuration is managed through environment variables. The following environment variables are available:

- `AWS_DEFAULT_REGION`: The AWS region.
- `AWS_ACCESS_KEY_ID`: The AWS access key ID.
- `AWS_SECRET_ACCESS_KEY`: The AWS secret access key.
- `JWT_SECRET`: The secret key used for JWT authentication.
- `BIND_ADDRESS`: The address the server will bind to (default is `:8080`).
- `ENVIRONMENT`: The environment in which the application is running (`local`, `development`, `production`, `test`).
- `DYNAMODB_ENDPOINT`: The endpoint for DynamoDB (used for local development with LocalStack).
- `S3_ENDPOINT`: The endpoint for S3 (used for local development with LocalStack).
- `CLOUDWATCH_ENDPOINT`: The endpoint for CloudWatch (used for local development with LocalStack).

## API Documentation

The API documentation is generated using OpenAPI and can be accessed at `http://localhost:8080/swagger/index.html` when the server is running.

## Development

To contribute to this project, follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-branch-name`.
3. Make your changes and commit them: `git commit -m 'Add some feature'`.
4. Push to the branch: `git push origin feature-branch-name`.
5. Submit a pull request.

## Running Tests

### Running Unit Tests

To run the unit tests, use the following command:

```sh
go test ./...
```

### Running e2e Tests

To run the e2e tests, please make sure that the service is running in locally or in docker, and use the following command:

```sh
go test ./... --tags=e2e
```

The following environment variables are available:
- `BASE_URL`: The base URL to run tests against (default is `http://localhost:8080`).