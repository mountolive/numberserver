package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
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
	Filename string
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
					newFilename = DEFAULT_LOG_FILE
					appenderFlag = tc.Appender
					logger = NewLogger(Appender(tc.Appender))
				} else {
					newFilename = DEFAULT_LOG_FILE
					logger = NewLogger()
				}
				// It will stop evaluation if this fails
				require.True(t, newFilename == logger.filename,
					genericError, newFilename, logger.filename)
				assert.True(t, appenderFlag == logger.appender,
					genericError, appenderFlag, logger.appender)
			})
		}
	})

	t.Run("Stream Write", func(t *testing.T) {
		genericError := "Got %v, Expected %v"

		// Checks for subset condition
		isSubset := func(got, expected []string) bool {
			countSet := make(map[string]int)
			for _, expectedWord := range expected {
				countSet[expectedWord] = 1
			}
			for _, gotWord := range got {
				if countSet[gotWord] == 0 {
					return false
				}
			}
			return true
		}

		// Reads the created file and checks if the lines were
		// correctly written
		wroteChecker := func(path string, linesExpected []string) (bool, error) {
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
			sizeDiff := len(linesExpected) - len(lines)
			// If negative, means we wrote more lines than expected
			// If > 1, means we're missing writing more than the expected last line
			// (After context cancellation)
			if sizeDiff < 0 || sizeDiff > 1 {
				return false, fmt.Errorf("Lines scanned: %d vs. Lines wrote: %d",
					len(lines), len(linesExpected))
			}
			// Checking subset condition
			return isSubset(lines, linesExpected), nil
		}

		testCases := []streamWriteCase{
			{
				Name:  "Write 1 line",
				Lines: []string{"one"},
			},
			{
				Name:  "Write several lines",
				Lines: []string{"one", "two", "three"},
			},
			{
				Name:     "Append one line",
				Lines:    []string{"one", "two"},
				Appender: true,
				// Using a brand new file for the test
				Filename: "./appender1.log",
			},
			{
				Name:     "Append several lines",
				Lines:    []string{"one", "two", "three"},
				Appender: true,
				Filename: "./appender2.log",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				logger := NewLogger(Appender(tc.Appender))
				if tc.Filename != "" {
					logger.setFilename(tc.Filename)
				}
				// Cleaning up the filesystem
				defer os.Remove(logger.filename)
				// Will be closed in the if-else blocks, below
				readStream := make(chan string)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				// Check what happens when we open an existing file for appending
				if tc.Appender {
					err := logger.StreamWrite(ctx, readStream)
					if err != nil {
						t.Error(err)
					}
					readStream <- tc.Lines[0]
					// Cancelling the context to stop the proccess
					cancel()
					// Recreating the stream digesting again
					err = logger.StreamWrite(context.Background(), readStream)
					if err != nil {
						t.Error(err)
					}
					for _, line := range tc.Lines[1:] {
						readStream <- line
					}
					// closing the channel
					close(readStream)
					check, err := wroteChecker(logger.filename, tc.Lines)
					if err != nil {
						t.Error(err)
					}
					assert.True(t, check, genericError, check, true)
					// Writing file from scratch
				} else {
					err := logger.StreamWrite(ctx, readStream)
					if err != nil {
						t.Error(err)
					}
					for _, line := range tc.Lines {
						readStream <- line
					}
					// closing the channel
					close(readStream)
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
