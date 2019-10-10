package util

import (
	"fmt"
)

// just check error and print if err is not nil
func CheckError(err error, info string) {
	if err != nil {
		logger.Println(info, err)
	}
}

// check error and exit if err is not nil
func CheckErrorOrExit(err error, info string) {
	if err != nil {
		logger.Fatalln(info, err)
	}
}

// get info about port occupied
func PortOccupiedInfo(port string) string {
	return fmt.Sprintf("Port %s has been occuried!", port)
}

func AcceptErrorInfo() string {
	return "Accept client error!"
}
