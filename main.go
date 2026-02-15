package main

import (
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"

	"github.com/jeremybumsted/bksprites/cmd/controller"
	"github.com/jeremybumsted/bksprites/cmd/runner"
)

var cli struct {
	Controller controller.ControllerCmd `cmd:"" help:"start an instance of the sprite stack controller"`
	Runner     runner.RunnerCmd         `cmd:"" help:"start an instance of the sprite stack runner"`
}

func main() {
	log.Info("Welcome to the BK Sprite Stack")

	ctx := kong.Parse(&cli,
		kong.Name("bksprites"),
		kong.Description("Run Buildkite agents as Fly.io Sprites"),
		kong.UsageOnError(),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
