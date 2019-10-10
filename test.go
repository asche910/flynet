package main

import (
	"fmt"
	"github.com/asche910/flynet/util"
	"log"
	"net"
)

func main() {
	fmt.Println("Start: ")

	util.SetEnableDebug(true)
	util.SetEnableLog(true)
	 util.GetLogger()
	//logger.Println("Hello, world!")

	_, err := net.Listen("tcp", ":80")
	util.CheckError(err, "hhhhh")

}

func TestPrint() {
	log.Println("Test Success!")
}
