package fly

import (
	"crypto/sha1"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"net"
	"strconv"
	"strings"
)

func StartSocks5(port string) {
	listener := ListenTCP(port)
	for {
		client, err := listener.Accept()
		if err != nil {
			logger.Println("accept failed!")
			continue
		}
		logger.Println("client accepted!")

		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	data := make([]byte, 1024)
	n, err := client.Read(data[:])
	if err != nil {
		logger.Println("read error!")
		return
	}

	if data[0] == 0x05 {
		// response the success of handshake to client
		_, err = client.Write([]byte{0x05, 0x00})
		if err != nil {
			logger.Println("response to client failed!")
			return
		}
		// read the detail request from client
		n, err = client.Read(data[:])
		if err != nil {
			logger.Println("read from client failed!")
			return
		}

		var host, port = parseSocksRequest(data, n)
		logger.Printf("start request %s:%s\n", host, port)

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

func Socks5ForClientByTCP(localPort, serverAddr, method, key string) {
	listener := ListenTCP(localPort)
	for {
		client, err := listener.Accept()
		if err != nil {
			logger.Println("accept failed --->", err)
			continue
		}
		logger.Println("client accepted!")

		go func() {
			buff := make([]byte, 1024)
			n, err := client.Read(buff)
			if err != nil {
				logger.Println("read handshake request failed --->", err)
				return
			}
			if buff[0] == 0x05 {
				if n, err = client.Write([]byte{0x05, 0x00}); err != nil {
					logger.Println("write handshake response failed --->", err)
					return
				}

				// read detail request
				//if n, err = client.Read(buff); err != nil {
				if n, err = io.ReadAtLeast(client, buff, 5); err != nil {
					logger.Println("read client quest failed --->", err)
					return
				}
				replyBy := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				if _, err = client.Write(replyBy); err != nil {
					logger.Println("write 'request success' failed --->", err)
					return
				}

				//
				//        +----+-----+-------+------+----------+----------+
				//        |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
				//        +----+-----+-------+------+----------+----------+
				//        | 1  |  1  | X'00' |  1   | Variable |    2     |
				//        +----+-----+-------+------+----------+----------+
				//
				server := DialWithAddr(serverAddr, method, key, buff[:n])
				if server == nil {
					return
				}

				go RelayTraffic(server, client)
				RelayTraffic(client, server)
			}
		}()
	}
}

func Socks5ForServerByTCP(localPort, method, key string) {
	listener := ListenTCP(localPort)
	for {
		logger.Println("waiting...")
		client, err := listener.Accept()
		if err != nil {
			logger.Println("server accept failed --->", err)
			continue
		}
		go func() {
			buff := make([]byte, 1024)
			conn := NewConn(client, NewCipherInstance(key, method))
			n, err := conn.Read(buff)
			if err != nil {
				logger.Println("read target address failed --->", err)
				return
			}

			host, port := parseSocksRequest(buff[:n], n)
			//logger.Printf("target server ------\n%s:%s\n------\n%d\n+++++++\n", host, port, buff[:n])
			logger.Printf("target server ------\n%s:%s\n------\n", host, port)

			// dial the target server
			server, err := net.Dial("tcp", net.JoinHostPort(host, port))
			if err != nil {
				logger.Println("dial target server failed --->", err)
				return
			}
			go RelayTraffic(server, conn)
			RelayTraffic(conn, server)
		}()
	}
}

func Socks5ForClientByUDP(localPort, serverAddr string) {
	listener := ListenTCP(localPort)
	for {
		con, err := listener.Accept()
		if err != nil {
			logger.Println("accept error: ", err)
			continue
		}
		logger.Println("client accepted!")

		go func() {
			var b [1024] byte
			_, err := con.Read(b[:])
			if err != nil {
				logger.Println("read error!")
				return
			}
			if b[0] == 0x05 {
				_, _ = con.Write([]byte{0x05, 0x00})

				// kcp can't resolve addr such as ":8080"
				if strings.HasPrefix(serverAddr, ":") {
					serverAddr = "127.0.0.1" + serverAddr
				}

				key := pbkdf2.Key([]byte("fly"), []byte("asche910"), 1024, 32, sha1.New)
				block, _ := kcp.NewAESBlockCrypt(key)
				session, err := kcp.DialWithOptions(serverAddr, block, 10, 3)
				if err != nil {
					logger.Println("connect targetServer failed! ", err)
					return
				}
				go TCPToUDP(session, con)
				UDPToTCP(con, session)
			}
		}()
	}
}

func Socks5ForServerByUDP(localPort string) {
	key := pbkdf2.Key([]byte("fly"), []byte("asche910"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	if listener, err := kcp.ListenWithOptions(":"+localPort, block, 10, 3); err == nil {
		logger.Printf("server listent udp at %s\n", localPort)
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
					var host, port = parseSocksRequest(data, n)

					// request to the target server
					server, err := net.Dial("tcp", net.JoinHostPort(host, port))
					logger.Printf("request remote %s:%s\n", host, port)

					if err != nil {
						logger.Println("dial failed!")
						CheckError(err, "dial remote failed!")
						return
					}

					// response request success to client
					by := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
					n, err = con.Write(by)
					if err != nil {
						logger.Println("response to client failed! ", err)
						return
					}
					go UDPToTCP(server, con)
					TCPToUDP(con, server)
				} else {
					logger.Println("unrecognized protocol!")
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
	// TODO sometimes there are some nums, such as '22 3 1 2 0 1 0 1 252...' after port. why?
	var host, port string
	var p1, p2 byte
	switch data[3] {
	case 0x01: // IPV4 address
		host = net.IPv4(data[4], data[5], data[6], data[7]).String()
		p1 = data[8]
		p2 = data[9]
	case 0x03: // domain
		host = string(data[5 : 5+data[4]]) // data[4] stands for the length of domain
		p1 = data[data[4]+5]
		p2 = data[data[4]+6]
	case 0x04: // IPV6 address
		host = net.IP{data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15], data[16], data[17], data[18], data[19]}.String()
		p1 = data[20]
		p2 = data[21]
	}
	port = strconv.Itoa(int(p1)<<8 | int(p2))
	return host, port
}
