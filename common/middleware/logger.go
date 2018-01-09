package middleware

import (
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func RollingLog(filename string) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   filename,
		MaxSize:    20, // megabytes
		MaxBackups: 5,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	}
}
