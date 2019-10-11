package client

import "github.com/asche910/flynet/util"

type FlyClient struct {
	mode string
	localhost string
	ports []string // ports[0] stands for the listening port; others are for reserve
	protocol string // tcp or udp protocol
	serverHost string
	serverPort string

}

func (client *FlyClient) LocalSocks5Proxy(port string) {
	util.StartSocks5(port)
}

func (client *FlyClient) LocalHttpProxy(port string) {
	util.StartHttp(port)
}

func (client *FlyClient) Socks5ProxyForTCP(localPort, serverAddr string) {
	util.Socks5ForClientByTCP(localPort, serverAddr)
}

func (client *FlyClient) Socks5ProxyForUDP(localPort, serverAddr string) {
	util.Socks5ForClientByUDP(localPort, serverAddr)
}

func (client *FlyClient) PortForward(laborPort, serverAddr string) {
	util.PortForwardForClient(laborPort, serverAddr)
}
