#!/bin/bash

# Deprecated, use go run cmd/codegen/main.go
 oapi-codegen -package apispec ./docs/swagger.yaml > pkg/api/apispec/generated.go