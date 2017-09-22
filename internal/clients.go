package internal

import (
	"github.com/gorilla/websocket"
	"github.com/zerocruft/flux/cluster"
	"github.com/zerocruft/flux/debug"
	"time"
)

type fluxClientConnection struct {
	token        string
	sendToClient chan []byte
	kill         bool
}

func NewClientConnection(token string, conn *websocket.Conn) {
	fcc := &fluxClientConnection{
		token:        token,
		sendToClient: make(chan []byte, 25),
	}

	go func() {
		for {
			time.Sleep(500 * time.Microsecond)
			if fcc.kill {
				break
			}

			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				debug.Log("read fail: client [" + fcc.token + "]")
				debug.Log(err)
				conn.Close()
				fcc.kill = true
				killClient(fcc.token)
				break
			}

			debug.Log(fcc.token + ": " + string(msgBytes))
			go cluster.PropagateToPeers(msgBytes)
			PropogateMsg(fcc.token, msgBytes)
		}
	}()

	go func() {
		for {
			time.Sleep(500 * time.Microsecond)
			if fcc.kill {
				break
			}

			select {
			case msgBytes := <-fcc.sendToClient:
				err := conn.WriteMessage(websocket.TextMessage, msgBytes)
				if err != nil {
					debug.Log("write fail: client[" + fcc.token + "]")
					debug.Log(err)
					fcc.kill = true
					conn.Close()
					killClient(fcc.token)
					break
				}

			default:
				continue
			}
		}
	}()

	persistClient(fcc)
}
