package health

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
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

func (hs *HealthService) GetState(ctx context.Context, state mine_state.State) api.State {
	if state == mine_state.Alive {
		return api.State_Alive
	} else if state == mine_state.Stopped {
		return api.State_Stopped
	} else {
		return api.State_Dead
	}
}
