package main

import "context"

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

// Writes streamed input to the configured log file
// throws error if file doesn't exist
func (l *Logger) StreamWrite(ctx context.Context, streamLines <-chan string) {
	// TODO Implement
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
