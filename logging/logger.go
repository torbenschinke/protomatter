package logging

import "log"

type Logger interface {
	Println(v ...interface{})
}

func NewLogger() DefaultLogger {
	return DefaultLogger{}
}

type DefaultLogger struct {
}

func (DefaultLogger) Println(v ...interface{}) {
	log.Println(v...)
}
