package main

import (
	"fmt"
	"log"
	"mihir/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store *Store
	quit  chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		store:          NewStore(storeOpts),
		FileServerOpts: opts,
		quit:           make(chan struct{}),
	}
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndStart(); err != nil {
		return err
	}

	fs.loop()

	return nil
}

func (fs *FileServer) loop() {

	defer func() {
		log.Printf("file server stopped due to user quit action")
	}()

	for {
		select {
		case msg := <-fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.quit:
			return
		}
	}
}

func (fs *FileServer) Stop() {
	close(fs.quit)
}
