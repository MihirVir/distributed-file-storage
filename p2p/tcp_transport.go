package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// conn is underlying connection of the peer
	conn net.Conn
	// if we dial and retrieve a connection -> outbound == true
	// if we accept and retrieve a connection -> outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	mu            sync.RWMutex
	peers         map[net.Addr]Peer
}

func NewTcpTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndStart() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)

	if err != nil {
		return err
	}

	// running startAcceptLoop to infinitely loop thru to check if any incoming connection
	// or handlePeer
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		// creating a new peer and accepting the incoming connection
		// infinite loop cuz we want to check if any other peer is trying
		// to connect
		go t.handlePeer(conn)
	}
}

func (t *TCPTransport) handlePeer(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("new incoming connection %v", peer)
}
