package main

import (
	"time"
)

func main() {

	go listen()
	go ticker()
	mainWG.Wait()

}

func ticker() {

	for {
		time.Sleep(1 * time.Second)
	}
}
