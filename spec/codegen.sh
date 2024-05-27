#!/bin/bash

# Generate the server side stubs for the openapi spec, excluding the service modules (very important)
openapi-generator generate -g go-server -i ./spec/spec.v1.yaml -o ../app --additional-properties=router=chi,packageName=openapi,outputAsLibrary=true,sourceFolder=openapi,onlyInterfaces=true