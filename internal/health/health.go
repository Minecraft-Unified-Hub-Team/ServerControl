package health

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_state"
)

func NewHealthService() (*HealthService, error) {
	return &HealthService{}, nil
}

type HealthService struct{}

func (hs *HealthService) Ping(ctx context.Context) error {
	var err error = nil

	return err
}

func (hs *HealthService) GetState(ctx context.Context, state *mine_state.State) (string, error) {
	var err error = nil

	return state.String(), err
}
