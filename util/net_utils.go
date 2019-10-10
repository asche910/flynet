package util

import (
	"fmt"
	"net"
)

// listen tcp port at the localPort
func ListenTCP(localPort string) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", localPort))
	if err != nil {
		logger.Panicf("The port %s has been used!---> %s \n", localPort, err.Error())
	}
	logger.Printf("Client listen tcp at: %s\n", localPort)
	return listener
}

// listen udp port at the localPort
func ListenUDP(localPort string) net.Listener {
	listener, err := net.Listen("udp", fmt.Sprintf(":%s", localPort))
	if err != nil {
		logger.Panicf("The port %s has been used!---> %s \n", localPort, err.Error())
	}
	logger.Printf("Client listen udp at: %s\n", localPort)
	return listener
}
