package p2p

import "net"

type HandshakerFunc func(net.Conn) error

func TCPHandshake(conn net.Conn) error { return nil }
