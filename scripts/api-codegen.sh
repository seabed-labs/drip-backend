#!/bin/bash

 oapi-codegen -package apispec ./docs/swagger.yaml > pkg/api/apispec/generated.go