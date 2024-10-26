package server_control

import (
	"context"
	"time"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
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
		Uninstall(context.Context) error
	}

	HealthInterface interface {
		Start(ctx context.Context, refreshTime time.Duration) error
		Stop(ctx context.Context) error

		GetState(context.Context) api.State
		Ping(context.Context) error
	}

	ConfigInterface interface {
		WriteEula(context.Context) error
		WriteJVM(context.Context) error
		WriteSettings(context.Context, map[string]string) error
		ReadSettings(context.Context) (map[string]string, error)
	}
)
