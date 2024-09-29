package server_control

import (
	"context"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/sirupsen/logrus"
)

func (sch *ServerControlHandler) WriteSettings(ctx context.Context, req *api.WriteSettingsRequest) (*api.WriteSettingsResponse, error) {
	var err error = nil

	logrus.Debug(req)

	err = sch.configService.WriteSettings(context.Background())
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.WriteSettingsResponse{}, err
}
