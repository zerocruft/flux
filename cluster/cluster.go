package cluster

import (
	"bytes"
	"fluxmq/debug"
	"github.com/zerocruft/capacitor"
	"net/http"
	"sync"
)

var (
	clusterPeers []capacitor.FluxNode
	mutex        sync.Mutex
)

func init() {
	clusterPeers = []capacitor.FluxNode{}
	mutex = sync.Mutex{}
}

func SetPeers(ps []capacitor.FluxNode) {
	mutex.Lock()
	clusterPeers = ps
	mutex.Unlock()
}

func copyOfPeers() (peers []capacitor.FluxNode) {
	peers = []capacitor.FluxNode{}
	mutex.Lock()
	defer mutex.Unlock()
	for _, peer := range clusterPeers {
		peers = append(peers, peer)
	}
	return
}

func PropagateToPeers(msgBytes []byte) {
	for _, peer := range copyOfPeers() {
		go postMsgToPeer(peer, msgBytes)
	}
}

func postMsgToPeer(node capacitor.FluxNode, msgBytes []byte) {
	resp, err := http.Post(node.Address, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		debug.Log(err)
		return
	}

	if resp.StatusCode != 200 {
		debug.Log("Attempting to send msg to peer: status code not 200")
		debug.Log(node.Address)
	}
	return
}
