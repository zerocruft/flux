package internal

import (
	"fmt"
	"sync"
	"time"
)

var (
	fccs        map[string]*fluxClientConnection
	topics      map[string][]string
	fccMutex    sync.Mutex
	topicsMutex sync.Mutex
)

func init() {
	fccs = map[string]*fluxClientConnection{}
	topics = map[string][]string{}
	fccMutex = sync.Mutex{}
	topicsMutex = sync.Mutex{}
	go stateTicker()
}

func stateTicker() {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println(len(fccs))
		fmt.Println(topics)
	}
}
