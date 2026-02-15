package sprites

import (
	"fmt"
	"os"

	sprites "github.com/superfly/sprites-go"
)

func RunJob(jobUUID string) error {
	spriteAuthToken := os.Getenv("SPRITE_API_TOKEN")
	client := sprites.New(spriteAuthToken)

	sprite := client.Sprite("bk-test-1")

	cmd := sprite.Command("buildkite-agent", "start", "--acquire-job", jobUUID, "--skip-checkout")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("Output: %s", output)
	return nil
}
