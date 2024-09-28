package action

import (
	"context"
	"fmt"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_os"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_state"
)

const (
	cd         = "cd"
	run        = "run.sh"
	serverPath = "/server"
)

func NewActionService() (*ActionService, error) {
	currentState, _ := mine_state.NewState(mine_state.Alive) // TODO set mine_state.Stopped here when install will be complited
	return &ActionService{
		State: currentState,
	}, nil
}

type ActionService struct {
	AliveCtx context.Context    // context that continues until server is stopped or dead
	stopCtx  context.CancelFunc // function that cancels server binary execution

	State *mine_state.State // channel that stores state of server
}

func (as *ActionService) Start(ctx context.Context) error {
	/* check that server has not been already started */
	if as.State.IsAlive() {
		return nil
		// return fmt.Errorf("server has been already started") // TODO verify that we use fmt.Errorf for creating errors
	}

	/* create aliveness context for server run */
	as.AliveCtx, as.stopCtx = context.WithCancel(context.Background())

	/* prepare command and arguments */
	command := "/bin/bash"
	args := append(
		make([]string, 0),
		"-c",
		fmt.Sprint(cd, serverPath, "&&", run),
	)

	/* start server in goroutine */
	go func() {
		as.State.Set(mine_state.Alive)
		status, _ := mine_os.ManagedExecCtx(as.AliveCtx, command, args)
		if status == 1 {
			as.State.Set(mine_state.Stopped)
		} else {
			as.State.Set(mine_state.Dead)
		}
	}()

	/* always okay */
	return nil
}

func (as *ActionService) Stop(ctx context.Context) error {
	/* call cancel context function */
	as.stopCtx()

	/* always okay */
	return nil
}
