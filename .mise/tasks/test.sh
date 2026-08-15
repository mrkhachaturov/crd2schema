#!/usr/bin/env bash
#MISE description="Run the unit tests"
#MISE dir="{{config_root}}"
#USAGE flag "--junit <path>" help="Also write a JUnit report"
set -euo pipefail

if [ -n "${usage_junit:-}" ]; then
    exec gotestsum --junitfile "${usage_junit}" -- ./...
fi

exec gotestsum -- ./...
