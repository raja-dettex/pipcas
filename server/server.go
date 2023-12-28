package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/raja-dettex/pipcas/p2p"
	"github.com/raja-dettex/pipcas/storage"
)

type Server interface {
	Start() error
}

type ServerOpts struct {
	TransportOpts p2p.TransportOpts
	StorageOpts   storage.StorageOpts
}

type StorageServer struct {
	transport p2p.Transport
	storage   storage.Storage
}

func OnPeer(peer p2p.Peer) error {
	return nil
}

func NewStorageServer(opts ServerOpts) *StorageServer {
	return &StorageServer{
		transport: p2p.NewTCPTransport(opts.TransportOpts, OnPeer),
		storage:   *storage.NewStorage(opts.StorageOpts),
	}
}

func (ss *StorageServer) Start() error {
	fmt.Println("lb ", ss.transport.LbAddress())
	conn, err := net.Dial("tcp", ss.transport.LbAddress())
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(fmt.Sprintf("add %s %s", strings.Split(ss.transport.Address(), ":")[0], strings.Split(ss.transport.Address(), ":")[1])))
	if err != nil {
		return err
	}
	if err := ss.transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}

func (ss *StorageServer) HandleConsume() {
	for rpc := range ss.transport.Consume() {
		fmt.Println(rpc)
		switch rpc.Header {
		case p2p.RPCWriteHeader:
			res, err := ss.storage.WriteStream(string(rpc.Key), bytes.NewReader(rpc.Paylaod))
			if err != nil {
				fmt.Println(err)
				continue
			}
			conn, ok := ss.transport.Conn(rpc.From)
			if !ok {
				fmt.Printf("conn %v does not exist \n", conn)
				continue
			}
			fmt.Println("here", res)
			_, err = (*conn).Write([]byte(res))
			if err != nil {
				fmt.Println("error occured while writing to conn", err)
			}
			ss.transport.Dump(rpc.From)
			(*conn).Close()
		case p2p.RPCReadHeader:
			r, err := ss.storage.Read(string(rpc.Key))
			if err != nil {
				fmt.Println(err)
				continue
			}
			bytes, err := ioutil.ReadAll(r)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("file content ", string(bytes))
			conn, ok := ss.transport.Conn(rpc.From)
			if !ok {
				fmt.Printf("conn %v does not exist \n", conn)
				continue
			}
			_, err = (*conn).Write([]byte(string(bytes)))
			if err != nil {
				fmt.Println("error occured while writing to conn", err)
			}
			ss.transport.Dump(rpc.From)
			(*conn).Close()
		case p2p.RPCRemoveHeader:
			if err := ss.storage.Delete(string(rpc.Key)); err != nil {
				fmt.Println(err)
				continue
			}
			conn, ok := ss.transport.Conn(rpc.From)
			if !ok {
				fmt.Printf("conn %v does not exist \n", conn)
				continue
			}
			_, err := (*conn).Write([]byte("Deleted"))
			if err != nil {
				fmt.Println("error occured while writing to conn", err)
			}
			ss.transport.Dump(rpc.From)
			(*conn).Close()
		}

	}
}
