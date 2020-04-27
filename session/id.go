package session

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	b "encoding/binary"
	"sync"
	"sync/atomic"

	"github.com/golangext/datastructures/mutex"
)

const reallocRatio = 10000

var reallocMutex sync.RWMutex
var seqNo = int64(0)
var randData []byte

func init() {
	randData = generateRandData()
}

func generateRandData() []byte {
	length := 200
	ret := make([]byte, length, length)
	rand.Reader.Read(ret)
	return ret
}

func NewID() string {
	lock := mutex.Reference(&reallocMutex)
	defer lock.Unlock()
	generator := sha512.New()
	intBuf := make([]byte, 8, 8)
	nextNo := atomic.AddInt64(&seqNo, 1)
	if nextNo%reallocRatio != 0 {
		lock.LockShared()
	} else {
		lock.LockExclusive()
		randData = generateRandData()
	}
	b.PutVarint(intBuf, nextNo)
	generator.Write(intBuf)
	generator.Write(randData)
	finalSum := generator.Sum(nil)
	encodedLen := base64.StdEncoding.EncodedLen(len(finalSum))
	finalByteArray := make([]byte, encodedLen, encodedLen)
	base64.StdEncoding.Encode(finalByteArray, finalSum)
	return string(finalByteArray)
}
