#!/usr/bin/env bash
set -euo pipefail

# Write next to this script regardless of the caller's working directory.
cd "$(dirname "$0")"

url="https://raw.githubusercontent.com/anyproto/anytype-heart/refs/heads/main/core/api/docs/openapi.json"

# -f fails the command on any HTTP error instead of saving the error body.
if ! curl -fsSL "$url" -o api_definition.json; then
    echo "Failed to fetch the API definition file." >&2
    exit 1
fi

echo "Updated api_definition.json"
