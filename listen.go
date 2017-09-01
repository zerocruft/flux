package main

import (
	"bufio"
	"github.com/zerocruft/flux/capacitor"
	"github.com/zerocruft/flux/debug"
	"github.com/zerocruft/flux/internal"
	"log"
	"net"
	"strings"
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
			continue
		}

		// The client is going to want a client token assigned to them. This first request should be that. If not, bail
		clientRequestBytes, err := bufio.NewReader(conn).ReadBytes('#')
		fluxMsgSections := strings.Split(string(clientRequestBytes), ":")
		if len(fluxMsgSections) != 4 || fluxMsgSections[0] != capacitor.FLUX_TYPE_CONNECTION_REQUEST {
			//TODO log something
			conn.Close()
			continue
		}

		newClientToken := newToken()
		responseMsg := capacitor.FluxConnectionResponseToBytes(newClientToken)
		_, err = conn.Write(responseMsg)
		if err != nil {
			debug.Log(err)
			conn.Close()
			continue
		}

		go internal.NewClientConnection(newClientToken, conn)
	}
}
