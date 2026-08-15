#!/usr/bin/env bash
#MISE description="Check Go formatting"
#MISE dir="{{config_root}}"
#USAGE flag "--write" help="Rewrite files instead of listing them"
set -euo pipefail

if [ "${usage_write:-false}" = "true" ]; then
    exec gofmt -w .
fi

unformatted="$(gofmt -l .)"
[ -z "${unformatted}" ] || { echo "gofmt needed:"; echo "${unformatted}"; exit 1; }
