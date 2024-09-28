package server_control

import (
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/action"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/health"
)

func NewServerControlHandler(
	actionService action.ActionInterface,
	healthService health.HealthInterface,
) (*ServerControlHandler, error) {
	return &ServerControlHandler{
		actionService: actionService,
		healthService: healthService,
	}, nil
}

type ServerControlHandler struct {
	api.UnimplementedActionServer
	api.UnimplementedHealthServer

	actionService action.ActionInterface
	healthService health.HealthInterface
}
