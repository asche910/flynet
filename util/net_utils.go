package util

import (
	"fmt"
	"net"
	"strconv"
)

// listen tcp port at the localPort
func ListenTCP(localPort string) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", localPort))
	if err != nil {
		logger.Panicf("the port %s has been used!---> %s \n", localPort, err.Error())
	}
	logger.Printf("listen tcp at: %s\n", localPort)
	return listener
}

// listen udp port at the localPort
func ListenUDP(localPort string) net.Listener {
	listener, err := net.Listen("udp", fmt.Sprintf(":%s", localPort))
	if err != nil {
		logger.Panicf("the port %s has been used!---> %s \n", localPort, err.Error())
	}
	logger.Printf("listen udp at: %s\n", localPort)
	return listener
}

func CheckPort(port string) string {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatalln("port is not a number!")
	}
	if portNum < 1 || portNum > 65535 {
		logger.Fatalln("port should be in range [1,65536)")
	}
	return port
}
