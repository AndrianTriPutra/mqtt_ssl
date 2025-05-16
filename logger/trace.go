package logger

import (
	"log"
)

func Trace(Header, value interface{}) {
	if Access == 6 {
		log.Printf("[TRACE] %s %v", Header, value)
	}
}
