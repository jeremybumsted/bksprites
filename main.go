package main

import (
	"github.com/alecthomas/kong"

	"github.com/jeremybumsted/bksprites/cmd/controller"
	"github.com/jeremybumsted/bksprites/cmd/spriteman"
)

var cli struct {
	Controller controller.ControllerCmd `cmd:"" help:"start an instance of the sprite stack controller"`
	Manager    spriteman.SpriteManCmd   `cmd:"" help:"start an instance of the sprite manager client"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("bksprites"),
		kong.Description("Run Buildkite agents as Fly.io Sprites"),
		kong.UsageOnError(),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
