package fly

import (
	"net"
)

type Conn struct {
	net.Conn
	*Cipher
}

func NewConn(con net.Conn, cipher *Cipher) *Conn {
	return &Conn{
		Conn:   con,
		Cipher: cipher,
	}
}

// dial server and send request addr of client
func DialWithAddr(server, method, key string, addr []byte) *Conn {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		logger.Println("Dial server failed --->", err)
		return nil
	}
	newConn := NewConn(conn, NewCipherInstance(key, method))
	if _, err := newConn.Write(addr); err != nil {
		logger.Println("write addr to server failed --->", err)
		return nil
	}
	return newConn
}

func (conn *Conn) Write(b []byte) (n int, err error) {
	buff := make([]byte, 1024)
	if len(buff) > len(b) {
		buff = buff[:len(b)]
	} else {
		buff = make([]byte, len(b))
	}

	conn.Encrypt(buff, b)
	n, err = conn.Conn.Write(buff)
	return
}

func (conn *Conn) Read(b []byte) (n int, err error) {
	buff := make([]byte, 1024)
	n, err = conn.Conn.Read(buff)

	//logger.Println("before", buff[:n])
	if n > 0 {
		b = b[:n]
		conn.Decrypt(b, buff[:n])
		//fmt.Println("after", b[:n])
	}
	return
}

func (conn *Conn) Close() error {
	err := conn.Conn.Close()
	return err
}
