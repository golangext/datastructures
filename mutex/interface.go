package mutex

import "sync"

type Lock interface {
	LockShared()
	LockExclusive()
	Upgrade() bool
	Unlock()
}

func Reference(m *sync.RWMutex) Lock {
	return &rwlock{mux: m}
}

func LockShared(m *sync.RWMutex) Lock {
	r := Reference(m)
	r.LockShared()
	return r
}

func LockExclusive(m *sync.RWMutex) Lock {
	r := Reference(m)
	r.LockExclusive()
	return r
}
