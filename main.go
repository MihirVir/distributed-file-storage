package main

import (
	"log"

	"github.com/anthdm/foreverstore/p2p"
)

func main() {
	tr := p2p.NewTcpTransport(":3000")

	if err := tr.ListenAndStart(); err != nil {
		log.Fatal(err)
	}

	select {}
}
