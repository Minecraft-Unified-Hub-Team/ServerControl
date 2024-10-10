package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/sirupsen/logrus"
)

func (sch *ServerControlHandler) Start(ctx context.Context, req *api.StartRequest) (*api.StartResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.actionService.Start(ctx)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StartResponse{}, err
}

func (sch *ServerControlHandler) Stop(ctx context.Context, req *api.StopRequest) (*api.StopResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.actionService.Stop(ctx)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.StopResponse{}, err
}

func (sch *ServerControlHandler) Install(ctx context.Context, req *api.InstallRequest) (*api.InstallResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.actionService.Install(ctx, req.Version)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	err = sch.configService.WriteEula(ctx)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	err = sch.configService.WriteJVM(ctx)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.InstallResponse{}, err
}

func (sch *ServerControlHandler) Uninstall(ctx context.Context, req *api.UninstallRequest) (*api.UninstallResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.actionService.Stop(ctx)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	err = sch.actionService.Uninstall(ctx)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.UninstallResponse{}, err
}
