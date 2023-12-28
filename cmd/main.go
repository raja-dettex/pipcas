package main

import (
	"log"
	"os"

	"github.com/raja-dettex/pipcas/p2p"
	"github.com/raja-dettex/pipcas/server"
	"github.com/raja-dettex/pipcas/storage"
)

var (
	addr   = os.Getenv("LISTEN_ADDR")
	lbAddr = os.Getenv("LB_ADDR")
)

func main() {
	opts := p2p.TransportOPts(addr, lbAddr)
	storageOpts := storage.StorageOpts{
		PathTransform: storage.CASTransformFunc,
	}
	sOpts := server.ServerOpts{TransportOpts: opts, StorageOpts: storageOpts}
	server := server.NewStorageServer(sOpts)
	go server.HandleConsume()
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
	// time.Sleep(time.Second * 1)
	// go sendWriteBytes()
	// time.Sleep(time.Second * 1)
	// go sendReadBytes()
	select {}
}
