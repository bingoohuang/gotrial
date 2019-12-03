package mutex

type Lock struct {
	ch chan struct{}
}

func NewLock() *Lock {
	mu := &Lock{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

func (m *Lock) Lock() { <-m.ch }
func (m *Lock) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Lock) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *Lock) IsLocked() bool { return len(m.ch) == 0 }

type Lock2 struct {
	ch chan struct{}
}

func NewLock2() *Lock2 {
	mu := &Lock2{make(chan struct{}, 1)}
	return mu
}

func (m *Lock2) Lock() { m.ch <- struct{}{} }
func (m *Lock2) Unlock() {
	select {
	case <-m.ch:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Lock2) TryLock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
	}
	return false
}

func (m *Lock2) IsLocked() bool { return len(m.ch) == 0 }
