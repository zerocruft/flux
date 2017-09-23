package cluster

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/zerocruft/capacitor"
	"github.com/zerocruft/flux/debug"
)

var (
	clusterPeers []capacitor.FluxNode
	mutex        sync.Mutex
	iam          string
)

func init() {
	clusterPeers = []capacitor.FluxNode{}
	mutex = sync.Mutex{}
}

func RegisterSelf(name string) {
	iam = name
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
		if peer.Name != iam {
			postMsgToPeer(peer, msgBytes)
		}
	}
}

func postMsgToPeer(node capacitor.FluxNode, msgBytes []byte) {
	c := http.DefaultClient
	r, err := http.NewRequest(http.MethodPost, "http://"+node.PeerEndpoint, bytes.NewReader(msgBytes))
	r.Close = true
	resp, err := c.Do(r)
	debug.Log(r)
	if err != nil {
		debug.Log("DEBUG: postMsgToPeer")
		debug.Log(node.ClientEndpoint)
		debug.Log(err)
		return
	}

	if resp.StatusCode != 200 {
		debug.Log("Attempting to send msg to peer: status code not 200")
		debug.Log(node.ClientEndpoint)
	}
	return
}
