package main

import (
	"io"
	"io/ioutil"
	"log"
)

const (
	logPrefix = "8queens "
)

var (
	l *log.Logger
)

// Necessary because we don't want to call functions on a nil pointer :D
func init() {
	l = log.New(ioutil.Discard, logPrefix, log.Lshortfile|log.Ltime)
}

// EnableLogging allows users of this library to enable logging to a specific io.Writer
func EnableLogging(dst io.Writer) {
	l.SetOutput(dst)
}
