package client

import "github.com/asche910/flynet/fly"

type FlyClient struct {
	Mode int
	localhost string
	Ports []string // ports[0] stands for the listening port; others are for reserve
	protocol string // tcp or udp protocol
	ServerAddr string

}

func (client *FlyClient) LocalSocks5Proxy(port string) {
	fly.StartSocks5(port)
}

func (client *FlyClient) LocalHttpProxy(port string) {
	fly.StartHttp(port)
}

func (client *FlyClient) Socks5ProxyForTCP(localPort, serverAddr string) {
	fly.Socks5ForClientByTCP(localPort, serverAddr)
}

func (client *FlyClient) Socks5ProxyForUDP(localPort, serverAddr string) {
	fly.Socks5ForClientByUDP(localPort, serverAddr)
}

func (client *FlyClient) PortForward(laborPort, serverAddr string) {
	fly.PortForwardForClient(laborPort, serverAddr)
}
