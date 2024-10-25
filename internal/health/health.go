package health

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	mine_ping "github.com/Minecraft-Unified-Hub-Team/ServerControl/utils/mine_ping/ping"
)

type State int

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

type HealthService struct {
	syncedState *SyncedState
}

func NewHealthService() (*HealthService, error) {
	currentState, _ := NewSyncedState(Stopped)
	return &HealthService{
		syncedState: currentState,
	}, nil
}

func (hs *HealthService) Start(ctx context.Context, refreshTime time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				hs.updateState(ctx)
			}
			time.Sleep(refreshTime)
		}
	}()
}

func (hs *HealthService) updateState(ctx context.Context) error {
	var err error = nil
	var errorFormat string = "HealthService.updateState(ctx): %w"

	hs.syncedState.mutex.Lock()
	defer hs.syncedState.mutex.Unlock()

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	_, status, err := mine_ping.Ping(hostname)
}

func (hs *HealthService) Ping(ctx context.Context) error {
	var err error = nil

	return err
}

func (hs *HealthService) GetState(ctx context.Context) api.State {

}
