package client

type FlyClient struct {
	mode string
	localhost string
	ports []string // ports[0] stands for the listening port; others are for reserve
	protocol string // tcp or udp protocol
	serverHost string
	serverPort string



}