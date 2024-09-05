package p2p

import "net"

// Message holds any arbitrary data that's being sent over
// each transport between two nodes in the network
type Message struct {
	// which node the msg came from
	From net.Addr
	// body of the data
	Payload []byte
}
