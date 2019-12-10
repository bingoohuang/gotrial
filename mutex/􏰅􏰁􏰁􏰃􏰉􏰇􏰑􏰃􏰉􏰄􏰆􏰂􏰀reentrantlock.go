package mutex

import (
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

type RecursiveMutex struct {
	sync.Mutex
	owner     int64
	recursion int32
}

func (m *RecursiveMutex) Lock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}

	m.Mutex.Lock()

	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) != gid {
		panic("unlock should be in same goroutine with its previous lock g")
	}

	m.recursion--
	if m.recursion != 0 {
		return
	}

	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

// TokenRecursionMutex ...
type TokenRecursionMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

func (m *TokenRecursionMutex) Lock(t int64) {
	if atomic.LoadInt64(&m.token) == t {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.token, t)
	m.recursion = 1
}

func (m *TokenRecursionMutex) Unlock(t int64) {
	if atomic.LoadInt64(&m.token) != t {
		panic("unlock should be in same token with its previous lock g")
	}

	m.recursion--
	if m.recursion != 0 {
		return
	}

	atomic.StoreInt64(&m.token, 0)
	m.Mutex.Unlock()
}
