package fly

import (
	"fmt"
	"net"
	"strconv"
)

// listen tcp port at the localPort
func ListenTCP(localPort string) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", localPort))
	if err != nil {
		logger.Panicln(PortOccupiedInfo(localPort), err.Error())
	}
	logger.Printf("listen tcp on: %s\n", localPort)
	return listener
}

// listen udp port at the localPort
func ListenUDP(localPort string) net.Listener {
	listener, err := net.Listen("udp", fmt.Sprintf(":%s", localPort))
	if err != nil {
		logger.Panicln(PortOccupiedInfo(localPort), err.Error())
	}
	logger.Printf("listen udp on: %s\n", localPort)
	return listener
}

func CheckPort(port string) {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatalln("port is not a number --->", err)
	}
	if portNum < 1 || portNum > 65535 {
		logger.Fatalln("port should be in range [1,65536)")
	}
}
