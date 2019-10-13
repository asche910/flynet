package relay

import (
	log2 "github.com/asche910/flynet/log"
	"github.com/xtaci/kcp-go"
	"io"
	"log"
	"net"
)

const key = "asche910-flynet-"

var (
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	logger   *log.Logger
)

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
			logger.Println("EncodeTo ---> read failed!", err)
			break
		}
		//fmt.Println(n)
		buff = Encrypt(buff, n)
		n, err = writer.Write(buff[:n])
		if err != nil {
			logger.Println("EncodeTo ---> write failed!", err)
		}
	}
}

func DecodeTo(writer io.Writer, reader io.Reader) {
	buff := make([]byte, 1024)
	for {
		n, err := reader.Read(buff[:])
		if err != nil {
			logger.Println("DecodeTo ---> read failed!", err)
			return
		}
		//fmt.Println(n)
		buff = DeCrypt(buff, n)
		n, err = writer.Write(buff[:n])
		if err != nil {
			logger.Println("DecodeTo ---> write failed!", err)
		}
	}
}

func Encrypt(by []byte, n int) []byte {
	/*	c, err := aes.NewCipher([]byte(key))
		if err != nil {
			logger.Println("aes.NewCipher failed!", err)
		}
		encrypter := cipher.NewCFBEncrypter(c, commonIV)
		var buff = make([]byte, 1024)
		encrypter.XORKeyStream(buff[:n], by[:n])*/
	// TODO add encrypt algorithm
	return by
}

func DeCrypt(by []byte, n int) []byte {
	/*		c, err := aes.NewCipher([]byte(key))
			if err != nil {
				logger.Println("aes.NewCipher failed!", err)
			}
			decrypter := cipher.NewCFBDecrypter(c, commonIV)
			var buff = make([]byte, 1024)
			decrypter.XORKeyStream(buff[:n], by[:n])*/
	// TODO add decrypt algorithm
	return by
}
