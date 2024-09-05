package main

import (
	"log"

	"github.com/anthdm/foreverstore/p2p"
)

func main() {
	tr := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	t := p2p.NewTcpTransport(tr)
	if err := t.ListenAndStart(); err != nil {
		log.Fatal(err)
	}

	select {}
}
