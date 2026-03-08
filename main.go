package main

import (
	"github.com/alecthomas/kong"

	"github.com/jeremybumsted/bksprites/cmd/controller"
	"github.com/jeremybumsted/bksprites/cmd/create"
	"github.com/jeremybumsted/bksprites/cmd/version"
)

// Version information - injected via ldflags at build time
var (
	Version   = "dev"
	CommitSHA = "none"
	BuildTime = "unknown"
)

var cli struct {
	Controller controller.ControllerCmd `cmd:"" help:"start an instance of the sprite stack controller"`
	Create     create.CreateCmd         `cmd:"" help:"create a new pre-configured sprite"`
	Version    version.VersionCmd       `cmd:"" help:"show version information"`
}

func main() {
	// Set version information for the version command
	version.Version = Version
	version.CommitSHA = CommitSHA
	version.BuildTime = BuildTime

	ctx := kong.Parse(&cli,
		kong.Name("bksprites"),
		kong.Description("Run Buildkite agents as Fly.io Sprites"),
		kong.UsageOnError(),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
