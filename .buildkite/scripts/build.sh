#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y
mise install

echo "--- :docker: Building and pushing container image with ko"

# Default to GitHub Container Registry if not specified
: "${KO_DOCKER_REPO:=ghcr.io/jeremybumsted/bksprites}"
export KO_DOCKER_REPO

# Determine image tags
# Priority: explicit TAG env var, git tag, git SHA
if [ -n "${TAG:-}" ]; then
  IMAGE_TAG="$TAG"
elif [ -n "${BUILDKITE_TAG:-}" ]; then
  IMAGE_TAG="$BUILDKITE_TAG"
else
  # Use short git SHA
  IMAGE_TAG="${BUILDKITE_COMMIT:0:7}"
fi

echo "Building image: ${KO_DOCKER_REPO}:${IMAGE_TAG}"
echo "Platforms: linux/amd64, linux/arm64"

# Build and push multi-arch image
# --bare: only push the image, no additional tags
# --platform: build for multiple architectures
# --tags: explicit tag to use
mise x ko@latest -- ko build \
  --platform=linux/amd64,linux/arm64 \
  --tags="${IMAGE_TAG}" \
  --bare \
  .

echo "--- :white_check_mark: Image built and pushed successfully"
echo "Image: ${KO_DOCKER_REPO}:${IMAGE_TAG}"

# Also tag as 'latest' if this is a git tag build
if [ -n "${BUILDKITE_TAG:-}" ]; then
  echo "--- :label: Tagging as latest"
  mise x ko@latest -- ko build \
    --platform=linux/amd64,linux/arm64 \
    --tags="latest" \
    --bare \
    .
  echo "Image: ${KO_DOCKER_REPO}:latest"
fi
