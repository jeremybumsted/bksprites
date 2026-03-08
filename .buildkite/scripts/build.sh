#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y
mise install

echo "--- :docker: Building container image with ko"

# Default to GitHub Container Registry if not specified
: "${KO_DOCKER_REPO:=ghcr.io/jeremybumsted/bksprites}"
export KO_DOCKER_REPO

# Determine build mode based on tag presence
if [ -n "${BUILDKITE_TAG:-}" ]; then
  echo "Building for tag release: $BUILDKITE_TAG"
  BUILD_MODE="release"
  LOCAL_BUILD=""
else
  echo "Building locally (no tag detected - will not push to registry)"
  BUILD_MODE="local"
  LOCAL_BUILD="--local"
fi

# Determine image tags
# Priority: explicit TAG env var, git tag, git SHA
if [ -n "${TAG:-}" ]; then
  IMAGE_TAG="$TAG"
elif [ -n "${BUILDKITE_TAG:-}" ]; then
  IMAGE_TAG="$BUILDKITE_TAG"
else
  # Use short git SHA for local builds
  IMAGE_TAG="${BUILDKITE_COMMIT:0:7}"
fi

# Build configuration based on GOARCH
# If GOARCH is set, build single-arch using env vars (ko requires both GOOS and GOARCH)
# Otherwise, build multi-arch using --platform flag
if [ -n "${GOARCH:-}" ]; then
  # Single-arch build - ko uses GOOS/GOARCH env vars, don't pass --platform
  export GOOS=linux
  export GOARCH
  IMAGE_TAG="${IMAGE_TAG}-${GOARCH}"
  echo "Building for single architecture: linux/${GOARCH}"

  if [ "$BUILD_MODE" = "local" ]; then
    echo "Building image locally: ${IMAGE_TAG} (single-arch)"
  else
    echo "Building image: ${KO_DOCKER_REPO}:${IMAGE_TAG} (single-arch)"
  fi

  # Build image (with or without push)
  mise x -- ko build \
    --tags="${IMAGE_TAG}" \
    --bare \
    ${LOCAL_BUILD} \
    .
else
  # Multi-arch build - use --platform flag
  PLATFORMS="linux/amd64,linux/arm64"
  echo "Building multi-arch image: $PLATFORMS"

  if [ "$BUILD_MODE" = "local" ]; then
    echo "Building image locally: ${IMAGE_TAG} (multi-arch)"
  else
    echo "Building image: ${KO_DOCKER_REPO}:${IMAGE_TAG} (multi-arch)"
  fi

  # Build image (with or without push)
  mise x -- ko build \
    --platform="${PLATFORMS}" \
    --tags="${IMAGE_TAG}" \
    --bare \
    ${LOCAL_BUILD} \
    .
fi

if [ "$BUILD_MODE" = "local" ]; then
  echo "--- :white_check_mark: Image built locally (not pushed)"
  echo "Local image tag: ${IMAGE_TAG}"
else
  echo "--- :white_check_mark: Image built and pushed successfully"
  echo "Image: ${KO_DOCKER_REPO}:${IMAGE_TAG}"
fi

# Also tag as 'latest' if this is a git tag build (only for multi-arch)
if [ -n "${BUILDKITE_TAG:-}" ] && [ -z "${GOARCH:-}" ]; then
  echo "--- :label: Tagging as latest"
  mise x -- ko build \
    --platform="${PLATFORMS}" \
    --tags="latest" \
    --bare \
    .
  echo "Image: ${KO_DOCKER_REPO}:latest"
fi
