package fly

import (
	"encoding/binary"
	"github.com/xtaci/kcp-go"
	"io"
	"net"
)

func TCPToUDP(session *kcp.UDPSession, conn net.Conn) {
	buff := make([]byte, 4096)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			logger.Println(err)
			break
		}
		//logs.Printf("TCPToUDP read %d byte\n", n)
		n, err = session.Write(buff[:n])
		if err != nil {
			logger.Println(err)
			break
		}
	}
}

func UDPToTCP(conn net.Conn, session *kcp.UDPSession) {
	buff := make([]byte, 4096)
	for {
		n, err := session.Read(buff)
		if err != nil {
			logger.Println(err)
			break
		}
		//logs.Printf("UDPToTCP read %d byte\n", n)
		n, err = conn.Write(buff[:n])
		if err != nil {
			logger.Println(err)
			break
		}
	}
}

func RelayTraffic(dst, src net.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := src.Read(buff)
		if n > 0 {
			m, err := dst.Write(buff[:n])
			if err != nil {
				logger.Println("RelayTraffic write failed:", err)
				break
			}
			if m != n {
				logger.Println("RelayTraffic short write:", err)
				break
			}
		}
		if err != nil {
			if err != io.EOF {
				logger.Println("RelayTraffic read failed:", err)
			}
			break
		}
	}
}

func RelayTrafficWithFlag(dst *Conn, src net.Conn, flag string) {
	buff := make([]byte, 1024)
	//conn, needDecrypt := src.(*Conn)
	//conn, needEncrypt := dst.(*Conn)
	//
	//var buffR *io.PipeReader
	//var buffW *io.PipeWriter
	//if needDecrypt {
	//	buffR, buffW = io.Pipe()
	//}

	for {
		logger.Println(flag, "start read...")
		n, err := src.Read(buff)
		if n > 0 {
			logger.Printf("%s write: \n-------------- \n%s \n-------------- \n ", flag, string(buff[:n]))
			m, err := dst.Write(buff[:n]) // m = n + 2
			if err != nil {
				logger.Println(flag, "RelayTraffic write failed:", err)
				break
			}
			if m != n+2 { // addition header size
				logger.Println(flag, "RelayTraffic short write:", err)
				break
			}
		}
		if err != nil {
			if err != io.EOF {
				logger.Println(flag, "RelayTraffic read failed:", err)
			} else {
				logger.Println(flag, "read EOF", n, err)
			}
			_ = dst.Close()
			break
		}
		//logger.Println(flag, "write one")
	}
	logger.Println(flag, "write done")
}

func RelayTrafficAndDecrypt(dst net.Conn, conn *Conn, flag string) {
	//pipeR, pipeW := io.Pipe()
	buff := make([]byte, 1024)

	// read pipeline
	go func() {
		sizeBuff := make([]byte, 2)
		var buffSize int
		var encryptBuff []byte
		var decryptBuff []byte

		for {
			n, err := io.ReadFull(conn.BufPipe, sizeBuff)
			if err != nil {
				logger.Println("ReadFull size error ", n, err)
				break
			}
			//binary.BigEndian.PutUint16(sizeBuff[:2], uint16(buffSize))
			buffSize = int(binary.BigEndian.Uint16(sizeBuff))
			logger.Println("Read header size", buffSize)
			logger.Println("Current pipe size", conn.BufPipe.Size())

			encryptBuff = make([]byte, buffSize)
			n, err = io.ReadFull(conn.BufPipe, encryptBuff)
			if err != nil {
				logger.Println("ReadFull data error ", n, err)
				break
			}
			decryptBuff = make([]byte, buffSize)
			conn.Cipher.Decrypt(decryptBuff, encryptBuff)
			logger.Println("Read decrypt body ", string(decryptBuff))

			n, err = dst.Write(decryptBuff)
			if err != nil {
				logger.Println("Write dst error ", n, err)
				break
			}
		}
	}()

	//io.ReadFull()
	//flag := "M->S"
	for {
		logger.Println(flag, "start read...")
		n, err := conn.Read(buff)
		if n > 0 {
			//logger.Printf("%s write: \n-------------- \n%s \n-------------- \n ", flag, string(buff[:n]))
		}
		if err != nil {
			if err != io.EOF {
				logger.Println(flag, "RelayTraffic read failed:", err)
			} else {
				logger.Println(flag, "read EOF", n, err)
			}
			_ = dst.Close()
			break
		}
		//go func() {
		n, err = conn.BufPipe.Write(buff[:n])
		if err != nil {
			logger.Println(flag, "RelayTraffic write pipe failed:", err)
		}
		//}()

		//logger.Println(flag, "write one")
	}
}

func RelayTrafficAndEncrypt(conn *Conn, src net.Conn, flag string) {
	RelayTrafficWithFlag(conn, src, flag)
	//io.Copy(conn, src)
}
