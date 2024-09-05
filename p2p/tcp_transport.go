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
type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTcpTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndStart() error {
	var err error
	// trying to establish a TCP connection
	t.listener, err = net.Listen("tcp", t.ListenAddr)

	// if there's an error establishing a TCP connection
	// we will throw the error
	if err != nil {
		return err
	}

	// running startAcceptLoop to infinitely loop thru to check if any incoming connection
	// or handlePeer
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	// after establishg we are looping thru (infinite loop) the incoming connections
	// and if there are any incoming connections we accept them and create a peer out of it.
	// Better explanation probs
	// SERVER (RUNNING) <- peer2 initiates connection
	// peer2 gets accepted
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
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}

		fmt.Printf("message: %+v\n", msg)
	}
}
