package version

import "fmt"

var (
	Version   string
	CommitSHA string
	BuildTime string
)

type VersionCmd struct{}

func (v *VersionCmd) Run() error {
	fmt.Printf("bksprites %s\n", Version)
	fmt.Printf("  commit: %s\n", CommitSHA)
	fmt.Printf("  built:  %s\n", BuildTime)
	return nil
}
