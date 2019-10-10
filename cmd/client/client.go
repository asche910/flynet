package main

import (
	"github.com/asche910/flynet/log"
	"github.com/asche910/flynet/relay"
	"github.com/asche910/flynet/util"
 	log2 "log"
)

var logger *log2.Logger

func main() {
	log.EnableLog(true)
	initLog()

	//util.Socks5ForClientByUDP("8848", "127.0.0.1:8888")
	util.PortForwardForClient("1080", ":7777")
}

func initLog()  {
	log.InitLog()
	util.InitLog()
	relay.InitLog()
	logger = log.GetLogger()
}