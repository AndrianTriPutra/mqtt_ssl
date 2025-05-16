package logger

import (
	"errors"
	"log"
)

var ErrorNotFound = errors.New("content not found")

var Access uint8 = 6

func Load(logging string) {
	switch logging {
	case "panic":
		Access = 0
	case "fatal":
		Access = 1
	case "error":
		Access = 2
	case "warning":
		Access = 3
	case "info":
		Access = 4
	case "debug":
		Access = 5
	case "trace":
		Access = 6
	}
}

func Level(level string, loc string, data interface{}) {
	switch level {
	case "panic":
		switch Access {
		case 0, 1, 2, 3, 4, 5, 6:
			log.Printf("[ PANIC ] [%s] -> ", loc)
			panic(data)
		}
	case "fatal":
		switch Access {
		case 1, 2, 3, 4, 5, 6:
			log.Fatalf("[ FATAL ] [%s] -> %v", loc, data)
		}
	case "error":
		switch Access {
		case 2, 3, 4, 5, 6:
			log.Printf("[ ERROR ] [%s] -> %s", loc, data)
		}
	case "warning":
		switch Access {
		case 3, 4, 5, 6:
			log.Printf("[WARNING] [%s] -> %s", loc, data)
		}
	case "info":
		switch Access {
		case 4, 5, 6:
			log.Printf("[  INFO ] [%s] -> %s", loc, data)
		}
	case "debug":
		switch Access {
		case 5, 6:
			log.Printf("[ DEBUG ] [%s] -> %s", loc, data)
		}
	}
}
