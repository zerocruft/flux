package main

import (
	"flag"
	"github.com/zerocruft/flux/debug"
	"sync"
)

var (
	flgDebug bool
	mainWG   sync.WaitGroup
)

func init() {
	parseFlags()
	initDebug()
	initWaitGroup()
}

func initWaitGroup() {
	mainWG = sync.WaitGroup{}
	mainWG.Add(1)
}
func initDebug() {
	debug.Init(func() bool {
		return flgDebug
	})
}

func parseFlags() {
	flag.BoolVar(&flgDebug, "debug", false, "debug mode")
	flag.Parse()
}
