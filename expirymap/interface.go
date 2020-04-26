package expirymap

import (
	"container/list"
	"time"
)

type ExpiryMap interface {
	Get(k interface{}, updateExpiry bool) interface{}
	Set(k interface{}, v interface{})
	Delete(k interface{})
}

func New(maxAge time.Duration) ExpiryMap {
	r := &expiryMap{
		maxAge:      maxAge,
		data:        make(map[interface{}]*expiryMapItem),
		expiryQueue: list.New(),
	}
	return r
}
