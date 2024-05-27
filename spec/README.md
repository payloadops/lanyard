# Payload's OpenAPI Specification

The OpenAPI v3 specification for Payload's API.

## What is OpenAPI?

From the [OpenAPI Specification](https://swagger.io/specification/):

> The OpenAPI Specification (OAS) defines a standard, language-agnostic interface to RESTful APIs which allows both humans and computers to discover and understand the capabilities of the service without access to source code, documentation, or through network traffic inspection. When properly defined, a consumer can understand and interact with the remote service with a minimal amount of implementation logic.

> An OpenAPI definition can then be used by documentation generation tools to display the API, code generation tools to generate servers and clients in various programming languages, testing tools, and many other use cases.

## Requirements

The following dependencies are necessary to run the Swagger UI and mock server.

* [Docker Compose >= 2.26.1](https://docs.docker.com/compose/install/)
* [Docker >= 26.0.0](https://docs.docker.com/get-docker/)
* [OpenAPI Generator >= 5.1.1](https://openapi-generator.tech/docs/installation)

## Development

Run docker compose in the root directory of this project, which will expose Swagger UI,
and the mock server.

```/bin/bash
docker-compose up
```

## Code Generation
To generate server-side stubs for the OpenAPI spec, excluding the service modules, use the OpenAPI Generator with the following command:

```/bin/bash
./codegen.sh
```