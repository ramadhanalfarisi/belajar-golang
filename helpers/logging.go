package helpers

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	filepath      = os.Getenv("LOCAL_FILEPATH")
	date          = time.Now().Format("2006-01-02")
	debug, _      = strconv.ParseBool(os.Getenv("DEBUG"))
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

func init() {
	var f *os.File
	log_path := filepath + "/logs/" + date + ".txt"
	if _, err := os.Stat(log_path); os.IsNotExist(err) {
		var errcreate error
		f, errcreate = os.Create(log_path)
		if errcreate != nil {
			log.Fatal(errcreate)
		}
	} else {
		var erropen error
		f, erropen = os.Open(log_path)
		if erropen != nil {
			log.Fatal(erropen)
		}
	}
	InfoLogger = log.New(f,"[INFO] ",log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(f,"[WARNING] ",log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f,"[ERROR] ",log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(f,"[DEBUG] ",log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string) {
	InfoLogger.Println(msg)
}

func Warning(msg string) {
	WarningLogger.Println(msg)
}

func Debug(msg string) {
	if debug {
		DebugLogger.Println(msg)
	}
}

func Error(err error) {
	var logs string
	pc, fn, line, _ := runtime.Caller(1)
	// Include function name if debugging
	if debug {
		logs = fmt.Sprintf("%s [%s:%s:%d]", err, runtime.FuncForPC(pc).Name(), fn, line)
	} else {
		logs = fmt.Sprintf("%s [%s:%d]", err, fn, line)
	}
	ErrorLogger.Println(logs)
}
