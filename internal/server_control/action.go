package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/sirupsen/logrus"
)

func (sch *ServerControlHandler) Start(ctx context.Context, req *api.StartRequest) (*api.StartResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.actionService.Start(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StartResponse{}, err
}

func (sch *ServerControlHandler) Stop(ctx context.Context, req *api.StopRequest) (*api.StopResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.actionService.Stop(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StopResponse{}, err
}

func (sch *ServerControlHandler) Install(ctx context.Context, req *api.InstallRequest) (*api.InstallResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.actionService.Install(context.Background(), req.Version)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	err = sch.configService.WriteEula(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	err = sch.configService.WriteJVM(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.InstallResponse{}, err
}
