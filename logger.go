package main

import "log"

var (
	LOG_ERROR   = 3
	LOG_INFO    = 2
	LOG_DEBUG   = 1
	LOG_WARNING = 0
)

func Log(lvl int, args ...interface{}) {
	var level string
	switch lvl {
	case 0:
		level = "[WARNING]"
	case 1:
		level = "[DEBUG]"
	case 2:
		level = "[INFO]"
	case 3:
		level = "[ERROR]"
	default:
		level = "[DEBUG]"
	}
	if ShouldLog {
		log.Println(level, args)
	}
}
