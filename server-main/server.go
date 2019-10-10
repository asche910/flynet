package main

import (
	"fmt"
	"github.com/asche910/flynet/relay"
	"github.com/asche910/flynet/util"
	"net"
	"strconv"
)

var(
	logger = util.GetLogger()
)

func main() {
	fmt.Println("Server start: ")


	port := "8088"

	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Println("Server error!")
	}
	fmt.Println("Server listen at " + port)

	for {
		logger.Println("Waiting...")
		client, err := ln.Accept()
		if err != nil {
			fmt.Println("server listener error:", err )
			continue
		}
		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	if client == nil {
		logger.Println("Client is nil!")
	}
	logger.Println("Connect success!")

	// 请求建立socks连接
	var b [1024] byte
	n, err := client.Read(b[:])
	if err != nil {
		logger.Println("Read error!")
		return
	}

	logger.Println(b[:])

	relay.Decrease(b[:])

	logger.Println(b[:])

	// asocks5 握手连接也经过加密
	// client的socks5握手报文由本地浏览器发送
	// client的代理只起到转发作用
	// 这样握手转到了本地进行
	if b[0] == 0x05 {
		// 服务器响应连接成功
		_, _ = client.Write(relay.Increase([]byte{0x05, 0x00}))
		// 服务器读取客户端的实际访问请求(如google...)
		n, err = client.Read(b[:])

		relay.Decrease(b[:])

		var host, port string
		switch b[3] {
		case 0x01: // IPV4 ADDRESS
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: // DOMAINNAME
			host = string(b[5:n-2]) // b[4]为域名长度
		case 0x04: // IPV6 ADDRESS
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		//TODO ?
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		// 服务器向目标网站发起请求
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			logger.Println("Dial failed!")
			return
		}
		// defer server.Close()


		//响应客户端请求成功
		by := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		relay.Increase(by)
		_, _ = client.Write(by)

		go relay.DecodeTo(server, client)
		relay.EncodeTo(client, server)

	}
}
