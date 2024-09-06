package main

import (
	"fmt"
	"log"
	"mihir/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {
	tr := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}

	t := p2p.NewTcpTransport(tr)
	if err := t.ListenAndStart(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg := <-t.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	select {}
}
