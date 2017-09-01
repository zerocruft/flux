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
	nano := strconv.Itoa(time.Now().Nanosecond())
	time.Sleep(2 * time.Nanosecond)
	tokenMutex.Unlock()
	return nano

}
