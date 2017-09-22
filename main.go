package main

import (
	"bytes"
	"encoding/json"
	"fluxmq/debug"
	"fmt"
	"github.com/zerocruft/capacitor"
	"github.com/zerocruft/flux/cluster"
	"github.com/zerocruft/flux/internal"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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
		if config.Balancer.BalancerAddress != "" {

			ping := capacitor.FluxPing{
				Node: capacitor.FluxNode{
					Address: config.Url + ":" + strconv.Itoa(config.Port) + "/flux",
				},
				NumberOfConnections: internal.NumberOfConnections(),
			}
			pingBytes, err := json.Marshal(ping)
			if err != nil {
				debug.Log("cluster tracker: pingBytes fail")
				debug.Log(err)
				return
			}
			balancerUrl := "http://" + config.Balancer.BalancerAddress + ":" + strconv.Itoa(config.Balancer.BalancerPort) + "/control/cluster/ping"
			resp, err := http.Post(balancerUrl, "application/json", bytes.NewReader(pingBytes))
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
			cluster.SetPeers(pong.Peers)
		}
	}
}
