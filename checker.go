package main

import (
	"fmt"
	"runtime"
)

var currOs = runtime.GOOS

type Checker interface {
	CheckTermination(string) bool
	ValidateInput(string) bool
}

// This would be used to check inputs
type NumberChecker struct {
	// making the fields private to let validation to Setters
	termination string
	numLimit    int
}

// Creates a NumberChecker with Termination string == "terminate"
// (and includes carriage) and NumLimit == 9
func NewDefaultNumberChecker() *NumberChecker {
	defaultTerminate := fmt.Sprintf("%s%s", "terminate", LINE_BREAK)
	return &NumberChecker{termination: defaultTerminate, numLimit: 9}
}

// Custom setter for termination word
// it adds the system's carriage character to the end of the string
// and assigns it to the NumberCheker
func (nc *NumberChecker) SetTermination(newTerminate string) {
	nc.termination = fmt.Sprintf("%s%s", newTerminate, LINE_BREAK)
}

// The limit of digits a string number can have
// it errors out if the newLimit is negative
func (nc *NumberChecker) SetNumLimit(newLimit int) error {
	if newLimit < 0 {
		return fmt.Errorf("NumLimit can't be a negative number: %d", newLimit)
	}
	nc.numLimit = newLimit
	return nil
}

// Basic getter of NumLimit of the NumberChecker
func (nc *NumberChecker) GetNumLimit() int {
	return nc.numLimit
}

// Basic getter of the Termination string of this NumberChecker
func (nc *NumberChecker) GetTermination() string {
	return nc.termination
}

// Basic Validation functions

// Checks the input passed and indicates whether it
// corresponds to the terminate string (input == nc.terminate) order.
func (nc *NumberChecker) CheckTermination(input string) bool {
	return input == nc.termination
}

// Validates whether the passed string corresponds
// to the expected format on the input numbers expected by
// the server. It also accounts for carriage character
// (this depends on the underlaying OS)
// As for this implementation, a number would be expected,
// with length 9 characters, by default (if not set differently
// in the Checker instance)
func (nc *NumberChecker) ValidateInput(input string) bool {
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
