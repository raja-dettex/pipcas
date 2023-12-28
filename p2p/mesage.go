package p2p

import "net"

type RPCHeader byte

const (
	RPCWriteHeader  RPCHeader = 0x01
	RPCReadHeader   RPCHeader = 0x02
	RPCRemoveHeader RPCHeader = 0x03
)

type RPCMessage struct {
	Header  RPCHeader
	From    net.Addr
	Key     []byte
	Paylaod []byte
}
