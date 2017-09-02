package main

import (
	"strconv"
	"sync"
	"time"
)

var (
	tokenMutex sync.Mutex
)

func newToken() string {
	tokenMutex.Lock()
	nano := strconv.FormatInt(time.Now().UnixNano(), 10)
	time.Sleep(2 * time.Nanosecond)
	tokenMutex.Unlock()
	return nano[9:]

}
