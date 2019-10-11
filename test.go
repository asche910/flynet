package main

import (
	"fmt"
	log "github.com/asche910/flynet/log"
	"github.com/asche910/flynet/util"
	"net"
)

func main() {
	fmt.Println("Start: ")

	log.EnableDebug(true)
	log.EnableLog(true)
	logger := log.GetLogger()
	logger.Println("Hello, world!")

	net.UDPAddr{}

	_, err := net.Listen("tcp", ":80")
	util.CheckError(err, "hhhhh")

}

