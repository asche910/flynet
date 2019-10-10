package relay

import (
	log2 "github.com/asche910/flynet/log"
	"github.com/xtaci/kcp-go"
	"io"
	"log"
	"net"
)

const CIP = 5

var logger *log.Logger

func InitLog() {
	logger = log2.GetLogger()
}

func TCPToUDP(session *kcp.UDPSession, conn net.Conn) {
	buff := make([]byte, 4096)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			logger.Println(err)
			break
		}
		//log.Printf("TCPToUDP read %d byte\n", n)
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
		//log.Printf("UDPToTCP read %d byte\n", n)
		n, err = conn.Write(buff[:n])
		if err != nil {
			logger.Println(err)
			break
		}
	}
}

func EncodeTo(writer io.Writer, reader io.Reader) {
	buff := make([]byte, 1024)
	for {
		n, err := reader.Read(buff[:])
		if err != nil {
			log.Println("Endode failed!")
			return
		}
		if n < 0 {
			return
		}
		buff = Increase(buff)
		n, err = writer.Write(buff[:n])
		if err != nil {
			log.Println("Write failed!")
		}
	}
}

func DecodeTo(writer io.Writer, reader io.Reader) {
	buff := make([]byte, 1024)
	for {
		n, err := reader.Read(buff[:])
		if err != nil {
			log.Println("Decode failed!")
			return
		}
		if n < 0 {
			return
		}
		buff = Decrease(buff)
		n, err = writer.Write(buff[:n])
		if err != nil {
			log.Println("Write failed!-")
		}
	}
}

func Increase(by []byte) []byte {
	for i := 0; i < len(by); i++ {
		by[i] = by[i] + CIP
	}
	return by
}

func Decrease(by []byte) []byte {
	for i := 0; i < len(by); i++ {
		by[i] = by[i] - CIP
	}
	return by
}
