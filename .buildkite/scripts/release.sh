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

# Ensure GITHUB_API_TOKEN is set for releases
# GoReleaser expects GITHUB_TOKEN, so we'll export it from GITHUB_API_TOKEN
if [ "$RELEASE_MODE" = "release" ]; then
  if [ -z "${GITHUB_API_TOKEN:-}" ]; then
    echo "ERROR: GITHUB_API_TOKEN is required for releases"
    echo "Please set GITHUB_API_TOKEN in Buildkite secrets"
    exit 1
  fi
  # Export as GITHUB_TOKEN for goreleaser
  export GITHUB_TOKEN="${GITHUB_API_TOKEN}"
fi

# Check for platform filtering via environment variables
CONFIG_FILE=".goreleaser.yaml"
CLEANUP_CONFIG=false

if [ -n "${GOOS:-}" ] && [ -n "${GOARCH:-}" ]; then
  echo "Building for platform: $GOOS/$GOARCH"

  # Create a temporary filtered config
  CONFIG_FILE=".goreleaser.tmp.yaml"
  CLEANUP_CONFIG=true

  # Generate filtered config with only the specified platform
  cat > "$CONFIG_FILE" << EOF
# Temporary goreleaser config for $GOOS/$GOARCH
version: 2

before:
  hooks:
    - go mod tidy

builds:
  - id: bksprites
    main: .
    binary: bksprites
    goos:
      - $GOOS
    goarch:
      - $GOARCH
    env:
      - CGO_ENABLED=0
    flags:
      - -trimpath
    ldflags:
      - -s
      - -w
      - -X main.Version={{.Version}}
      - -X main.CommitSHA={{.ShortCommit}}
      - -X main.BuildTime={{.Date}}

archives:
  - id: default
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- .Os }}_
      {{- .Arch }}
    formats:
      - tar.gz
    format_overrides:
      - goos: windows
        formats: [zip]
    files:
      - README.md
      - LICENSE

checksum:
  name_template: 'checksums.txt'
  algorithm: sha256

snapshot:
  version_template: "{{ incpatch .Version }}-next"

metadata:
  mod_timestamp: '{{ .CommitTimestamp }}'
EOF
elif [ -n "${GOOS:-}" ] || [ -n "${GOARCH:-}" ]; then
  echo "ERROR: Both GOOS and GOARCH must be set together for platform filtering"
  exit 1
else
  echo "Building for all platforms"
fi

# Run goreleaser
if [ "$RELEASE_MODE" = "release" ]; then
  echo "Building and releasing binaries..."
  mise x goreleaser@latest -- goreleaser release --clean --config="$CONFIG_FILE"
else
  echo "Building snapshot (not publishing)..."
  mise x goreleaser@latest -- goreleaser release --snapshot --clean --skip=publish --config="$CONFIG_FILE"
fi

# Cleanup temporary config
if [ "$CLEANUP_CONFIG" = true ]; then
  rm -f "$CONFIG_FILE"
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
