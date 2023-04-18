// log is used for logging
package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (
	mu       sync.RWMutex
	logLevel int
)

func init() {
	log.SetFlags(0)
}

// SetLogLevel sets which log level to do, -1 = no logging, 0 = debug, 1 = info, 2 = warn, 3 = error
func SetLogLevel(level int) {
	logLevel = level
}

// LogToFile enables logging to files
func LogToFile() {
	os.Remove("quail.log")
	f, err := os.OpenFile("quail.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed opening log file:", err.Error())
		os.Exit(1)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
}

func Println(v ...interface{}) {
	Infoln(v...)
}

func Printf(format string, v ...interface{}) {
	Infof(format, v...)
}

func logPrintf(format string, v ...interface{}) {
	if logLevel > 0 {
		log.Printf(format, v...)
		return
	}
	mu.Lock()
	var ok bool
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	mu.Unlock()
	file = strings.ReplaceAll(file, "/Users/xackery/Documents/code/projects/quail/", "")
	log.Printf(fmt.Sprintf("./%s:%d %s", file, line, format), v...)
}

func logPrintln(v ...interface{}) {
	if logLevel > 0 {
		log.Println(v...)
		return
	}
	mu.Lock()
	var ok bool
	_, file, line, ok := runtime.Caller(0)
	if !ok {
		file = "???"
		line = 0
	}
	mu.Unlock()
	file = strings.ReplaceAll(file, "/Users/xackery/Documents/code/projects/quail/", "")
	v = append([]interface{}{fmt.Sprintf("%s:%d ", file, line)}, v...)

	log.Println(v...)
}

func Debugln(v ...interface{}) {
	if logLevel > 0 {
		return
	}
	logPrintln(v...)
}

func Infoln(v ...interface{}) {
	if logLevel > 1 {
		return
	}
	logPrintln(v...)
}

func Warnln(v ...interface{}) {
	if logLevel > 2 {
		return
	}
	logPrintln(v...)
}

func Errorln(v ...interface{}) {
	if logLevel > 3 {
		return
	}
	logPrintln(v...)
}

func Debugf(format string, v ...interface{}) {
	if logLevel > 0 {
		return
	}
	logPrintf(format, v...)
}

func Infof(format string, v ...interface{}) {
	if logLevel > 1 {
		return
	}
	logPrintf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	if logLevel > 2 {
		return
	}
	logPrintf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	if logLevel > 3 {
		return
	}
	logPrintf(format, v...)
}
