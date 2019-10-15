package client

import "github.com/asche910/flynet/flynet"

type FlyClient struct {
	Mode int
	localhost string
	Ports []string // ports[0] stands for the listening port; others are for reserve
	protocol string // tcp or udp protocol
	ServerAddr string

}

func (client *FlyClient) LocalSocks5Proxy(port string) {
	flynet.StartSocks5(port)
}

func (client *FlyClient) LocalHttpProxy(port string) {
	flynet.StartHttp(port)
}

func (client *FlyClient) Socks5ProxyForTCP(localPort, serverAddr string) {
	flynet.Socks5ForClientByTCP(localPort, serverAddr)
}

func (client *FlyClient) Socks5ProxyForUDP(localPort, serverAddr string) {
	flynet.Socks5ForClientByUDP(localPort, serverAddr)
}

func (client *FlyClient) PortForward(laborPort, serverAddr string) {
	flynet.PortForwardForClient(laborPort, serverAddr)
}
