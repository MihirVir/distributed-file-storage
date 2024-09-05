package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpTrasnport(t *testing.T) {
	listenAdd := ":3001"
	tr := NewTcpTransport(listenAdd)

	assert.Equal(t, tr.listenAddress, listenAdd)

	// Server
	// tr.Start()
	assert.Nil(t, tr.ListenAndStart())

	select {}
}
