package main

import (
	"context"
	"fmt"
	"runtime"
)

const LOG_FILE = "./numbers.log"

type Appender interface {
	AddToLog(ctx context.Context, path, input string) error
}

var terminate = fmt.Sprintf("%s%s", "terminate", LINE_BREAK)
var currOs = runtime.GOOS

// Basic Validation functions

// Checks the input passed and indicates whether it
// corresponds to a terminate (input == terminate) order.
func CheckTermination(input string) bool {
	return input == terminate
}

// Validates whether the passed string corresponds
// to the expected format on the input numbers expected by
// the server. It should be 9 characters long, ended by
// carriage character (this depends on the underlaying OS)
func ValidateInputNum(input string) bool {
	substr := 1
	if currOs == "windows" {
		substr = 2
	}
	newEnd := len(input) - substr
	if newEnd != 9 {
		return false
	}
	// Omitting the carriage in the evaluation
	for _, char := range input[:newEnd] {
		value := int(char)
		// '0' == 48 and '9' == 57
		if value < 48 || value > 57 {
			return false
		}
	}
	return true
}
