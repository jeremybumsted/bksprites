# bksprites - Run Buildkite agents on Fly.io sprites

An experimental implementation of the [Buildkite Stacks API](https://buildkite.com/docs/apis/agent-api/stacks)
to run agents using [Sprites](https://docs.sprites.dev/)

## Why?

Buildkite allows pretty much unlimited freedom in terms of runner architecture
so long as the buildkite-agent can run on it. Sprites are an interesting stateful
sandbox offered by the folks over at [fly.io](fly.io). Since you can create (basically)
unlimited sprites, and they can maintain state, using them for autoscaling Buildkite
agent workloads makes them a compelling offer for folks who want fast CI at
small-medium scale (and maybe even larger)

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/jeremybumsted/bksprites/releases).

Available for:
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

Example for Linux:
```bash
# Download and extract (replace VERSION with the actual version)
curl -L https://github.com/jeremybumsted/bksprites/releases/download/VERSION/bksprites_VERSION_linux_amd64.tar.gz | tar xz

# Move to your PATH
sudo mv bksprites /usr/local/bin/

# Verify installation
bksprites version
```

### Container Image

Container images are available from GitHub Container Registry for both `linux/amd64` and `linux/arm64`:

```bash
docker pull ghcr.io/jeremybumsted/bksprites:latest
```

## Quickstart

Run the controller:

```bash
bksprites controller --agent-token=$BUILDKITE_SPRITE_AGENT_TOKEN --queue="sprites"
```

Or with Docker:

```bash
docker run -e BUILDKITE_SPRITE_AGENT_TOKEN ghcr.io/jeremybumsted/bksprites:latest \
  controller --agent-token=$BUILDKITE_SPRITE_AGENT_TOKEN --queue="sprites"
```

## Controller

The controller connects to the stacks API and polls for scheduled jobs on the configured
queue; assigning them to a sprite that has been configured with the buildkite-agent

To keep things simple, and self-contained, the controller

### Options

Add the options

## Development

This project uses `mise-en-place` to manage dependencies. Run `mise install`
after cloning the repo to ensure your dependencies are up to date.

### Building Binaries Locally

This project uses [GoReleaser](https://goreleaser.com/) for building release binaries:

```bash
# Install dependencies (includes goreleaser)
mise install

# Build snapshot (local testing, no publishing)
mise x goreleaser -- goreleaser release --snapshot --clean

# Built binaries will be in dist/ directory
./dist/bksprites_darwin_arm64/bksprites version
```

You can also use standard Go tooling:
```bash
go build -o bksprites .
./bksprites version
```

Configuration is in `.goreleaser.yaml` which specifies build targets, ldflags, and release options.

### Building Container Images Locally

This project uses [ko](https://ko.build/) for building container images:

```bash
# Build and push to a registry (defaults to ghcr.io/jeremybumsted/bksprites)
mise x ko -- ko build .

# Build for local testing without pushing
mise x ko -- ko build --local .

# Build with custom tag
mise x ko -- ko build --tags=v1.0.0 .
```

Configuration is in `.ko.yaml` which specifies the base image (Chainguard static), platforms (amd64/arm64), and build flags.

## Contributing

PRs welcome 🙂 this is just an experimental repo to build my
understanding of Sprites and the Stacks API

## License

See [LICENSE](/LICENSE)
