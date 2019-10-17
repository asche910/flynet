package logs

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

func EnableLog(flag bool) {
	logFlag = flag
}

func EnableDebug(flag bool) {
	debugFlag = flag
}

// return a logger
func GetLogger() *log.Logger {
	if logger != nil {
		return logger
	}
	var f *os.File
	var err error
	if logFlag {
		f, err = os.OpenFile("flynet.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println("Create logger file failed!", err)
		}
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
	// cancel log.Lshortfile
	logs := log.New(targetWriter, "", log.Ldate|log.Ltime)
	logger = logs
	return logs
}
