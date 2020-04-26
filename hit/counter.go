package hit

import (
	"container/list"
	"sync"
	"time"

	"github.com/golangext/datastructures/mutex"
)

type counterMap struct {
	mux         sync.RWMutex
	data        map[interface{}]*counterMapItem
	expiryQueue *list.List
	duration    time.Duration
	delta       time.Duration
}

type counterMapItem struct {
	bucketTime         time.Time
	expiryQueueElement *list.Element
	k                  interface{}
	current            int64
	delta              int64
}

func item(e *list.Element) *counterMapItem {
	if e == nil {
		return nil
	}

	return e.Value.(*counterMapItem)
}

func (m *counterMap) _delete(e *counterMapItem) {
	m.expiryQueue.Remove(e.expiryQueueElement)
	if cur, ok := m.data[e.k]; ok {
		cur.current -= e.delta
		if cur == e {
			delete(m.data, e.k)
		}
	}
}

func (m *counterMap) run(now time.Time) {
	watermark := now.Add(-m.duration)
	// anything before watermark needs to go
	toExpire := item(m.expiryQueue.Front())
	for toExpire != nil && toExpire.bucketTime.Before(watermark) {
		m._delete(toExpire)
	}
}

func (m *counterMap) Get(k interface{}) int64 {
	defer mutex.LockExclusive(&m.mux).Unlock()
	now := time.Now()
	m.run(now)

	if item, ok := m.data[k]; ok {
		return item.current
	}

	return 0
}

func (m *counterMap) Add(k interface{}, count int64) {
	defer mutex.LockExclusive(&m.mux).Unlock()
	now := time.Now()
	m.run(now)

	item, ok := m.data[k]

	if ok && item.bucketTime.Add(m.delta).After(now) {
		item.delta += count
		item.current += count
		return
	}

	// else alloc a new bucket
	newItem := &counterMapItem{bucketTime: now, k: k, delta: count, current: count}
	m.data[k] = newItem
	element := m.expiryQueue.PushBack(newItem)
	newItem.expiryQueueElement = element

	if ok {
		newItem.current += item.current
	}
}

func (m *counterMap) Delete(k interface{}) {
	defer mutex.LockExclusive(&m.mux).Unlock()
	now := time.Now()
	m.run(now)

	if item, ok := m.data[k]; ok {
		m._delete(item)
	}
}
