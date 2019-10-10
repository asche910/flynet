package util

import (
	"crypto/sha1"
	"fmt"
	"github.com/asche910/flynet/relay"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"net"
	"strconv"
)

func StartSocks5(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	CheckErrorOrExit(err, PortOccupiedInfo(port))

	for {
		client, err := listener.Accept()
		if err != nil {
			logger.Println("Accept failed!")
			continue
		}
		logger.Println("Client accepted!")

		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	var b [1024] byte
	n, err := client.Read(b[:])
	if err != nil {
		logger.Println("Read error!")
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
			host = string(b[5 : n-2]) // b[4] stands for the length of domain
		case 0x04: // IPV6 address
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		// request to the target server
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			CheckError(err, "request to the target server failed!")
			return
		}

		// response request success to client
		by := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		_, err = client.Write(by)
		CheckError(err, "response 'request success' to client failed!")

		go io.Copy(server, client)
		io.Copy(client, server)
	}
}

func Socks5ForClientByTCP(localPort, serverAddr string) {
	listener := ListenTCP(localPort)
	for {
		client, err := listener.Accept()
		if err != nil {
			logger.Println("Accept failed!")
			continue
		}
		logger.Println("Client accepted!")

		go func() {
			server, err := net.Dial("tcp", serverAddr)
			if err != nil {
				logger.Println("Connect remote failed!")
				return
			}
			go relay.EncodeTo(server, client)
			relay.DecodeTo(client, server)
		}()
	}
}

func Socks5ForServerByTCP(localPort string) {
	listener := ListenTCP(localPort)
	for {
		logger.Println("Waiting...")
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("server listener error:", err)
			continue
		}
		go func() {
			data := make([]byte, 1024)
			n, err := client.Read(data[:])
			if err != nil {
				logger.Println("Read error!")
				return
			}

			//logger.Println(data[:])
			relay.Decrease(data[:])
			//logger.Println(data[:])

			if data[0] == 0x05 {
				// response the success of handshake to client
				_, _ = client.Write(relay.Increase([]byte{0x05, 0x00}))
				// read the detail request from client
				n, err = client.Read(data[:])
				relay.Decrease(data[:])

				var host, port = parseSocksRequest(data, n)
				// request to the target server
				server, err := net.Dial("tcp", net.JoinHostPort(host, port))
				if err != nil {
					logger.Println("Dial failed!")
					return
				}
				// response request success to client
				by := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				relay.Increase(by)
				_, _ = client.Write(by)

				go relay.DecodeTo(server, client)
				relay.EncodeTo(client, server)
			}
		}()
	}
}

func Socks5ForClientByUDP(localPort, serverAddr string) {
	listener := ListenTCP(localPort)
	for {
		con, err := listener.Accept()
		if err != nil {
			logger.Println("Accept error: ", err)
			continue
		}
		logger.Println("Client accepted!")

		go func() {
			var b [1024] byte
			_, err := con.Read(b[:])
			if err != nil {
				logger.Println("Read error!")
				return
			}
			if b[0] == 0x05 {
				_, _ = con.Write([]byte{0x05, 0x00})

				key := pbkdf2.Key([]byte("flynet"), []byte("asche910"), 1024, 32, sha1.New)
				block, _ := kcp.NewAESBlockCrypt(key)
				session, err := kcp.DialWithOptions(serverAddr, block, 10, 3)
				if err != nil {
					logger.Println("connect targetServer failed! ", err)
					return
				}
				go relay.TCPToUDP(session, con)
				relay.UDPToTCP(con, session)
			}
		}()
	}
}

func Socks5ForServerByUDP(localPort string) {
	key := pbkdf2.Key([]byte("flynet"), []byte("asche910"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	if listener, err := kcp.ListenWithOptions(":"+localPort, block, 10, 3); err == nil {
		logger.Printf("Server listent udp at %s\n", localPort)
		for {
			con, err := listener.AcceptKCP()
			if err != nil {
				logger.Fatal(err)
			}
			go func() {
				data := make([]byte, 1024)
				n, err := con.Read(data)
				if err != nil {
					logger.Println(err)
					return
				}
				//logger.Println(string(data[:n]))
				if data[0] == 5 && data[1] == 1 && data[2] == 0 {
					// 解析并请求内容
					var host, port = parseSocksRequest(data, n)

					// 服务器向目标网站发起请求
					server, err := net.Dial("tcp", net.JoinHostPort(host, port))
					logger.Printf("Request remote %s:%s\n", host, port)

					if err != nil {
						logger.Println("Dial failed!")
						CheckError(err, "Dial remote failed!")
						return
					}

					//响应客户端请求成功
					by := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
					n, err = con.Write(by)
					if err != nil {
						logger.Println("Response to client failed! ", err)
						return
					}
					go relay.UDPToTCP(server, con)
					relay.TCPToUDP(con, server)
				} else {
					logger.Println("Unrecognized protocol!")
					return
				}
			}()
		}
	} else {
		logger.Panicln(err)
	}
}

// parse socks5 request for target host and port
func parseSocksRequest(data []byte, n int) (string, string) {
	// 解析并请求内容
	var host, port string
	switch data[3] {
	case 0x01: // IPV4 address
		host = net.IPv4(data[4], data[5], data[6], data[7]).String()
	case 0x03: // domain
		host = string(data[5 : n-2]) // data[4] stands for the length of domain
	case 0x04: // IPV6 address
		host = net.IP{data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15], data[16], data[17], data[18], data[19]}.String()
	}
	port = strconv.Itoa(int(data[n-2])<<8 | int(data[n-1]))
	return host, port
}
