package main

import (
	"flag"
	"github.com/zerocruft/flux/debug"
	"sync"
)

var (
	flgDebug  bool
	flgConfig string
	flgPort   int
	mainWG    sync.WaitGroup
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
	debug.Init(flgDebug)
}

func parseFlags() {
	flag.BoolVar(&flgDebug, "debug", false, "debug mode")
	flag.StringVar(&flgConfig, "config", "flux.toml", "location of flux config")
	flag.IntVar(&flgPort, "port", 8080, "port")
	flag.Parse()
}
