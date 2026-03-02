#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y

mise install

echo "--- Running :go: Tests"

mise x gotestsum@latest -- gotestsum --format testname --junitfile unit-tests.xml --junitfile-testcase-classname relative -- -coverprofile=cover.out ./...

echo "--- Uploading artifacts"
buildkite-agent artifact upload "cover.out"

cat cover.out | buildkite-agent annotate --style "info" --context "gotestsum"

echo "--- Annotating build"
exit 0
