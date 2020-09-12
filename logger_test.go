package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type newLoggerCase struct {
	Name     string
	Filename string
	Appender bool
}

type streamWriteCase struct {
	Name     string
	Appender bool
	Lines    []string
}

func TestLogger(t *testing.T) {
	t.Run("New Logger", func(t *testing.T) {
		genericError := "Got %v, Expected %v"
		testCases := []newLoggerCase{
			{
				Name:     "Appender and Filename",
				Filename: "./other.log",
				Appender: true,
			},
			{
				Name:     "Only Appender",
				Appender: true,
			},
			{
				Name:     "Only Filename",
				Filename: "./other.log",
			},
			{
				Name: "Default",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				var logger *Logger
				var newFilename string
				var appenderFlag bool
				if tc.Filename != "" && tc.Appender {
					newFilename = tc.Filename
					appenderFlag = tc.Appender
					logger = NewLogger(Filename(tc.Filename), Appender(tc.Appender))
				} else if tc.Filename != "" {
					newFilename = tc.Filename
					appenderFlag = tc.Appender
					logger = NewLogger(Filename(tc.Filename))
				} else if tc.Appender {
					newFilename = LOG_FILE
					appenderFlag = tc.Appender
					logger = NewLogger(Appender(tc.Appender))
				} else {
					newFilename = LOG_FILE
					logger = NewLogger()
				}
				require.True(t, newFilename == logger.filename,
					genericError, newFilename, logger.filename)
				assert.True(t, appenderFlag == logger.appender,
					genericError, appenderFlag, logger.appender)
			})
		}
	})

	t.Run("Stream Write", func(t *testing.T) {
		genericError := "Got %v, Expected %v"
		// Reads the created file and checks if the lines were
		// correctly written
		wroteChecker := func(path string, wroteLines []string) (bool, error) {
			file, err := os.Open(path)
			defer file.Close()
			if err != nil {
				return false, fmt.Errorf("An error occurred while reading the log file")
			}
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			var lines []string
			// Created lines
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			// Comparing lines
			result := true
			sort.Strings(lines)
			sort.Strings(wroteLines)
			for i := 0; i < len(lines); i++ {
				result = lines[i] == wroteLines[i]
			}
			return result, nil
		}
		testCases := []streamWriteCase{
			{
				Name:  "Write 1 line",
				Lines: []string{"one"},
			},
			{
				Name:  "Write several line",
				Lines: []string{"one", "two", "three"},
			},
			{
				Name:     "Append one line",
				Lines:    []string{"one", "two"},
				Appender: true,
			},
			{
				Name:     "Append several line",
				Lines:    []string{"one", "two", "three"},
				Appender: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				logger := NewLogger(Appender(tc.Appender))
				// Giving it capacity so that it doens't block the write
				readStream := make(chan string, 10)
				defer close(readStream)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tc.Appender {
					logger.StreamWrite(ctx, readStream)
					readStream <- tc.Lines[0]
					// Cancelling the context to stop the proccess
					cancel()
					// Recreating the stream digesting again
					ctx, cancel = context.WithCancel(context.Background())
					logger.StreamWrite(ctx, readStream)
					for _, line := range tc.Lines[1:] {
						readStream <- line
					}
					check, err := wroteChecker(logger.filename, tc.Lines)
					if err != nil {
						t.Error(err)
					}
					assert.True(t, check, genericError, check, true)
				} else {
					logger.StreamWrite(ctx, readStream)
					for _, line := range tc.Lines {
						readStream <- line
					}
					check, err := wroteChecker(logger.filename, tc.Lines)
					if err != nil {
						t.Error(err)
					}
					assert.True(t, check, genericError, check, true)
				}
			})
		}
	})
}
