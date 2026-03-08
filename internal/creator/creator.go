package creator

import (
	"github.com/jeremybumsted/bksprites/internal/models"
	"github.com/superfly/sprites-go"
)

type Creator struct {
	Client sprites.Client
}

func CreateSprite() (models.Sprite, error) {
	var sprite models.Sprite
	sprite.Name = "foo"

	return sprite, nil
}
