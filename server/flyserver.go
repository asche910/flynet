package server

import (
	"net"
)

type FlyServer struct {
	mode      string // which function to use
	localHost string
	ports     []string // ports[0] stands for the listening port; others are for reserve
	protocol  string   // tcp or udp protocol

	clients map[string]net.Conn

}

func (server *FlyServer) localSocks5() {

}
