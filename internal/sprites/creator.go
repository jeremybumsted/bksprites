package sprites

import (
	"context"
)

func (s *SpriteHandler) CreateAgentSprite(name string) (*AgentSprite, error) {
	ctx := context.Background()

	r, err := s.Client.CreateSprite(ctx, name, nil)
	if err != nil {
		return nil, err
	}

	as := AgentSprite{
		Name:    r.Name(),
		Address: r.URL,
	}

	return &as, nil
}
