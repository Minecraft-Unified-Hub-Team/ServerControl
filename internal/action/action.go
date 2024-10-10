package action

import (
	"context"
	"fmt"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_os"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_state"
	"github.com/sirupsen/logrus"
)

const (
	cd            = "cd"
	run           = "run.sh"
	serverPath    = "/server"
	baseURL       = "https://maven.minecraftforge.net/net/minecraftforge/forge/%s"
	installerName = "/forge-%s-installer.jar"
)

type ActionService struct {
	aliveCtx context.Context    // context that continues until server is stopped or dead
	stopCtx  context.CancelFunc // function that cancels server binary execution

	syncedState *mine_state.SyncedState // channel that stores state of server
}

func NewActionService() (*ActionService, error) {
	currentState, _ := mine_state.NewSyncedState(mine_state.Stopped) 
	return &ActionService{
		syncedState: currentState,
	}, nil
}

func (as *ActionService) downloadJar(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.downloadJar(ctx, %s)", version) + ": %w"

	url := fmt.Sprintf(baseURL+installerName, version, version)

	command := "wget"
	args := append(
		make([]string, 0),
		"-P",
		serverPath,
		url,
	)

	err = mine_os.ExecCtx(ctx, command, args)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) installJar(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.installJar(ctx, %s)", version) + ": %w"

	command := "java"
	args := append(
		make([]string, 0),
		"-jar",
		fmt.Sprintf(serverPath+installerName, version),
		"--installServer",
		serverPath,
	)

	err = mine_os.ExecCtx(ctx, command, args)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) removeJar(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.removeJar(ctx, %s)", version) + ": %w"

	command := "rm"
	args := append(
		make([]string, 0),
		fmt.Sprintf(serverPath+installerName, version),
	)

	err = mine_os.ExecCtx(ctx, command, args)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) Install(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.Install(ctx, %s)", version) + ": %w"

	err = as.downloadJar(ctx, version)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	err = as.installJar(ctx, version)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	err = as.removeJar(ctx, version)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) Uninstall(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "ActionService.Uninstall(ctx): %w"

	command := "rm"
	args := append(
		make([]string, 0),
		"-rf",
		serverPath+"/*",
	)

	err = mine_os.ExecCtx(ctx, command, args)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) Start(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "ActionService.Start(ctx): %w"

	if as.syncedState.IsAlive() {
		return fmt.Errorf(errorFormat, "server has been already started") // TODO verify that we use fmt.Errorf for creating errors
	}

	as.aliveCtx, as.stopCtx = context.WithCancel(context.Background())

	command := "/bin/bash"
	args := append(
		make([]string, 0),
		"-c",
		fmt.Sprintf("%s %s && ./%s", cd, serverPath, run),
	)
	logrus.Debugln(command, args)

	go func() {
		as.syncedState.Set(mine_state.Alive)
		status, err := mine_os.ManagedExecCtx(as.aliveCtx, command, args)
		if err != nil {
			logrus.Debugln("get error in managed start:", err)
		}
		if status == mine_os.NO_ERROR {
			as.syncedState.Set(mine_state.Stopped)
		} else {
			as.syncedState.Set(mine_state.Dead)
		}
	}()

	return err
}

func (as *ActionService) Stop(ctx context.Context) error {
	as.stopCtx()
	return nil
}

func (as *ActionService) GetState(ctx context.Context) mine_state.State {
	return as.syncedState.State()
}
