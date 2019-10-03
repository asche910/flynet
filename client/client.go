package main

import (
	"fmt"
	"github.com/asche910/SocketsProxy/relay"
	"log"
	"net"
)

func main() {
	fmt.Println("Client start: ")

	listener, err := net.Listen("tcp", ":8848")
	if err != nil {
		log.Fatalln("The port has been used!", err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("Accept failed!")
			continue
		}
		log.Println("Client accepted!")

		go handle(client)
	}
}

// 本地中转	---> 连接远程 & 服务本地
func handle(client net.Conn) {
	server, err := net.Dial("tcp", ":8088")
	if err != nil {
		log.Println("Connect remote failed!")
		return
	}

	go relay.EncodeTo(server, client)

	relay.DecodeTo(client, server)
}
