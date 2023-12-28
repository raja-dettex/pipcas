package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (peer *TCPPeer) Close() error {
	if err := peer.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (p *TCPPeer) String() string {
	return fmt.Sprintf("TCPPeer{connection : %d, outbound : %v}", p.conn, p.outbound)
}

type TransportOpts struct {
	listenaddr string
	lbAddr     string
	listener   net.Listener
}

func TransportOPts(addr string, lbAddr string) TransportOpts {
	return TransportOpts{listenaddr: addr, lbAddr: lbAddr}
}

type TCPTransport struct {
	opts           TransportOpts
	HandshakerFunc HandshakerFunc
	Decoder        Decoder
	rpcCh          chan RPCMessage
	OnPeer         func(Peer) error

	mu       sync.RWMutex
	ConnPool map[net.Addr]*net.Conn
}

func (transport *TCPTransport) Address() string {
	return transport.opts.listenaddr
}

func (transport *TCPTransport) LbAddress() string {
	return transport.opts.lbAddr
}

func (transport *TCPTransport) Conn(addr net.Addr) (*net.Conn, bool) {
	conn, ok := transport.ConnPool[addr]
	return conn, ok
}

func (transport *TCPTransport) Dump(addr net.Addr) {
	delete(transport.ConnPool, addr)
}

func NewTCPTransport(opts TransportOpts, onPeer func(Peer) error) *TCPTransport {
	return &TCPTransport{
		opts:           opts,
		HandshakerFunc: TCPHandshake,
		Decoder:        &DefaultDecoder{},
		rpcCh:          make(chan RPCMessage),
		OnPeer:         onPeer,
		ConnPool:       make(map[net.Addr]*net.Conn),
	}
}

func (t *TCPTransport) Consume() <-chan RPCMessage {
	return t.rpcCh
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.opts.listenaddr)
	if err != nil {
		return err
	}
	t.opts.listener = ln
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.opts.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, false)
	fmt.Printf("new incoming tcp conn %v\n", *peer)
	if err := t.HandshakerFunc(conn); err != nil {
		fmt.Println("dropping tcp connection ,err ", err)
		conn.Close()
		return
	}
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			fmt.Println("dropping tcp connection ,err ", err)
			conn.Close()
			return
		}
	}
	rpc := RPCMessage{From: conn.RemoteAddr()}

	if err := t.Decoder.Decode(conn, &rpc); err != nil {
		fmt.Println("dropping tcp connection ,err ", err)
		conn.Close()
		return
	}
	t.mu.Lock()
	if _, ok := t.ConnPool[conn.RemoteAddr()]; ok {
		t.mu.Unlock()
		fmt.Println(t.ConnPool)
		fmt.Println("dropping tcp connection , connection already exists in pool ")
		conn.Close()
		return
	}
	t.ConnPool[conn.RemoteAddr()] = &conn
	t.mu.Unlock()
	fmt.Println(rpc)
	t.rpcCh <- rpc
	//return

}
