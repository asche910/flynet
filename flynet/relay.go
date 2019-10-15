package flynet

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
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

		fmt.Println("before:", buff)
		buff = DeCrypt(buff, n)
		fmt.Println("after:", buff)

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
		}*/
	//encrypter := cipher.NewCFBEncrypter(c, commonIV)
	var buff = make([]byte, 1024)
	//encrypter.XORKeyStream(buff[:n], by[:n])
	// TODO add encrypt algorithm

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	stream := cipher.NewCTR(block, commonIV)
	stream.XORKeyStream(buff[:n], by[:n])
	return buff
}

func DeCrypt(by []byte, n int) []byte {
	/*c, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.Println("aes.NewCipher failed!", err)
	}
	decrypter := cipher.NewCFBDecrypter(c, commonIV)*/
	//var buff = make([]byte, 1024)
	//decrypter.XORKeyStream(buff[:n], by[:n])
	//TODO add decrypt algorithm

	return Encrypt(by, n)
}

func RelayTraffic(dst, src net.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := src.Read(buff)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Read", n)
		if n > 0 {
			if n, err = dst.Write(buff[:n]); err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}
