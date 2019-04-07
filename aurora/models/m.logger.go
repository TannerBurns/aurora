package models

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	handleName string
	Debug      *log.Logger
	Info       *log.Logger
	Warning    *log.Logger
	Error      *log.Logger
	Fatal      *log.Logger
	Log        *log.Logger
}

func NewLogger() *Logger {
	logger := &Logger{}

	logger.InitLogging("Cage", os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr, os.Stdout)
	return logger
}

func (logger *Logger) InitLogging(
	handleName string,
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	fatalHandle io.Writer,
	basicHandle io.Writer,
) {
	logger.handleName = handleName
	logger.Debug = log.New(debugHandle, fmt.Sprintf("%s - DEBUG: ", logger.handleName),
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Info = log.New(infoHandle, fmt.Sprintf("%s - ", logger.handleName),
		0)
	logger.Warning = log.New(warningHandle, fmt.Sprintf("%s - WARNING: ", logger.handleName),
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Error = log.New(errorHandle, fmt.Sprintf("%s - ERROR: ", logger.handleName),
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Fatal = log.New(fatalHandle, fmt.Sprintf("%s - FATAL: ", logger.handleName),
		log.Ldate|log.Ltime|log.Lshortfile)
	logger.Log = log.New(basicHandle, "", 0)
}

func (logger *Logger) Logging(req *http.Request, statuscode int) {
	logger.Info.Println(fmt.Sprintf(`%s - [%s] "%s %s %s" %d "%s" "%s"`,
		req.Host,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		req.Method,
		req.RequestURI,
		req.Proto,
		statuscode,
		req.Referer(),
		req.UserAgent()))
	logger.Log.Println(fmt.Sprintf(`%s - [%s] "%s %s %s" %d "%s" "%s"`,
		req.Host,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		req.Method,
		req.RequestURI,
		req.Proto,
		statuscode,
		req.Referer(),
		req.UserAgent()))
}
