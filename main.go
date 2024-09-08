package main

import (
	"log"
	"mihir/p2p"
	"time"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: on peer func
	}

	tcpTransport := p2p.NewTcpTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransfromFunc,
		Transport:         tcpTransport,
	}
	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
