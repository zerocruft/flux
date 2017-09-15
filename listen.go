package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/zerocruft/capacitor"
	"github.com/zerocruft/flux/debug"
	"github.com/zerocruft/flux/internal"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func listen() {

	router := mux.NewRouter()
	router.HandleFunc("/test", handleTest).Methods("GET")
	router.HandleFunc("/flux", handleFluxConnection).Methods("GET")
	router.HandleFunc("/control/cluster", handleControlCluster).Methods("POST")

	err := http.ListenAndServe(":"+strconv.Itoa(flgPort), router)
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

func handleControlCluster(response http.ResponseWriter, request *http.Request) {

}

func handleTest(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("HELLO DANG WORLD!!! BOOOM"))
}
