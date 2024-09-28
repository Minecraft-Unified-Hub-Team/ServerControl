package mine_state

type State uint

const (
	Alive   State = 0
	Stopped State = 1
	Dead    State = 2
)

func (s *State) IsAlive() bool {
	return *s == 0
}

func (s *State) IsStopped() bool {
	return *s == 1
}

func (s *State) IsDead() bool {
	return *s == 2
}
