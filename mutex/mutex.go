package mutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Mutex struct {
	sync.Mutex
}

const (
	mutexLocked = 1 << iota
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
)

// TryLock tries to lock in a fast way.
func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32(m.internalState(), 0, mutexLocked)
}

// Count counts the number the waiters and the owner on the lock.
func (m *Mutex) Count() int {
	state := atomic.LoadInt32(m.internalState())

	return int(state>>mutexWaiterShift + state&mutexLocked)
}

// IsWoken tells if the mutex is in woken state.
func (m *Mutex) IsWoken() bool {
	return atomic.LoadInt32(m.internalState())&mutexWoken == mutexWoken
}

// IsStarving tells if the mutex is in starving state.
func (m *Mutex) IsStarving() bool {
	return m.getState()&mutexStarving == mutexStarving
}

func (m *Mutex) getState() int32 {
	return atomic.LoadInt32(m.internalState())
}

func (m *Mutex) internalState() *int32 {
	return (*int32)(unsafe.Pointer(&m.Mutex))
}
