#!/bin/bash

# This script will execute when scaling out your sprites
# The controller will use this script as the provision script
# by default, but you're free to provide your own provision
# script to best suit your needs

# This script is used to configure a sprite to be ready to
# Run our jobs to build bksprites :)
echo "=== Installing Dependencies ===\n"

# The Buildkite agent install uses the documented installation steps from
# https://buildkite.com/docs/agent/self-hosted/install/linux
echo "Configuring the buildkite Agent\n"
TOKEN="${BUILDKITE_SPRITE_AGENT_TOKEN}" bash -c "$(curl -sL https://raw.githubusercontent.com/buildkite/agent/main/install.sh)"

echo "Configure mise-en-place\n"
curl https://mise.run | sh
mise --version
