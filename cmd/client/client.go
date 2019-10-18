package main

import (
	"fmt"
	"github.com/asche910/flynet/client"
	"github.com/asche910/flynet/fly"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

var (
	flyClient = client.FlyClient{}
	ModeMap   = map[int]string{1: "http", 2: "socks5", 3: "socks5-tcp", 4: "socks5-udp", 5: "forward"}
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
		flyClient.Socks5ProxyForTCP(flyClient.Ports[0], flyClient.ServerAddr, flyClient.Method, flyClient.Password)
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
  -m, --method      choose a encrypt method, which must be one of ['aes-128-cfb','aes-192-cfb',
                    'aes-256-cfb', 'aes-128-ctr', 'aes-192-ctr', 'aes-256-ctr', 'rc4-md5', 
                    'rc4-md5-6', 'chacha20', 'chacha20-ietf'], default is 'aes-256-cfb'
  -P, --password    password of server
  -V, --verbose     output detail info
  -l, --log        output detail info to log file
  -H, --help        show detail usage

Mail bug reports and suggestions to <asche910@gmail.com>
or github: https://github.com/asche910/flynet`)
}

func parseArgs(args []string) {
	if len(args) < 1 {
		return
	}
	switch args[0] {
	case "--mode", "-M":
		if len(args) > 1 {
			switch args[1] {
			case ModeMap[1]:
				flyClient.Mode = 1
			case ModeMap[2], "socks":
				flyClient.Mode = 2
			case ModeMap[3], "socks-tcp":
				flyClient.Mode = 3
			case ModeMap[4], "socks-udp":
				flyClient.Mode = 4
			case ModeMap[5]:
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
			fly.CheckPort(args[1])
			if len(args) > 2 && !strings.HasPrefix(args[2], "-") {
				fly.CheckPort(args[2])
				flyClient.Ports = []string{args[1], args[2]}
				parseArgs(args[3:])
			} else {
				flyClient.Ports = []string{args[1]}
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
		} else {
			fmt.Println("flynet: no correct serverAddr found!")
			printHelp()
			os.Exit(1)
		}
	case "-m", "--method":
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			flyClient.Method = args[1]
			parseArgs(args[2:])
		} else {
			fmt.Println("flynet: no password found!")
			printHelp()
			os.Exit(1)
		}
	case "--password", "-P":
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			flyClient.Password = args[1]
			parseArgs(args[2:])
		} else {
			fmt.Println("flynet: no password found!")
			printHelp()
			os.Exit(1)
		}
	case "--verbose", "-V":
		fly.EnableDebug(true)
		parseArgs(args[1:])
	case "--logs", "-l":
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			flyClient.LogPath = args[1]
		}
		fly.EnableLog(true)
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

func checkArgs() {
	mode := flyClient.Mode
	if mode == 0 {
		fmt.Println("Please choose a mode!")
		printHelp()
		os.Exit(1)
	} else if mode == 3 || mode == 4 || mode == 5 {
		if flyClient.ServerAddr == "" {
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

// load config file
func readConfig(path string) {
	conf, err := ini.Load(path)
	if err != nil {
		fmt.Println("load config file fail --->", err)
		fmt.Println("you could refer the example config file: https://github.com/asche910/flynet")
		os.Exit(1)
	}
	section, err := conf.NewSection("client")
	if err != nil {
		fmt.Println("read client tag failed --->", err)
	}
	mode := getAttr(section, "mode")
	port := getAttr(section, "port")
	serverAddr := getAttr(section, "serverAddr")
	method := getAttr(section, "method")
	password := getAttr(section, "password")
	pacOn := getAttr(section, "pac-on")
	verbose := getAttr(section, "verbose")
	logs := getAttr(section, "log")

	switch mode {
	case ModeMap[1]:
		flyClient.Mode = 1
	case ModeMap[2], "socks":
		flyClient.Mode = 2
	case ModeMap[3], "socks-tcp":
		flyClient.Mode = 3
	case ModeMap[4], "socks-udp":
		flyClient.Mode = 4
	case ModeMap[5]:
		flyClient.Mode = 5
	default:
		fmt.Println("flynet: no correct mode found!")
		printHelp()
		os.Exit(1)
	}

	fly.CheckPort(port)
	flyClient.Ports = []string{port}

	flyClient.ServerAddr = serverAddr

	flyClient.Method = method

	flyClient.Password = password

	if pacOn == "true" {
		 flyClient.PACMode = true
	}else {
		flyClient.PACMode = false
	}

	if verbose == "true" {
		flyClient.Verbose = true
		fly.EnableDebug(true)
	}else {
		flyClient.Verbose = false
	}

	if logs != "" {
		flyClient.LogPath = logs
		fly.EnableLog(true)
	}
}

// read attr from config file
func getAttr(section *ini.Section, name string) string {
	key, e := section.GetKey(name)
	if e != nil {
		return ""
	}
	return key.String()
}

func initLog() {
	fly.InitLog()
}
