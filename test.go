package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main(){
	fmt.Println("Start: ")

	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		fmt.Println("error!")
	}
	fmt.Println(conn)
	str := "hello,server!"
	conn.Write([]byte(str[:]) )
	io.Copy(os.Stdout, conn)
}

func TestPrint(){
	log.Println("Test Success!")
}