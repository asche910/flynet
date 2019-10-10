package util

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

// start http proxy at port
func StartHttp(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	CheckErrorOrExit(err, PortOccupiedInfo(port))

	for {
		client, err := ln.Accept()
		if err != nil {
			CheckError(err, AcceptErrorInfo())
			continue
		}
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		logger.Println(err)
		return
	}
	var method, host, address string
	_, _ = fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		logger.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" { //https request
		address = hostPortURL.Scheme + ":443"
	} else {                                            //http request
		if strings.Index(hostPortURL.Host, ":") == -1 { // host not end with a port， add the default port 80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	// having already get host and port, let's dial to the target server
	server, err := net.Dial("tcp", address)
	if err != nil {
		logger.Println(err)
		return
	}
	if method == "CONNECT" {
		_, _ = fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		_, _ = server.Write(b[:n])
	}
	// forward the traffic
	go io.Copy(server, client)
	io.Copy(client, server)
}
