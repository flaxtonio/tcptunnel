package main

import (
	"fmt"
	"net"
	"tcptunnel"
)

func main() {
	tunnel := tcptunnel.CreateTunnel(":8888")
	tunnel.GetAddr = func(conn *net.TCPConn) string {
		//TODO: Need to write some logic
		return "flaxton.io:80"
	}
	tunnel.ErrorHandling = func(er error) {
		fmt.Println(er.Error())
	}
	tunnel.Disconnected = func(addr string) {
		fmt.Println("Disconnected:", addr)
	}
	tunnel.Listen()
}
