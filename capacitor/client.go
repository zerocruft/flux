package capacitor

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"net"
)


// Returns writeChannel, readChannel, error
func NewClient(req FluxConnectParameters) (FluxClient, error) {
	var clientConnection net.Conn
	clientConnection, err := dial(req.FluxAddress)
	if err != nil {
		return FluxClient{}, err
	}

	//Ask Flux service to respond with a capacitor token
	connectRequestBytes := fluxConnectionRequestToBytes()
	_, err = clientConnection.Write(connectRequestBytes)
	if err != nil {
		clientConnection.Close()
		return FluxClient{}, err
	}

	responseBytes, err := bufio.NewReader(clientConnection).ReadBytes('#')
	fluxServiceResponse, success := BytesToFluxObject(bytes.TrimRight(responseBytes, "#"))
	if !success {
		clientConnection.Close()
		return FluxClient{}, errors.New("Internal Flux Object Error. Invalid msg format")
	}
	//Check to make sure flux object is of connection response type
	if fluxServiceResponse.GetType() != FLUX_TYPE_CONNECTION_RESPONSE {
		clientConnection.Close()
		log.Println(fluxServiceResponse)//TODO do I need to do this?
		return FluxClient{}, errors.New("Flux Service failed to respond")
	}

	//Now.. subscribe to channels
	for _, channel := range req.Topics {
		channelSubscribeReq := fluxTopicSubscriptionRequestToBytes(fluxServiceResponse.GetClientToken(), channel)
		clientConnection.Write(channelSubscribeReq)
	}

	clientWriteChannel := make(chan FluxMessage, 3)
	clientReadChannel := make(chan FluxMessage, 3)
	connections[fluxServiceResponse.GetClientToken()] = clientConnection
	clientWriteChannels[fluxServiceResponse.GetClientToken()] = clientWriteChannel
	clientReadChannels[fluxServiceResponse.GetClientToken()] = clientReadChannel

	// Write Channel
	go func() {
		for {
			newMsg := <-clientWriteChannel
			flxMsgBytes := fluxMessageToBytes(fluxServiceResponse.GetClientToken(), newMsg)
			_, err := clientConnection.Write(flxMsgBytes)
			if err != nil {
				//cLOSE CONNECTION, REMOVE CLIENT, etc...
			}
		}
	}()

	// Read Channel
	go func() {
		for {
			newMsgBytes, err := bufio.NewReader(clientConnection).ReadBytes('#')
			if err != nil {
				//cLOSE CONNECTION, REMOVE CLIENT, etc...
			}
			clientReadChannel <- bytesToFluxMessage(newMsgBytes)

		}
	}()

	return FluxClient{
		clientToken: fluxServiceResponse.GetClientToken(),
		send: clientWriteChannel,
		receive: clientReadChannel,
	}, nil
}



type FluxClient struct {
	clientToken string
	send    chan FluxMessage
	receive chan FluxMessage
}
