package health

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	mine_ping "github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_ping/ping"
)

type State int

const (
	portField    = "server-port"
	serverPath   = "/server"
	exitcodeFile = "minecraft_exitcode"
)

func fileExists(ctx context.Context, path string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("health.FileExists(ctx, %s)", path) + " %w"

	_, err = os.Stat(path)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func checkMinecraftProcessStarted(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "health.checkMinecraftProcessStarted(ctx): %w"
	var out bytes.Buffer

	cmd := exec.Command("sh", "-c", "ps aux | grep java | grep -v grep")
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	if out.Len() > 0 {
		return err
	}

	return fmt.Errorf(errorFormat, fmt.Errorf("no any java process started"))
}

func getMinecraftExitcode(ctx context.Context, path string) (int, error) {
	var err error = nil
	var errorFormat string = fmt.Sprintf("health.getMinecraftExitcode(ctx, %s)", path) + ": %w"
	var exitcode int = 0

	file, err := os.Open(path)
	if err != nil {
		return -1, fmt.Errorf(errorFormat, err)
	}
	defer file.Close()

	_, err = fmt.Fscanf(file, "%d", &exitcode)
	if err != nil {
		return -1, fmt.Errorf(errorFormat, err)
	}

	return exitcode, err
}

const (
	Alive State = iota
	Stopped
	Dead
)

func (s State) String() string {
	return [...]string{"Alive", "Stopped", "Dead"}[s]
}

func (s State) EnumIndex() int {
	return int(s)
}

type SyncedState struct {
	value State
	mutex sync.Mutex
}

func NewSyncedState(state State) (*SyncedState, error) {
	var err error = nil

	return &SyncedState{
		value: state,
	}, err
}

func (ss *SyncedState) Set(value State) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.value = value
}

func (ss *SyncedState) Get() State {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	return ss.value
}

type HealthService struct {
	syncedState *SyncedState

	aliveCtx context.Context
	stopCtx  context.CancelFunc
}

func NewHealthService() (*HealthService, error) {
	currentState, _ := NewSyncedState(Stopped)
	return &HealthService{
		syncedState: currentState,
	}, nil
}

func (hs *HealthService) Start(ctx context.Context, refreshTime time.Duration) error {
	hs.aliveCtx, hs.stopCtx = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-hs.aliveCtx.Done():
				hs.syncedState.Set(Stopped) // if server stopped correctly with context cancelation => server state is stopped
				return
			default:
				hs.updateState(ctx)
			}
			time.Sleep(refreshTime)
		}
	}()

	return nil
}

func (hs *HealthService) Stop(ctx context.Context) error {
	hs.stopCtx()
	return nil
}

func (hs *HealthService) updateState(ctx context.Context) {
	var err error = nil

	err = fileExists(ctx, serverPath+"/"+exitcodeFile)
	if err != nil { // file does not exist => server stopped or alive
		err = checkMinecraftProcessStarted(ctx)
		if err != nil {
			hs.syncedState.Set(Stopped)
			return
		}
		hs.syncedState.Set(Alive)
		return
	}

	exitcode, err := getMinecraftExitcode(ctx, serverPath+"/"+exitcodeFile)
	if err == nil && exitcode == 0 {
		hs.syncedState.Set(Stopped)
		return
	}

	hs.syncedState.Set(Dead)
}

func (hs *HealthService) Ping(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "HealthService.Ping(ctx): %w"

	_, status, _ := mine_ping.StaticPing(ctx)
	if status != -1 {
		return fmt.Errorf(errorFormat, err)
	}

	return err
}

func (hs *HealthService) GetState(ctx context.Context) api.State {
	state := hs.syncedState.Get()
	switch state {
	case Alive:
		return api.State_Alive
	case Stopped:
		return api.State_Stopped
	default:
		return api.State_Dead
	}
}
