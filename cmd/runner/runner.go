package runner

import (
	"fmt"
)

type RunnerCmd struct{}

func (r *RunnerCmd) Run() error {
	fmt.Printf("TODO: Build the runner :)")
	return nil
}
