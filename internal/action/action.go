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

	State *mine_state.State // channel that stores state of server
}

func NewActionService() (*ActionService, error) {
	currentState, _ := mine_state.NewState(mine_state.Stopped) // TODO set mine_state.Stopped here when install will be completed
	return &ActionService{
		State: currentState,
	}, nil
}

func (as *ActionService) downloadJar(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.downloadJar(ctx, %s)", version) + ": %w"

	/* configure url address for downloading correct forge jar */
	url := fmt.Sprintf(baseURL+installerName, version, version)

	/* prepare command and arguments */
	command := "wget"
	args := append(
		make([]string, 0),
		"-P",
		serverPath,
		url,
	)

	/* dowload jar file */
	err = mine_os.ExecCtx(ctx, command, args)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) installJar(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.installJar(ctx, %s)", version) + ": %w"

	/* prepare command and arguments */
	command := "java"
	args := append(
		make([]string, 0),
		"-jar",
		fmt.Sprintf(serverPath+installerName, version),
		"--installServer",
		serverPath,
	)

	/* configure command for installing */
	err = mine_os.ExecCtx(ctx, command, args)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) removeJar(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.removeJar(ctx, %s)", version) + ": %w"

	/* prepare command and arguments */
	command := "rm"
	args := append(
		make([]string, 0),
		fmt.Sprintf(serverPath+installerName, version),
	)

	/* remove used jar file */
	err = mine_os.ExecCtx(ctx, command, args)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) Install(ctx context.Context, version string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("ActionService.Install(ctx, %s)", version) + ": %w"

	/* download installer for minecraft server with setted version */
	err = as.downloadJar(ctx, version)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	/* install server file via java installer for jar files */
	err = as.installJar(ctx, version)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	/* remove used jar file after installation */
	err = as.removeJar(ctx, version)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (as *ActionService) Start(ctx context.Context) error {
	var err error = nil

	/* check that server has not been already started */
	if as.State.IsAlive() {
		// return err
		return fmt.Errorf("server has been already started") // TODO verify that we use fmt.Errorf for creating errors
	}

	/* create aliveness context for server run */
	as.aliveCtx, as.stopCtx = context.WithCancel(context.Background())

	/* prepare command and arguments */
	command := "/bin/bash"
	args := append(
		make([]string, 0),
		"-c",
		fmt.Sprintf("%s %s && ./%s", cd, serverPath, run),
	)
	logrus.Debugln(command, args)

	/* start server in goroutine */
	go func() {
		as.State.Set(mine_state.Alive)
		status, err := mine_os.ManagedExecCtx(as.aliveCtx, command, args)
		logrus.Debugln("get error in start:", err)
		if status == 1 {
			as.State.Set(mine_state.Stopped)
		} else {
			as.State.Set(mine_state.Dead)
		}
	}()

	/* always okay */
	return err
}

func (as *ActionService) Stop(ctx context.Context) error {
	var err error = nil

	/* call cancel context function */
	as.stopCtx()

	/* always okay */
	return err
}
