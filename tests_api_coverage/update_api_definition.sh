#!/bin/bash

curl https://raw.githubusercontent.com/anyproto/anytype-heart/refs/heads/main/core/api/docs/openapi.json -o api_definition.json
if [ $? -ne 0 ]; then
    echo "Failed to fetch the API definition file."
    exit 1
fi