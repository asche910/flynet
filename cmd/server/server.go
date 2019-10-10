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

	util.Socks5ForServerByUDP("8888")
}

func initLog()  {
	log.InitLog()
	util.InitLog()
	relay.InitLog()
	logger = log.GetLogger()
}