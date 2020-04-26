package mutex

import "sync"

type rwlock struct {
	rlocked bool
	wlocked bool
	mux     *sync.RWMutex
}

func (m *rwlock) LockShared() {
	if m.rlocked || m.wlocked {
		return
	}

	m.rlocked = true
	m.mux.RLock()
}

func (m *rwlock) LockExclusive() {
	if m.rlocked {
		panic("Can not wlock an rlocked mutex")
	}

	if m.wlocked {
		return
	}

	m.wlocked = true
	m.mux.Lock()
}

func (m *rwlock) Upgrade() bool {
	if m.wlocked {
		return false
	}

	if !m.rlocked {
		panic("Can not upgrade unlocked mutex")
	}

	m.rlocked = false
	m.mux.RUnlock()
	m.wlocked = true
	m.mux.Lock()

	return true
}

func (m *rwlock) Unlock() {
	if m.rlocked {
		m.rlocked = false
		m.mux.RUnlock()
	}

	if m.wlocked {
		m.wlocked = false
		m.mux.Unlock()
	}
}
