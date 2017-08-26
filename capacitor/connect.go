package capacitor

import "net"

type FluxConnectParameters struct {
	FluxAddress string
	Topics      []string
}

func dial(address string) (net.Conn, error) {

	return net.Dial("tcp", address)
}
