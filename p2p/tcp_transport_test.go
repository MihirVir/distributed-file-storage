package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpTrasnport(t *testing.T) {
	listenAdd := ":3001"

	tcpOpts := &TCPTransportOpts{
		ListenAddr: ":3000",
	}
	tr := NewTcpTransport(*tcpOpts)

	assert.Equal(t, tr.ListenAddr, listenAdd)

	// Server
	// tr.Start()
	assert.Nil(t, tr.ListenAndStart())

	select {}
}
