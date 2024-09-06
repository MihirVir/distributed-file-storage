package p2p

import (
	"fmt"
	"net"
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
	OnPeer        func(Peer) error
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implements Peer Interface, which is used to close the remote connection
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpch     chan RPC
}

// Consume implements the Transport Interface, which will return read-only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpch
}

func NewTcpTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpch:             make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndStart() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

/*
startAcceptLoop is an infinite loop that listens for new connections.
When a connection is detected, it accepts it and creates a new peer.
*/
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handlePeer(conn)
	}
}

func (t *TCPTransport) handlePeer(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	/*
		The purpose of this is to allow the user of the TCPTransport to define what happens when a new peer joins,
		like logging the event, adding the peer to a list, or sending an initial message
	*/
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)

		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()
		t.rpch <- rpc

		fmt.Printf("message: %+v\n", rpc)
	}
}
