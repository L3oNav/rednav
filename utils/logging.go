package utils

import "log"

var logger *log.Logger

func init() {
	logger = log.New(log.Writer(), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
