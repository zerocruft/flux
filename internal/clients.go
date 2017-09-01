package internal

import (
	"bufio"
	"github.com/zerocruft/flux/debug"
	"net"
	"time"
)

type fluxClientConnection struct {
	token      string
	connection net.Conn
	created    time.Time
	lastPinged time.Time
	stack      chan []byte
	kill       bool
}

func NewClientConnection(token string, conn net.Conn) {
	fcc := &fluxClientConnection{
		token:      token,
		connection: conn,
		created:    time.Now(),
		lastPinged: time.Now(),
		stack:      make(chan []byte, 10),
	}

	go func() {
		for {
			time.Sleep(10 * time.Microsecond)
			if fcc.kill {
				break
			}

			incomingBytes, err := bufio.NewReader(fcc.connection).ReadBytes('#')
			if err != nil {
				debug.Log("read fail: client [" + fcc.token + "]")
				debug.Log(err)
				fcc.kill = true
				break
			}

			propogateMsg(fcc.token, incomingBytes)
		}
	}()

	go func() {
		for {
			if fcc.kill {
				break
			}

			_, err := fcc.connection.Write(<-fcc.stack)
			if err != nil {
				debug.Log("write fail: client[" + fcc.token + "]")
				debug.Log(err)
				fcc.kill = true
				break
			}
		}

	}()

	persistClient(fcc)
}
