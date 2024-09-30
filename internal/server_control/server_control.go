package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_state"
)

func NewServerControlHandler(
	actionService ActionInterface,
	healthService HealthInterface,
	configService ConfigInterface,
) (*ServerControlHandler, error) {
	return &ServerControlHandler{
		actionService: actionService,
		healthService: healthService,
		configService: configService,
	}, nil
}

type (
	ServerControlHandler struct {
		api.UnimplementedActionServer
		api.UnimplementedHealthServer
		api.UnimplementedConfigServer

		actionService ActionInterface
		healthService HealthInterface
		configService ConfigInterface
	}

	ActionInterface interface {
		Start(context.Context) error
		Stop(context.Context) error
		Install(context.Context, string) error

		State(context.Context) *mine_state.State
	}

	HealthInterface interface {
		GetState(context.Context, *mine_state.State) (string, error)
		Ping(context.Context) error
	}

	ConfigInterface interface {
		WriteEula(context.Context) error
		WriteJVM(context.Context) error
		WriteSettings(context.Context) error
	}
)
