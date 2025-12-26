package main

import "fmt"

/*
HANDLER
*/
type Logger interface {
	SetNext(Logger)
	Log(level string, message string)
}

/*
BASE HANDLER
*/
type BaseLogger struct {
	next Logger
}

func (b *BaseLogger) SetNext(next Logger) {
	b.next = next
}

/*
CONCRETE HANDLERS
*/
type InfoLogger struct{ BaseLogger }

func (l *InfoLogger) Log(level, message string) {
	if level == "INFO" {
		fmt.Println("INFO:", message)
	} else if l.next != nil {
		l.next.Log(level, message)
	}
}

type ErrorLogger struct{ BaseLogger }

func (l *ErrorLogger) Log(level, message string) {
	if level == "ERROR" {
		fmt.Println("ERROR:", message)
	} else if l.next != nil {
		l.next.Log(level, message)
	}
}

func main() {
	info := &InfoLogger{}
	errorLogger := &ErrorLogger{}

	info.SetNext(errorLogger)

	info.Log("INFO", "Application started")
	info.Log("ERROR", "Something failed")
}
