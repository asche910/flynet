package util

import (
	"io"
	"net"
	"sync"
)

func PortForwardForClient(localPort, serverAddr string) {
	for {
		var  localCon, serverCon net.Conn
		var err error
		for {
			localCon, err = net.Dial("tcp", ":"+localPort)
			if err != nil {
				logger.Println(err, "Dial to the target failed!")
			} else {
				logger.Printf("Dial to localhost:%s success!\n", localPort)
				break
			}
		}

		for {
			serverCon, err = net.Dial("tcp", serverAddr)
			if err != nil {
				logger.Println(err, "Dial to the target failed!")
			} else {
				logger.Printf("Dial to %s success!\n", serverAddr)
				break
			}
		}
		logger.Println("Connect success, start forwarding...")
		forward( serverCon, localCon)
	}
}

func PortForwardForServer(laborPort, queryPort string)  {
	laborLn := ListenTCP(laborPort)
	queryLn := ListenTCP(queryPort)

	for {
		laborConn, e1 := laborLn.Accept()
		queryConn, e2 := queryLn.Accept()
		if e1 != nil || e2 != nil {
			CheckError(e1, "")
			CheckError(e2, "")
			continue
		}
		logger.Println("Connect success, start forwarding...")
		forward(laborConn, queryConn)
	}
}

func forward(con1, con2 net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go copyConn(con1,con2, &wg)
	go copyConn(con2,con1, &wg)
	wg.Wait()
}

func copyConn(con1, con2 net.Conn, wg *sync.WaitGroup)  {
	io.Copy(con1, con2)
	con1.Close()
	wg.Done()
}