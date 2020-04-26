package expirymap

import (
	"container/list"
	"sync"
	"time"

	"github.com/golangext/datastructures/mutex"
)

type expiryMap struct {
	mux         sync.RWMutex
	data        map[interface{}]*expiryMapItem
	expiryQueue *list.List
	maxAge      time.Duration
}

type expiryMapItem struct {
	updatedTime        time.Time
	expiryQueueElement *list.Element
	k                  interface{}
	v                  interface{}
}

func item(e *list.Element) *expiryMapItem {
	if e == nil {
		return nil
	}

	return e.Value.(*expiryMapItem)
}

func (m *expiryMap) _delete(e *expiryMapItem) {
	delete(m.data, e.k)
	m.expiryQueue.Remove(e.expiryQueueElement)
}

func (m *expiryMap) run(now time.Time) {
	watermark := now.Add(-m.maxAge)
	// anything before watermark needs to go
	toExpire := item(m.expiryQueue.Front())
	for toExpire != nil && toExpire.updatedTime.Before(watermark) {
		m._delete(toExpire)
	}
}

func (m *expiryMap) Get(k interface{}, updateExpiry bool) interface{} {
	defer mutex.LockExclusive(&m.mux).Unlock()
	now := time.Now()
	m.run(now)

	if item, ok := m.data[k]; ok {
		if updateExpiry {
			m.expiryQueue.MoveToBack(item.expiryQueueElement)
			item.updatedTime = now
		}
		return item.v
	}

	return nil
}

func (m *expiryMap) Set(k interface{}, v interface{}) {
	defer mutex.LockExclusive(&m.mux).Unlock()
	now := time.Now()
	m.run(now)

	if item, ok := m.data[k]; ok {
		m._delete(item)
	}

	item := &expiryMapItem{updatedTime: now, k: k, v: v}
	m.data[k] = item
	element := m.expiryQueue.PushBack(item)
	item.expiryQueueElement = element
}

func (m *expiryMap) Delete(k interface{}) {
	defer mutex.LockExclusive(&m.mux).Unlock()
	now := time.Now()
	m.run(now)

	if item, ok := m.data[k]; ok {
		m._delete(item)
	}
}
