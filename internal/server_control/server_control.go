package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_state"
)

func NewServerControlHandler(
	actionService ActionInterface,
	healthService HealthInterface,
) (*ServerControlHandler, error) {
	return &ServerControlHandler{
		actionService: actionService,
		healthService: healthService,
	}, nil
}

type (
	ServerControlHandler struct {
		api.UnimplementedActionServer
		api.UnimplementedHealthServer

		actionService ActionInterface
		healthService HealthInterface
	}

	ActionInterface interface {
		Start(context.Context) error
		Stop(context.Context) error
		Install(context.Context, string) error
	}

	HealthInterface interface {
		GetState(context.Context) (mine_state.State, error)
		Ping(context.Context) error
	}
)
