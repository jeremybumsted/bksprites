# bksprites - Run Buildkite agents on Fly.io sprites

An experimental implementation of the [Buildkite Stacks API](https://buildkite.com/docs/apis/agent-api/stacks)
to run agents using [Sprites](https://docs.sprites.dev/), bksprites is split
into two components, the Controller and the Manager

## Controller

The Controller connects to the stacks API and polls for scheduled jobs
to be assigned to a Sprite configured with the bksprites Manager.

It requires an agent token and queue

### Options

Add the options

## Manager

The Manager runs on a Sprite that has been configured to run the Buildkite agent.
It creates a simple http server that receives job run requests, in addition to
a /health endpoint to quickly start a warm or cold sprite during autoscaling operations

The manager also ensures that the number of running agents is limited by monitoring
against the sprites resources and the max agents value.

# Development

This project uses `mise-en-place` to manage dependencies. Run `mise install` after cloning the repo
to ensure your dependencies are up to date.

# Contributing

PRs welcome, this is just an experimental repo to play with Sprites and the Stacks API

# License

See [LICENSE.md]
