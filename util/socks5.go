package util

import (
	"github.com/asche910/flynet/relay"
	"io"
	"log"
	"net"
	"strconv"
)

func Start(port string) error{
	listener, err := net.Listen("tcp", ":" + port)
	CheckErrorOrExit(err, "The port has been used!")

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

func handle(client net.Conn)  {
	var b [1024] byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println("Read error!")
		return
	}

	if b[0] == 0x05 {
		// response the success of handshake to client
		_, _ = client.Write(relay.Increase([]byte{0x05, 0x00}))
		// read the detail request from client
		n, err = client.Read(b[:])

		var host, port string
		switch b[3] {
		case 0x01: // IPV4 address
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: // domain
			host = string(b[5:n-2]) // b[4] stands for the length of domain
		case 0x04: // IPV6 address
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		// request to the target server
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		CheckErrorOrExit(err, "request to the target server failed!")

		// response request success to client
		by := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		_, err = client.Write(by)
		CheckError(err, "response request success to client failed!")

		go io.Copy(server, client)
		io.Copy(client, server)
	}
}