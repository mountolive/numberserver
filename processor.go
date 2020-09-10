package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

const LOG_FILE = "./numbers.log"

var terminate = fmt.Sprintf("%s%s", "terminate", LINE_BREAK)
var currOs = runtime.GOOS

type Appender interface {
	AddToLog(ctx context.Context, path, input string) error
}

type Statistics struct {
	sync.Mutex
	Received   int
	Duplicates int
	Total      int
}

// Prints to STDOUT the current statistics of the server,
// regarding received numbers, number of duplicates and
// total number of unique numbers received by the server
func (s *Statistics) PrintCurrent() {
	s.Lock()
	fmt.Printf("Received: %d unique numbers, %d duplicates, "+
		"Unique totals: %d \n", s.Received, s.Duplicates, s.Total)
	s.Unlock()
}

// Updates Total statistics based on the values recv and dups,
// Received and Duplicates. If one of them is negative, it
// silently exits
func (s *Statistics) Update(recv, dups int) {
	if recv < 0 || dups < 0 {
		return
	}
	if dups > recv {
		return
	}
	s.Lock()
	defer s.Unlock()
	s.Received = recv
	s.Duplicates = dups
	s.Total = s.Total + (recv - dups)
}

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
