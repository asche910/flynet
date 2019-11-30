package main

import (
	"fmt"
	"github.com/asche910/flynet/fly"
	"github.com/asche910/flynet/server"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

var (
	flyServer = server.FlyServer{}
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

	switch flyServer.Mode {
	case 1:
		flyServer.LocalHttpProxy(flyServer.Ports[0])
	case 2:
		flyServer.LocalSocks5Proxy(flyServer.Ports[0])
	case 3:
		flyServer.Socks5ProxyForTCP(flyServer.Ports[0], flyServer.Method, flyServer.Password)
	case 4:
		flyServer.Socks5ProxyForUDP(flyServer.Ports[0])
	case 5:
		ports := flyServer.Ports
		flyServer.PortForward(ports[0], ports[1])
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
  -m, --method      choose a encrypt method, which must be one of ['aes-128-cfb','aes-192-cfb',
                    'aes-256-cfb', 'aes-128-ctr', 'aes-192-ctr', 'aes-256-ctr', 'rc4-md5', 
                    'rc4-md5-6', 'chacha20', 'chacha20-ietf'], default is 'aes-256-cfb'
  -P, --password    password for client to connect
  -C, --config      read from config file
  -V, --verbose     output detail info
  -l, --log         output detail info to log file, which if not indicated "flynet.log" will 
                    be default
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
				flyServer.Mode = 1
			case ModeMap[2], "socks":
				flyServer.Mode = 2
			case ModeMap[3], "socks-tcp":
				flyServer.Mode = 3
			case ModeMap[4], "socks-udp":
				flyServer.Mode = 4
			case ModeMap[5]:
				flyServer.Mode = 5
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
		if hasNAttributes(args, 1) {
			fly.CheckPort(args[1])
			if hasNAttributes(args, 2) {
				fly.CheckPort(args[2])
				flyServer.Ports = []string{args[1], args[2]}
				parseArgs(args[3:])
			} else {
				flyServer.Ports = []string{args[1]}
				parseArgs(args[2:])
			}
		} else {
			fmt.Println("flynet: no port found!")
			printHelp()
			os.Exit(1)
		}
	case "-m", "--method":
		if hasNAttributes(args, 1) {
			flyServer.Method = args[1]
			parseArgs(args[2:])
		} else {
			fmt.Println("fly: no password found!")
			printHelp()
			os.Exit(1)
		}
	case "--password", "-P":
		if hasNAttributes(args, 1) {
			flyServer.Password = args[1]
			parseArgs(args[2:])
		} else {
			fmt.Println("flynet: no password found!")
			printHelp()
			os.Exit(1)
		}
	case "--config", "-C":
		if hasNAttributes(args, 1) {
			readConfig(args[1])
		} else {
			fmt.Println("flynet: no config file found!")
			printHelp()
			os.Exit(1)
		}
	case "--verbose", "-V":
		fly.EnableDebug(true)
		parseArgs(args[1:])
	case "--logs", "-l":
		if hasNAttributes(args, 1) {
			fly.SetLogName(args[1])
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

// judge if the current args option has  n attributes at least.
// for example:
// args = []string{"--listen", "1080", "8888"}
// hasNAttributes(args, 2) will return true
func hasNAttributes(args []string, n int) bool {
	if len(args) > n {
		for i := 1; i <= n; i++ {
			if strings.HasPrefix(args[i], "-") {
				return false
			}
		}
		return true
	}
	return false
}

func checkArgs() {
	if flyServer.Mode == 0 {
		fmt.Println("Please choose a mode!")
		printHelp()
		os.Exit(1)
	} else if flyServer.Mode == 5 {
		if len(flyServer.Ports) != 2 {
			fmt.Println("flynet: please input two port!")
			printHelp()
			os.Exit(1)
		}
	} else {
		if len(flyServer.Ports) != 1 {
			fmt.Println("flynet: please choose a port to listen!")
			printHelp()
			os.Exit(1)
		}
	}
}

// load config file
func readConfig(path string) {
	conf, err := ini.Load(path)
	if err != nil {
		fmt.Println("load config file fail --->", err)
		fmt.Println("you could refer to the example config file: https://github.com/asche910/flynet")
		os.Exit(1)
	}
	section, err := conf.NewSection("server")
	if err != nil {
		fmt.Println("read client tag failed --->", err)
	}
	mode := getAttr(section, "mode")
	portStr := getAttr(section, "port")
	method := getAttr(section, "method")
	password := getAttr(section, "password")
	verbose := getAttr(section, "verbose")
	logs := getAttr(section, "log")

	switch mode {
	case ModeMap[1]:
		flyServer.Mode = 1
	case ModeMap[2], "socks":
		flyServer.Mode = 2
	case ModeMap[3], "socks-tcp":
		flyServer.Mode = 3
	case ModeMap[4], "socks-udp":
		flyServer.Mode = 4
	case ModeMap[5]:
		flyServer.Mode = 5
	default:
		fmt.Println("flynet: no correct mode found!")
		printHelp()
		os.Exit(1)
	}

	// in case of forward mode, which needs two port
	ports := strings.Fields(portStr)
	if len(ports) == 2 {
		fly.CheckPort(ports[0])
		fly.CheckPort(ports[1])
		flyServer.Ports = []string{ports[0], ports[1]}
	}else if len(ports) == 1 {
		fly.CheckPort(ports[0])
		flyServer.Ports = []string{ports[0]}
	}
	flyServer.Method = method
	flyServer.Password = password

	if verbose == "true" {
		fly.EnableDebug(true)
	}

	if logs != "" {
		fly.SetLogName(logs)
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
