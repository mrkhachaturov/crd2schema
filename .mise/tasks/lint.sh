#!/usr/bin/env bash
#MISE description="Run golangci-lint"
#MISE dir="{{config_root}}"
set -euo pipefail
exec golangci-lint run ./...
