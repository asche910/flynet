package server

import (
	"github.com/asche910/flynet/fly"
	"net"
)

type FlyServer struct {
	Mode      int // which function to use, [1:'http', 2:'socks5', 3:'socks5-tcp', 4:'socks5-udp', 5:'forward']
	localHost string
	Ports     []string // ports[0] stands for the listening port; others are for reserve
	protocol  string   // tcp or udp protocol

	clients map[string]net.Conn
}

func (server *FlyServer) LocalSocks5Proxy(port string) {
	fly.StartSocks5(port)
}

func (server *FlyServer) LocalHttpProxy(port string) {
	fly.StartHttp(port)
}

func (server *FlyServer) Socks5ProxyForTCP(localPort string) {
	fly.Socks5ForServerByTCP(localPort)
}

func (server *FlyServer) Socks5ProxyForUDP(localPort string) {
	fly.Socks5ForServerByUDP(localPort)
}

func (server *FlyServer) PortForward(laborPort, queryPort string) {
	fly.PortForwardForServer(laborPort, queryPort)
}
