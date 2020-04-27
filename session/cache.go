package session

import (
	"time"

	"github.com/golangext/datastructures/expirymap"
)

type cache struct {
	data expirymap.ExpiryMap
}

type Cache interface {
	NewSession() Session
	Get(id string) Session
	Resurect(session Session)
}

func NewCache(expiry time.Duration) Cache {
	c := &cache{expirymap.New(expiry)}
	return c
}

func (c *cache) NewSession() Session {
	sess := NewSession()
	c.data.Set(sess.ID(), sess)
	return sess
}

func (c *cache) Resurect(session Session) {
	c.data.Set(session.ID(), session)
}

func Resurect(id string) Session {
	dat := &session{id: id, data: make(map[interface{}]interface{})}
	return dat
}

func (c *cache) Get(id string) Session {
	dat := c.data.Get(id, true)
	if dat == nil {
		return nil
	}

	return dat.(*session)
}
