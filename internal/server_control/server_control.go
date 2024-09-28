package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
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
	}

	HealthInterface interface {
		GetState(context.Context) error
		Ping(context.Context) error
	}
)
