package util

import (
	"io"
	"log"
	"os"
)

var (
	logFlag   = false
	debugFlag = false
	logger *log.Logger
)

func init(){
	// init logger avoid for not call GetLogger() function
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func SetEnableLog(flag bool) {
	logFlag = flag
}

func SetEnableDebug(flag bool) {
	debugFlag = flag
}

// return a logger
func GetLogger() *log.Logger {
	f, err := os.OpenFile("flynet.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Create logger file failed!", err)
	}

	var writers []io.Writer
	if logFlag && debugFlag {
		writers = []io.Writer{f, os.Stdout}
	} else if !logFlag && !debugFlag {
		writers = []io.Writer{}
	} else if debugFlag {
		writers = []io.Writer{os.Stdout}
	} else {
		writers = []io.Writer{f}
	}

	targetWriter := io.MultiWriter(writers...)
	logg := log.New(targetWriter, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger = logg
	return logg
}
