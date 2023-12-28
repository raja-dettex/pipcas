package p2p

import "net"

type Peer interface {
	Close() error
}

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPCMessage
	Conn(net.Addr) (*net.Conn, bool)
	Dump(net.Addr)
	Address() string
	LbAddress() string
}
