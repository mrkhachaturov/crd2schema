#!/usr/bin/env bash
#MISE description="Lint Go, shell, yaml, md and toml via flint"
#MISE dir="{{config_root}}"
#USAGE flag "--fix" help="Apply what flint can fix, report the rest"
set -euo pipefail

if [ "${usage_fix:-false}" = "true" ]; then
    exec flint run --full --fix --allow-fixed
fi

exec flint run --full
