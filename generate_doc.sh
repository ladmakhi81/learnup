#!/usr/bin/env bash

set -eu
set -o pipefail
directories="$(find {internals,pkg} -type d | paste -d',' -s -)"
swag fmt -g ./cmd/api/main.go --dir "$directories"
swag init --parseDependency -g ../cmd/api/main.go --dir "$directories" --output ./docs

