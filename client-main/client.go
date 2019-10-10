package main

import (
	"fmt"
	"github.com/asche910/flynet/relay"
	"github.com/asche910/flynet/util"
	"net"
)

var(
	logger = util.GetLogger()
)

func main() {
	fmt.Println("Client start: ")


	listener, err := net.Listen("tcp", ":8848")
	if err != nil {
		logger.Fatalln("The port has been used!", err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			logger.Println("Accept failed!")
			continue
		}
		logger.Println("Client accepted!")

		go handle(client)
	}
}

// 本地中转	---> 连接远程 & 服务本地
func handle(client net.Conn) {
	server, err := net.Dial("tcp", ":8088")
	if err != nil {
		logger.Println("Connect remote failed!")
		return
	}

	go relay.EncodeTo(server, client)

	relay.DecodeTo(client, server)
}
