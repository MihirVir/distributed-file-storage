package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpTrasnport(t *testing.T) {
	tcpOpts := &TCPTransportOpts{
		ListenAddr:    ":3001",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTcpTransport(*tcpOpts)

	assert.Equal(t, tr.ListenAddr, ":3001")

	// Server
	// tr.Start()
	assert.Nil(t, tr.ListenAndStart())
}
