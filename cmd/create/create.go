package create

import "github.com/jeremybumsted/bksprites/internal/models"

type CreateCmd struct {
	sprite models.Sprite
}

func (c *CreateCmd) Run() {
}
