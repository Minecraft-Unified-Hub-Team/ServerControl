package impl

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/sirupsen/logrus"
)

func (sch *ServiceControlHandler) Start(ctx context.Context, req *api.StartRequest) (*api.StartResponse, error) {
	logrus.Debug(req)

	err := sch.actionService.Start(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StartResponse{}, nil
}
