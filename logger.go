package main

import (
	"context"
	"fmt"
	"os"
)

const LOG_FILE = "./numbers.log"

type Logger struct {
	filename string
	appender bool
}

// Creates a new Logger. If no option is passed, creates it
// so that the filename to log would be ./numbers.log
// and will overwrite the file (if existing) for every fresh start
// (using options allows us to avoid creating several New* functions)
// example usages: NewLogger(Filename("new.log"), Appender(true))
//                 NewLogger(Appender(true))
func NewLogger(options ...func(*Logger)) *Logger {
	defaultLogger := &Logger{filename: LOG_FILE, appender: false}
	for _, option := range options {
		option(defaultLogger)
	}
	return defaultLogger
}

// Option for setting a logger's filename
func Filename(name string) func(*Logger) {
	return func(logger *Logger) {
		logger.setFilename(name)
	}
}

// Option for setting whether a logger should append
// or create new files everytime it starts writing logs
func Appender(appender bool) func(*Logger) {
	return func(logger *Logger) {
		logger.setAppender(appender)
	}
}

// Writes streamed input to the configured log file,
// adding a carriage new-line character to the end of each passed line
// throws error if file doesn't exist
func (l *Logger) StreamWrite(ctx context.Context, streamLines <-chan string) error {
	var file *os.File
	var err error
	// Checking context before opening file
	select {
	case <-ctx.Done():
		return fmt.Errorf("Context passed to stream writer is cancelled")
	default:
		if l.appender {
			// Appending existing file, creating if it doesn't exist
			file, err = os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			file, err = os.Create(l.filename)
		}
		if err != nil {
			return fmt.Errorf("An error occurred while retrieving/creating the logfile: %w", err)
		}
	}
	// Start consuming input
	go func() {
		defer file.Close()
		for line := range streamLines {
			lineToPrint := []byte(fmt.Sprintf("%v%s", line, LINE_BREAK))
			select {
			case <-ctx.Done():
				// Attempting writing lastly received line
				fmt.Println("Writing last input, before exiting")
				if _, err = file.Write(lineToPrint); err != nil {
					fmt.Printf("An error occurred while writing to the file: %v\n", err)
					return
				}
				fmt.Println("Canceled writing")
				return
			default:
				if _, err = file.Write(lineToPrint); err != nil {
					fmt.Printf("An error occurred while writing to the file: %v\n", err)
					return
				}
			}
		}
	}()
	return nil
}

// Sets new filename to be written by the logger
func (l *Logger) setFilename(name string) {
	l.filename = name
}

// Sets whether the logger should create a file or
// append to an existing one
func (l *Logger) setAppender(appender bool) {
	l.appender = appender
}
