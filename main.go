package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/zerocruft/capacitor"
	"github.com/zerocruft/flux/cluster"
	"github.com/zerocruft/flux/debug"
	"github.com/zerocruft/flux/internal"
)

func main() {
	fmt.Println(config)
	go listen()
	go clusterTracker()
	mainWG.Wait()

}

func clusterTracker() {

	for {
		time.Sleep(10 * time.Second)
		cluster.RegisterSelf(config.Iam)
		if config.Balancer.BalancerAddress != "" {

			ping := capacitor.FluxPing{
				Node: capacitor.FluxNode{
					ClientEndpoint: config.Url + ":" + strconv.Itoa(config.Port) + "/flux",
					Name:           config.Iam,
					PeerEndpoint:   config.Url + ":" + strconv.Itoa(config.Port) + "/control/peer-chat",
				},
				NumberOfConnections: internal.NumberOfConnections(),
			}
			pingBytes, err := json.Marshal(ping)
			if err != nil {
				debug.Log("cluster tracker: pingBytes fail")
				debug.Log(err)
				return
			}
			balancerURL := "http://" + config.Balancer.BalancerAddress + ":" + strconv.Itoa(config.Balancer.BalancerPort) + "/control/cluster/ping"
			resp, err := http.Post(balancerURL, "application/json", bytes.NewReader(pingBytes))
			if err != nil {
				debug.Log("cluster tracker: pong response failed")
				debug.Log(err)
				return
			}

			defer resp.Body.Close()
			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				debug.Log("cluster tracker: response body failed")
				debug.Log(err)
				return
			}

			pong := capacitor.FluxPong{}
			err = json.Unmarshal(respBytes, &pong)
			if err != nil {
				debug.Log(err)
			}
			debug.Log(pong)
			cluster.SetPeers(pong.Peers)
		}
	}
}
