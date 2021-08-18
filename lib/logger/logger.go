package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	dirName     = "srm"
	fileName    = "log"
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger
	LogFile     *os.File
)

func CreateDirectory(dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func Init() error {

	CreateDirectory("/var/log/" + dirName)

	now := time.Now()
	year := fmt.Sprintf("%d", now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())
	date := year + month + day

	LogFile, err := os.OpenFile("/var/log/"+dirName+"/"+
		fileName+"_"+date+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	InfoLogger = log.New(io.MultiWriter(LogFile, os.Stdout), "[INFO]\t", log.Ldate|log.Ltime)
	WarnLogger = log.New(io.MultiWriter(LogFile, os.Stdout), "[WARN]\t", log.Ldate|log.Ltime)
	ErrorLogger = log.New(io.MultiWriter(LogFile, os.Stdout), "[ERROR]\t", log.Ldate|log.Ltime)
	FatalLogger = log.New(io.MultiWriter(LogFile, os.Stdout), "[FATAL]\t", log.Ldate|log.Ltime)

	return nil
}

func Info(message string) {
	InfoLogger.Println(message)
}

func Warn(message string) {
	WarnLogger.Println(message)
}

func Error(err error) {
	ErrorLogger.Println(err)
}

func Fatal(message string) {
	FatalLogger.Println(message)
}

func End() {
	_ = LogFile.Close()
}
