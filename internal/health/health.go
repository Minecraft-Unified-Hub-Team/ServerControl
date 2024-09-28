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
	return nil
}

func (hs *HealthService) GetState(ctx context.Context) (mine_state.State, error) {
	return mine_state.State{}, nil
}
