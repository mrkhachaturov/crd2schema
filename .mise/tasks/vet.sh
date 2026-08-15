#!/usr/bin/env bash
#MISE description="Run go vet"
#MISE dir="{{config_root}}"
set -euo pipefail
exec go vet ./...
