package hit

import (
	"container/list"
	"time"
)

const bucketCount = 100

type Counter interface {
	Get(k interface{}) int64
	Add(k interface{}, count int64)
}

func NewCounter(duration time.Duration) Counter {
	if duration <= 0 {
		panic("Invalid duration for hit counter")
	}

	r := &counterMap{
		duration:    duration,
		delta:       duration / time.Duration(bucketCount),
		data:        make(map[interface{}]*counterMapItem),
		expiryQueue: list.New(),
	}
	return r
}
