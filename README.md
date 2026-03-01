# bksprites - Run Buildkite agents on Fly.io sprites

An experimental implementation of the [Buildkite Stacks API](https://buildkite.com/docs/apis/agent-api/stacks)
to run agents using [Sprites](https://docs.sprites.dev/)

## Controller

The Controller connects to the stacks API and polls for scheduled jobs
to be assigned to a Sprite configured with the bksprites Manager.

It requires an agent token and queue

### Options

Add the options

## Development

This project uses `mise-en-place` to manage dependencies. Run `mise install`
after cloning the repo to ensure your dependencies are up to date.

## Contributing

PRs welcome 🙂 this is just an experimental repo
to play with Sprites and the Stacks API

## License

See [LICENSE](/LICENSE)
