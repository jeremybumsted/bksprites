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

## Quickstart

Install from the latest release and run the controller:

```bash
bksprites controller --agent-token=$BUILDKITE_SPRITE_AGENT_TOKEN --queue="sprites"
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

## Contributing

PRs welcome 🙂 this is just an experimental repo to build my
understanding of Sprites and the Stacks API

## License

See [LICENSE](/LICENSE)
