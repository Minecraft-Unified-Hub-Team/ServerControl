package action

import (
	"context"
	"fmt"
	"os"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_os"
	"github.com/sirupsen/logrus"
)

const (
	cd            = "cd"
	run           = "run.sh"
	exitcodeFile  = "minecraft_exitcode"
	serverPath    = "/server"
	baseURL       = "https://maven.minecraftforge.net/net/minecraftforge/forge/%s"
	installerName = "/forge-%s-installer.jar"
)

type ActionService struct {
	aliveCtx context.Context    // context that continues until server is stopped or dead
	stopCtx  context.CancelFunc // function that cancels server binary execution
}

func NewActionService() (*ActionService, error) {
	return &ActionService{}, nil
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

func (as *ActionService) modifyRun(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "ActionService.modifyRun(ctx): %w"

	bashCode := `
	echo $? > minecraft_exitcode
	`

	file, err := os.OpenFile(serverPath+"/"+run, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	defer file.Close()

	_, err = file.WriteString(bashCode)
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

	err = as.modifyRun(ctx)
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

func (as *ActionService) removeExitcode(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "ActionService.removeExitcode(ctx): %w"

	command := "rm"
	args := append(
		make([]string, 0),
		"-f",
		serverPath+"/"+exitcodeFile,
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

	err = as.removeExitcode(ctx)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
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
		_, err := mine_os.ManagedExecCtx(as.aliveCtx, command, args)
		if err != nil {
			logrus.Debugln("get error in managed start:", err)
		}
	}()

	return err
}

func (as *ActionService) Stop(ctx context.Context) error {
	as.stopCtx()
	return nil
}
