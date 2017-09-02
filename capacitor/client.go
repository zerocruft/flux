package capacitor

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zerocruft/flux/debug"
	"log"
	"strconv"
)

// Returns writeChannel, readChannel, error
func NewClient(req FluxConnectParameters) (FluxClient, error) {
	debug.Init(req.Debug)

	fluxUrl := "ws://" + req.FluxAddress + ":" + strconv.Itoa(req.FluxPort) + "/flux"
	conn, _, err := websocket.DefaultDialer.Dial(fluxUrl, nil)
	if err != nil {
		fmt.Println(err)
		return FluxClient{}, err
	}

	mt, payload, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		debug.Log(mt)
		debug.Log(err)
		return FluxClient{}, err
	}
	debug.Log("response: " + string(payload))

	fluxServiceResponse, success := bytesToFluxObject(payload)
	if !success {
		conn.Close()
		return FluxClient{}, errors.New("Internal Flux Object Error. Invalid msg format")
	}

	//Check to make sure flux object is of connection response type
	if fluxServiceResponse.GetType() != FLUX_TYPE_CONNECTION_RESPONSE {
		conn.Close()
		log.Println(fluxServiceResponse) //TODO do I need to do this?
		return FluxClient{}, errors.New("Flux Service failed to respond")
	}

	//Now.. subscribe to channels
	for _, channel := range req.Topics {
		channelSubscribeReq := fluxTopicSubscriptionRequestToBytes(fluxServiceResponse.GetClientToken(), channel)
		conn.WriteMessage(websocket.TextMessage, channelSubscribeReq)
	}

	clientWriteChannel := make(chan FluxMessage, 25)
	clientReadChannel := make(chan FluxMessage, 25)
	connections[fluxServiceResponse.GetClientToken()] = conn
	clientWriteChannels[fluxServiceResponse.GetClientToken()] = clientWriteChannel
	clientReadChannels[fluxServiceResponse.GetClientToken()] = clientReadChannel

	// Write Channel
	go func() {
		for {
			newMsg := <-clientWriteChannel
			debug.Log(newMsg)
			flxMsgBytes := fluxMessageToBytes(fluxServiceResponse.GetClientToken(), newMsg)
			err := conn.WriteMessage(websocket.TextMessage, flxMsgBytes)
			if err != nil {
				debug.Log(err)
				conn.Close()
			}
		}
	}()

	// Read Channel
	go func() {
		for {
			mt, payload, err := conn.ReadMessage()
			if err != nil {
				debug.Log(err)
				debug.Log(mt)
				conn.Close()
			}
			clientReadChannel <- bytesToFluxMessage(payload)
		}
	}()

	return FluxClient{
		clientToken: fluxServiceResponse.GetClientToken(),
		send:        &clientWriteChannel,
		receive:     &clientReadChannel,
	}, nil
}

type FluxClient struct {
	clientToken string
	send        *chan FluxMessage
	receive     *chan FluxMessage
}

func (fc *FluxClient) Send() chan FluxMessage {
	return *fc.send
}

func (fc *FluxClient) Receive() chan FluxMessage {
	return *fc.receive
}
