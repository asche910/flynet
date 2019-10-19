package fly

import (
	"io"
	"log"
	"os"
)

var (
	logFlag   = false
	debugFlag = false
	logger    *log.Logger
	logName = "flynet.log"
)

func InitLog() {
	logger = GetLogger()
}

func EnableLog(flag bool) {
	logFlag = flag
}

func EnableDebug(flag bool) {
	debugFlag = flag
}

// set log file name, which can include absolute path
func SetLogName(name string)  {
	logName = name
}

// return a logger
func GetLogger() *log.Logger {
	if logger != nil {
		return logger
	}
	var f *os.File
	var err error
	if logFlag {
		f, err = os.OpenFile(logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println("create or open log file failed --->", err)
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
	return logs
}
