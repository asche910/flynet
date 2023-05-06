package fly

import (
	"encoding/binary"
	"io"
	"net"
)

type Conn struct {
	Conn       net.Conn
	Cipher     *Cipher
	BufPipe    *BufferedPipe
	PipeReader *io.PipeReader
	PipeWriter *io.PipeWriter
}

func NewConn(con net.Conn, cipher *Cipher) *Conn {
	pr, pw := io.Pipe()
	return &Conn{
		Conn:       con,
		Cipher:     cipher,
		BufPipe:    NewBufferedPipe(4096),
		PipeReader: pr,
		PipeWriter: pw,
	}
}

// DialWithAddr pipeW dial server and send request addr of client
func DialWithAddr(server, method, key string, addr []byte) *Conn {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		logger.Println("Dial server failed --->", err)
		return nil
	}
	newConn := NewConn(conn, NewCipherInstance(key, method))

	//  2 byte of len + request data
	//body := make([]byte, len(addr)+2)
	//binary.BigEndian.PutUint16(body[:2], uint16(len(addr)))
	//copy(body[2:], addr)
	//fmt.Println("body ", body)

	if _, err := newConn.Write(addr); err != nil {
		logger.Println("Write addr to server failed --->", err)
		return nil
	}
	return newConn
}

func (conn *Conn) Write(b []byte) (n int, err error) {
	dataSize := len(b)
	var buff []byte
	buff = make([]byte, dataSize+2)
	//buff = make([]byte, 1024)
	//if len(buff) > len(b) {
	//	buff = buff[:len(b)]
	//} else {
	//	buff = make([]byte, len(b))
	//}

	binary.BigEndian.PutUint16(buff[:2], uint16(dataSize))

	//logger.Println("Before encrypt: ", b[:])
	//logger.Println("before", b)
	conn.Cipher.Encrypt(buff[2:], b)
	//logger.Println("after", buff)
	//logger.Println("After encrypt: ", len(buff))
	//n, err = conn.Conn.Write(b)
	logger.Println("Conn write")
	n, err = conn.Conn.Write(buff)
	logger.Println("Conn write len", n)

	if err != nil {
		logger.Println("Conn write failed", err)
	}
	if n != len(buff) {
		logger.Println("Conn short write ", n, len(buff))
	}
	return n, err
}

func (conn *Conn) Read(b []byte) (n int, err error) {
	return conn.Conn.Read(b)
	//buff := make([]byte, 1024)
	//logger.Println("Conn read")
	//n, err = conn.Conn.Read(buff)
	//if err != nil {
	//	logger.Println("Conn read failed", err)
	//}
	////logger.Println("before", buff[:n])
	//if n > 0 {
	//	logger.Println("Conn read len", n)
	//	b = b[:n]
	//	conn.Decrypt(b, buff[:n])
	//	//logger.Println("after", b[:n])
	//}
	//return n, err
}

func (conn *Conn) Close() error {
	err := conn.Conn.Close()
	return err
}
