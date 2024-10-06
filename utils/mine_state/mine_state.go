package mine_state

import (
	"sync"
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

func (s *SyncedState) IsAlive() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value == Alive
}

func (s *SyncedState) IsStopped() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value == Stopped
}

func (s *SyncedState) IsDead() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value == Dead
}

func (s *SyncedState) Set(state State) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = state
}

func (s *SyncedState) State() State {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value
}
