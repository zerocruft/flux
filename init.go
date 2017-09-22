package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/zerocruft/flux/debug"
	"os"
	"sync"
)

var (
	flgDebug  bool
	flgConfig string
	mainWG    sync.WaitGroup
	config    FluxConfig
)

func init() {
	parseFlags()
	loadConfig()
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
	flag.Parse()
}

func loadConfig() {
	_, err := toml.DecodeFile(flgConfig, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
