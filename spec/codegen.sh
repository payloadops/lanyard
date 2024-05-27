#!/bin/bash

openapi-generator generate -g go-server -i ./spec/spec.v1.yaml -o ../api --additional-properties=router=chi,packageName=openapi,outputAsLibrary=true,sourceFolder=openapi,onlyInterfaces=true