package mine_state

import (
	"sync"
)

const (
	Alive   = 0
	Stopped = 1
	Dead    = 2
)

type State struct {
	value int
	mutex sync.Mutex
}

func NewState(state int) (*State, error) {
	return &State{
		value: Stopped,
	}, nil
}

func (s *State) IsAlive() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value == Alive
}

func (s *State) IsStopped() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value == Stopped
}

func (s *State) IsDead() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value == Dead
}

func (s *State) Set(state int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = state
}
