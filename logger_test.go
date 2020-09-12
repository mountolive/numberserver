package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type newLoggerCase struct {
	Name     string
	Filename string
	Appender bool
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
	})
}
