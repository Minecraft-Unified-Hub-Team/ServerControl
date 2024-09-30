package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/sirupsen/logrus"
)

func (sch *ServerControlHandler) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.healthService.Ping(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.PingResponse{}, err
}

func (sch *ServerControlHandler) GetState(ctx context.Context, req *api.StateRequest) (*api.StateResponse, error) {
	var err error = nil

	stateString, err := sch.healthService.GetState(context.Background(), sch.actionService.State(ctx))
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StateResponse{State: stateString}, err
}
