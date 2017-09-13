package main

import (
	"time"
	"fmt"
)

func main() {
	fmt.Println(config)
	go listen()
	go ticker()
	mainWG.Wait()

}

func ticker() {

	for {
		time.Sleep(1 * time.Second)
	}
}
