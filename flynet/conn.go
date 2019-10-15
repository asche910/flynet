package flynet

import "net"

type Conn struct {
	net.Conn
	*Cipher
}

func (conn *Conn) Write(b []byte) (n int, err error) {
	buff := make([]byte, 1024)
	if len(buff) > len(b){
		buff = buff[:len(b)]
	}else {
		buff = make([]byte, len(b))
	}

	conn.Encrypt(buff, b)
	n, err = conn.Conn.Write(buff)
	return
}

func (conn *Conn) Read(b []byte) (n int, err error) {
	buff := make([]byte, 1024)
	if len(buff) > len(b){
		buff = buff[:len(b)]
	}else {
		buff = make([]byte, len(b))
	}

	n, err = conn.Conn.Read(buff)
	if n > 0 {
		conn.Decrypt(b[:n], buff[:n])
	}
	return
}

func (conn *Conn) Close() error {
	err := conn.Conn.Close()
	return err
}
