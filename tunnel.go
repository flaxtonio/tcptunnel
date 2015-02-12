package tcptunnel

import (
	"net"
)

type Tunnel struct {
	SourceConn    *net.TCPConn
	RemoteConn    *net.TCPConn
	EventLoop     bool // if True then tunneling will be working using async Event Loop based on channels
	StopReadChan  chan bool
	StopWriteChan chan bool
	Terminate     chan bool
	BlackList     []string // List of blacklisted IP addresses
	//TODO: Add more configurations
}

func CreateTunnel(source, remote *net.TCPConn) Tunnel {
	return Tunnel{
		SourceConn: source,
		RemoteConn: remote,
	}
}

func rw_socket(t *Tunnel, SourceConn, RemoteConn *net.TCPConn) { // Read Write to Socket
	buf_receive := make([]byte, 1024)
	for {
		if t.Terminate {
			return
		}
		t.StopReadChan <-
		rlen, err := SourceConn.Read(buf_receive)
		if err != nil {
			//TODO: Write error handlers
			return
		}
		RemoteConn.Write(buf_receive[:rlen])
		//TODO: Write functionality of read bytes transfer
	}
}

func (t *Tunnel) Start() {
	go rw_socket(t, t.SourceConn, t.RemoteConn)
	go rw_socket(t, t.RemoteConn, t.SourceConn)
}
