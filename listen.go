package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/zerocruft/capacitor"
	"github.com/zerocruft/flux/debug"
	"github.com/zerocruft/flux/internal"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func listen() {

	router := mux.NewRouter()
	router.HandleFunc("/flux", handleFluxConnection).Methods("GET")
	router.HandleFunc("/control/peer-chat", handlePeerChat).Methods("POST")

	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleFluxConnection(resp http.ResponseWriter, req *http.Request) {

	//TODO check header stuff for validation?

	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	token := newToken() //TODO this needs to be even more unique if part of flux cluster
	responseMsg := capacitor.FluxConnectionResponseToBytes(token)

	// Send the client their token. Signifying an accepted connection
	err = conn.WriteMessage(websocket.TextMessage, responseMsg)
	if err != nil {
		//Client was unable to receive response msg. Close connection
		debug.Log("client token response failed")
		debug.Log(err)
		conn.Close()
	}

	go internal.NewClientConnection(token, conn)

}

func handlePeerChat(response http.ResponseWriter, request *http.Request) {
	//TODO validate payload make sure its a valid FluxMsg
	//TODO valudate auth in header (perhaps a generated token that only the balancer knows)

	defer request.Body.Close()
	msgBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		debug.Log(err)
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	internal.PropogateMsg("NODE_TALK", msgBytes)
	return
}
