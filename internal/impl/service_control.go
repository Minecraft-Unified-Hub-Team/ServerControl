package impl

import (
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/app"
)

func NewServiceControlHandler(actionService app.ActionService) (*ServiceControlHandler, error) {
	return &ServiceControlHandler{
		actionService: actionService,
	}, nil
}

type ServiceControlHandler struct {
	api.UnimplementedActionServiceServer
	actionService app.ActionService
}
