package main

import (
	"io"
	"io/ioutil"
	"log"
)

const (
	// The log prefix for output logging
	logPrefix = "8queens "
)

var (
	// The global logging reciever
	l *log.Logger
)

/* init initializes the variables in this source file
 */
func init() {
	l = log.New(ioutil.Discard, logPrefix, log.Lshortfile|log.Ltime)
}

/* EnableLogging takes in a destination io.Writer to write all log information to
 * for this package
 * @param dst the writer to writer all logs to
 */
func EnableLogging(dst io.Writer) {
	l.SetOutput(dst)
}
