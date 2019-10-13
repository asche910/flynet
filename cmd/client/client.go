package main

import (
	"fmt"
	"github.com/asche910/flynet/client"
	"github.com/asche910/flynet/log"
	"github.com/asche910/flynet/relay"
	"github.com/asche910/flynet/util"
	log2 "log"
	"os"
	"strings"
)

var (
	logger *log2.Logger
	flyClient = client.FlyClient{}
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	parseArgs(os.Args[1:])
	checkArgs()
	initLog()

	switch flyClient.Mode {
	case 1:
		flyClient.LocalHttpProxy(flyClient.Ports[0])
	case 2:
		flyClient.LocalSocks5Proxy(flyClient.Ports[0])
	case 3:
		flyClient.Socks5ProxyForTCP(flyClient.Ports[0], flyClient.ServerAddr)
	case 4:
		flyClient.Socks5ProxyForUDP(flyClient.Ports[0], flyClient.ServerAddr)
	case 5:
		ports := flyClient.Ports
		flyClient.PortForward(ports[0], flyClient.ServerAddr)
	default:
		fmt.Println("flynet: unknown error!")
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Usage: flynet [options]
  -M, --mode        choose which mode to run. the mode must be one of['http', 'socks5', 
                    'socks5-tcp', 'socks5-udp', 'forward']
  -L, --listen      choose which port(s) to listen or forward
  -S, --server      the server address client connect to
  -V, --verbose     output detail info
  -l, --log         output detail info to log file
  -H, --help        show detail usage

Mail bug reports and suggestions to <asche910@gmail.com>
or github: https://github.com/asche910/flynet`)
}

func parseArgs(args []string) {
	if len(args) < 1{
		return
	}
	switch args[0] {
	case "--mode", "-M":
		if len(args) > 1 {
			switch args[1] {
			case "http":
				flyClient.Mode = 1
			case "socks5", "socks":
				flyClient.Mode = 2
			case "socks5-tcp", "socks-tcp":
				flyClient.Mode = 3
			case "socks5-udp", "socks-udp":
				flyClient.Mode = 4
			case "forward":
				flyClient.Mode = 5
			default:
				fmt.Println("flynet: no correct mode found!")
				printHelp()
				os.Exit(1)
			}
			parseArgs(args[2:])
		} else {
			fmt.Println("flynet: no detail mode found!")
			printHelp()
			os.Exit(1)
		}
	case "--listen", "-L":
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			port1 := util.CheckPort(args[1])
			if len(args) > 2 && !strings.HasPrefix(args[2], "-") {
				port2 := util.CheckPort(args[2])
				flyClient.Ports = []string{port1, port2}
				parseArgs(args[3:])
			} else {
				flyClient.Ports = []string{port1}
				parseArgs(args[2:])
			}
		} else {
			fmt.Println("flynet: no port found!")
			printHelp()
			os.Exit(1)
		}
	case "--server", "-S":
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			flyClient.ServerAddr = args[1]
			parseArgs(args[2:])
		}else {
			fmt.Println("flynet: no correct serverAddr found!")
			printHelp()
			os.Exit(1)
		}
	case "--verbose", "-V":
		log.EnableDebug(true)
		parseArgs(args[1:])
	case "--log", "-l":
		log.EnableLog(true)
		parseArgs(args[1:])
	case "--help", "-H":
		printHelp()
		parseArgs(args[1:])
		os.Exit(0)
	default:
		fmt.Println("flynet: please input correct command!")
		printHelp()
		os.Exit(1)
	}
}

func checkArgs()  {
	mode := flyClient.Mode
	if mode == 0 {
		fmt.Println("Please choose a mode!")
		printHelp()
		os.Exit(1)
	}else if mode == 3 || mode == 4 || mode == 5{
		if flyClient.ServerAddr == ""{
			fmt.Println("flynet: please input serverAddr!")
			printHelp()
			os.Exit(1)
		}
	}

	if len(flyClient.Ports) != 1 {
		fmt.Println("flynet: please choose a port!")
		printHelp()
		os.Exit(1)
	}
}

func initLog() {
	log.InitLog()
	util.InitLog()
	relay.InitLog()
	logger = log.GetLogger()
}
