package ulogger

import (
	"bytes"
	"fmt"
)

// Error displays an error message
func (log *Logger) Error(args ...interface{}) {
	if log.logLevelCode > 4 {
		if log.RemoteAvailable {
			// Create the logMessage struct here
			logStruct, _ := generateTimestamp("ERROR")
			ch := make(chan int)
			go sendLogMessageFromWrite(logStruct, ch, args...)
			<-ch
		}
		return
	}
	ch := make(chan int)
	go func(ch chan int) {
		write(errorPrefix, log, log.ErrorColor, ch, args...)
	}(ch)
	<-ch
}

// Errorf displays an error message
func (log *Logger) Errorf(format string, args ...interface{}) {
	if log.logLevelCode > 4 {

		if log.RemoteAvailable {
			// Create the logMessage struct here
			logStruct, _ := generateTimestamp("ERROR")
			ch := make(chan int)
			go sendLogMessageFromWritef(logStruct, ch, format, args...)
			<-ch
		}
		return
	}
	ch := make(chan int)
	go func(ch chan int) {
		writef(errorPrefix, log, log.ErrorColor, ch, format, args...)
	}(ch)
	<-ch
}

// Errorln displays an error message
func (log *Logger) Errorln(args ...interface{}) {
	if log.logLevelCode > 4 {

		if log.RemoteAvailable {
			// Create the logMessage struct here
			logStruct, _ := generateTimestamp("ERROR")
			ch := make(chan int)
			go sendLogMessageFromWriteln(logStruct, ch, args...)
			<-ch
		}
		return
	}
	ch := make(chan int)
	go func(ch chan int) {
		writeln(errorPrefix, log, log.ErrorColor, ch, args...)
	}(ch)
	<-ch
}

// ErrorDump displays the dump of the variables passed using the go-spew library
func (log *Logger) ErrorDump(args ...interface{}) {
	// Don't stream this to the remote server
	ch := make(chan int)
	go func(ch chan int) {
		writeDump(errorPrefix, log, log.ErrorColor, ch, args...)
	}(ch)
	<-ch
}

// Returns a string, along with a logMessage after prefixing the timestamp and the type of log
func errorPrefix(log *Logger) (*bytes.Buffer, logMessage) {
	buf := new(bytes.Buffer)
	logStruct, timestamp := generateTimestamp("ERROR")
	logStruct.OrganizationName = log.OrganizationName
	logStruct.ApplicationName = log.ApplicationName
	log.ErrorTimeColor.Fprint(buf, timestamp.Format(timeFormat))
	fmt.Fprint(buf, " ")
	log.ErrorMessageTypeColor.Fprint(buf, logStruct.MessageType)
	fmt.Fprint(buf, " ")
	return buf, logStruct
}
