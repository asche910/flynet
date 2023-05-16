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
		logger.Debugln(flag, "start read...")
		n, err := src.Read(buff)
		if err != nil {
			if err != io.EOF {
				logger.Errorln(flag, "relayTraffic read failed:", err)
			} else {
				logger.Debugln(flag, "read EOF", n, err)
			}
			_ = dst.Close()
			break
		}
		if n > 0 {
			logger.Debugf("%s write: \n-------------- \n%s \n-------------- \n ", flag, string(buff[:n]))
			m, err := dst.Write(buff[:n]) // m = n + 4
			if err != nil {
				logger.Errorln(flag, "relayTraffic write failed", err)
				break
			}
			if m != n+DataOffset { // addition header size
				logger.Errorln(flag, "relayTraffic short write:", err)
				break
			}
			logger.Debugln(flag, "write", m)
		}

		//logger.Println(flag, "write one")
	}
	logger.Debugln(flag, "write done")
}

func RelayTrafficAndDecrypt(dst net.Conn, conn *Conn, flag string) {
	//pipeR, pipeW := io.Pipe()
	buff := make([]byte, 1024)

	// read pipeline
	go func() {
		headerBuff := make([]byte, DataOffset)
		var buffSize int
		var encryptBuff []byte
		var decryptBuff []byte

		for {
			logger.Debugln(flag, "start read header...")
			n, err := io.ReadFull(conn.BufPipe, headerBuff)
			if err != nil {
				logger.Errorln(flag, "readFull size error ", n, err)
				break
			}
			//binary.BigEndian.PutUint16(headerBuff[:2], uint16(buffSize))
			buffSize = int(binary.BigEndian.Uint16(headerBuff[2:DataOffset]))
			logger.Debugln(flag, "read header size", buffSize)
			logger.Debugln(flag, "current pipe size", conn.BufPipe.Size())

			// Check magic number
			if headerBuff[0] != 255 || headerBuff[1] != 255 {
				logger.Errorln(flag, "header error !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! ", headerBuff)
				// TODO other side stop
				break
			}

			encryptBuff = make([]byte, buffSize)
			n, err = io.ReadFull(conn.BufPipe, encryptBuff)
			if err != nil {
				logger.Errorln(flag, "readFull data error ", n, err)
				break
			}
			decryptBuff = make([]byte, buffSize)
			conn.Cipher.Decrypt(decryptBuff, encryptBuff)
			logger.Debugf("%s read decrypt body: \n-------------- \n%s \n-------------- \n", flag, string(decryptBuff))

			n, err = dst.Write(decryptBuff)
			if err != nil {
				logger.Errorln(flag, "write dst error ", n, err)
				break
			}
			logger.Debugln(flag, "write dst", n)
		}
	}()

	//io.ReadFull()
	//flag := "M->S"
	for {
		logger.Debugln(flag, "start read...")
		n, err := conn.Read(buff)
		if n > 0 {
			//logger.Printf("%s write: \n-------------- \n%s \n-------------- \n ", flag, string(buff[:n]))
		}
		if err != nil {
			if err != io.EOF {
				logger.Errorln(flag, "relayTraffic read failed:", err)
			} else {
				logger.Debugln(flag, "read EOF", n, err)
			}
			_ = dst.Close()
			break
		}
		//go func() {
		logger.Debugln(flag, " start write pipe...")
		n, err = conn.BufPipe.Write(buff[:n])
		logger.Debugln(flag, " write pipe ", n, "size ", conn.BufPipe.Size())
		if err != nil {
			logger.Errorln(flag, "RelayTraffic write pipe failed:", err)
		}
		//}()

		//logger.Println(flag, "write one")
	}
}

func RelayTrafficAndEncrypt(conn *Conn, src net.Conn, flag string) {
	RelayTrafficWithFlag(conn, src, flag)
	//io.Copy(conn, src)
}
