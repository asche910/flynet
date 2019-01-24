package relay

import (
	"io"
	"log"
)

const CIP = 5

func EncodeTo( writer io.Writer, reader io.Reader) {
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
