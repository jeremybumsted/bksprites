#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y

mise install

echo "--- Running :go: Tests"

cd app &&
  mise x gotestsum@latest -- gotestsum --format testname --junitfile unit-tests.xml --junitfile-testcase-classname relative -- -coverprofile=cover.out ./...

echo "--- Done"
exit 0
