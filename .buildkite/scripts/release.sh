#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y
mise install

echo "--- :package: Running GoReleaser"

# Check if this is a tag build or a snapshot
if [ -n "${BUILDKITE_TAG:-}" ]; then
  echo "Creating release for tag: $BUILDKITE_TAG"
  RELEASE_MODE="release"
else
  echo "Creating snapshot build (no tag detected)"
  RELEASE_MODE="snapshot"
fi

# Ensure GITHUB_TOKEN is set for releases
if [ "$RELEASE_MODE" = "release" ] && [ -z "${GITHUB_TOKEN:-}" ]; then
  echo "ERROR: GITHUB_TOKEN is required for releases"
  echo "Please set GITHUB_TOKEN in Buildkite secrets"
  exit 1
fi

# Run goreleaser
if [ "$RELEASE_MODE" = "release" ]; then
  echo "Building and releasing binaries..."
  mise x goreleaser@latest -- goreleaser release --clean
else
  echo "Building snapshot (not publishing)..."
  mise x goreleaser@latest -- goreleaser release --snapshot --clean --skip=publish
fi

echo "--- :white_check_mark: GoReleaser completed successfully"

# Show artifacts
if [ -d "dist" ]; then
  echo "--- :package: Built artifacts:"
  ls -lh dist/ | grep -E '\.(tar\.gz|zip|txt)$' || true
fi

# Upload artifacts to Buildkite
if [ "$RELEASE_MODE" = "snapshot" ] && [ -d "dist" ]; then
  echo "--- Uploading snapshot artifacts"
  buildkite-agent artifact upload "dist/*.tar.gz"
  buildkite-agent artifact upload "dist/*.zip"
  buildkite-agent artifact upload "dist/checksums.txt"
fi
