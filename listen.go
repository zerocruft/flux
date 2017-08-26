package main

import (
	"log"
	"bufio"
	"net"
	"github.com/zerocruft/flux/debug"
	"github.com/zerocruft/flux/capacitor"
	"fmt"
)

func listen() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		debug.Log(err)
		log.Println("Failed to init listener")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			debug.Log(err)
			log.Println("Connection Failure")
		}

		// The client is going to want a client token assigned to them. This first request should be that. If not, bail
		clientRequestBytes, err := bufio.NewReader(conn).ReadBytes('#')
		rawFluxObject, success := capacitor.BytesToFluxObject(clientRequestBytes)
		if !success {
			//TODO log something
			conn.Close()
			return
		}
		fmt.Println(rawFluxObject)
	}
}