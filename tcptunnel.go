package tcptunnel

import (
	"io"
	"net"
)

type RemoteAddrTrigger func(*net.TCPConn) string
type ErrorHandler func(error)
type DisconnectCallback func(string)

type Tunnel struct {
	ListenAddress string
	GetAddr       RemoteAddrTrigger // Run logic for load balancing and get server address
	ErrorHandling ErrorHandler
	Disconnected  DisconnectCallback
}

func CreateTunnel(address string) Tunnel {
	return Tunnel{
		ListenAddress: address,
		ErrorHandling: func(err error) {},
		Disconnected:  func(adr string) {},
	}
}

func (t *Tunnel) Listen() {
	addr, addr_err := net.ResolveTCPAddr("tcp", t.ListenAddress)
	if addr_err != nil {
		t.ErrorHandling(addr_err)
		return
	}

	socket, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.ErrorHandling(err)
		return
	}
	for {
		conn, packet_err := socket.AcceptTCP()
		if packet_err != nil {
			t.ErrorHandling(packet_err)
			continue
		}
		go t.Start(conn)
	}
}

func (t *Tunnel) Start(conn *net.TCPConn) {
	var remote_addr string
	remote_addr = t.GetAddr(conn)
	addr, addr_err := net.ResolveTCPAddr("tcp", remote_addr)
	if addr_err != nil {
		t.ErrorHandling(addr_err)
		return
	}
	rem, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		t.ErrorHandling(addr_err)
		return
	}

	go io.Copy(conn, rem)
	io.Copy(rem, conn)
	t.Disconnected(remote_addr) // Call Disconnected callback with address
}
