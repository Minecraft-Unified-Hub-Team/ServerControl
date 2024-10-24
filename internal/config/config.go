package config

import (
	"context"
	"fmt"
	"os"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_settings"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	serverPath  = "/server"
	eulaConfig  = "/eula.txt"
	eulaContent = "eula=true"
	jvmConfig   = "/user_jvm_args.txt"

	bytesInGB = 1024 * 1024 * 1024
	bytesInMB = 1024 * 1024
)

type ConfigService struct{}

func NewConfigService() (*ConfigService, error) {
	return &ConfigService{}, nil
}

func (cs *ConfigService) WriteEula(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "ConfigService.WriteEula(ctx): %w"

	eulaFile, err := os.OpenFile(serverPath+eulaConfig, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	defer eulaFile.Close()

	_, err = eulaFile.WriteString(eulaContent)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (cs *ConfigService) WriteJVM(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "ConfigService.WriteJVM(ctx): %w"

	jvmFile, err := os.OpenFile(serverPath+jvmConfig, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	defer jvmFile.Close()

	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	maxMemory := int(v.Total / bytesInMB)
	_, err = jvmFile.WriteString(fmt.Sprintf("-Xmx%dM", maxMemory))
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (cs *ConfigService) WriteSettings(ctx context.Context, m map[string]string) error {
	return mine_settings.WriteSettingsConfig(ctx, "/server", m)
}

func (cs *ConfigService) ReadSettings(ctx context.Context) (map[string]string, error) {
	return mine_settings.ReadSettingsConfig(ctx, "/server")
}
