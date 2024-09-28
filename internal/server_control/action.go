package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/sirupsen/logrus"
)

func (sch *ServerControlHandler) Start(ctx context.Context, req *api.StartRequest) (*api.StartResponse, error) {
	logrus.Debug(req)

	err := sch.actionService.Start(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StartResponse{}, nil
}

func (sch *ServerControlHandler) Stop(ctx context.Context, req *api.StopRequest) (*api.StopResponse, error) {
	logrus.Debug(req)

	err := sch.actionService.Stop(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StopResponse{}, nil
}
