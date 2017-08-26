package capacitor

import "net"

var (
	connections         map[string]net.Conn
	clientWriteChannels map[string]chan FluxMessage
	clientReadChannels  map[string]chan FluxMessage
)

func init() {
	connections = map[string]net.Conn{}
	clientWriteChannels = map[string]chan FluxMessage{}
	clientReadChannels = map[string]chan FluxMessage{}
}
