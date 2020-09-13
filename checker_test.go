package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type validateInputTestCase struct {
	Name     string
	Input    string
	Expected bool
	ErrorMsg string
}

type setTerminationCase struct {
	Name         string
	NewTerminate string
	Expected     string
}

type setNumLimitCase struct {
	Name     string
	NumLimit int
	Errored  bool
}

type getIntValueCase struct {
	Name     string
	Input    string
	Expected int
	Errored  bool
}

func TestNumberChecker(t *testing.T) {
	t.Run("Canary test", func(t *testing.T) {
		var _ Checker = &NumberChecker{}
	})

	t.Run("Set Terminate", func(t *testing.T) {
		genericError := "Got: %s, Expected: %s"
		numberChecker := NewDefaultNumberChecker()
		testCases := []setTerminationCase{
			{
				Name:         "Basic termination 1",
				NewTerminate: "I'm thinking of ending things",
				Expected:     "I'm thinking of ending things",
			},
			{
				Name:         "Basic termination 2",
				NewTerminate: "This is the end",
				Expected:     "This is the end",
			},
		}
		for _, tc := range testCases {
			numberChecker.SetTermination(tc.NewTerminate)
			newTerminate := numberChecker.GetTermination()
			assert.True(t, newTerminate == tc.Expected, genericError, newTerminate, tc.Expected)
		}
	})

	t.Run("Set Num Limit", func(t *testing.T) {
		genericError := "Got: %v, Expected: %v"
		numberChecker := NewDefaultNumberChecker()
		testCases := []setNumLimitCase{
			{
				Name:     "Correct limit",
				NumLimit: 23,
			},
			{
				Name:     "Negative limit",
				NumLimit: -6,
				Errored:  true,
			},
		}
		for _, tc := range testCases {
			isErrored := numberChecker.SetNumLimit(tc.NumLimit) != nil
			if isErrored {
				assert.True(t, isErrored == tc.Errored, genericError, isErrored, tc.Errored)
			} else {
				newNumLimit := numberChecker.GetNumLimit()
				assert.True(t, newNumLimit == tc.NumLimit, genericError, newNumLimit, tc.NumLimit)
			}
		}
	})

	t.Run("Check Termination", func(t *testing.T) {
		numberChecker := NewDefaultNumberChecker()
		generalErrorMsg := "The passed word %v should have prompted %v"
		testCases := []validateInputTestCase{
			{
				Name:     "Correct termination",
				Input:    "terminate",
				Expected: true,
			},
			{
				Name:     "Not termination",
				Input:    "anotherword",
				Expected: false,
			},
			{
				Name:     "Composed, bad 1",
				Input:    "terminate hello",
				Expected: false,
			},
			{
				Name:     "Composed, bad 2",
				Input:    "hello terminate",
				Expected: false,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				result := numberChecker.CheckTermination(tc.Input)
				assert.True(t, tc.Expected == result, generalErrorMsg, tc.Input, tc.Expected)
			})
		}
	})

	t.Run("Validate Input", func(t *testing.T) {
		numberChecker := NewDefaultNumberChecker()
		testCases := []validateInputTestCase{
			{
				Name:     "Valid Input Num 1",
				Input:    "314159265",
				Expected: true,
				ErrorMsg: "The number passed %s should have been %v",
			},
			{
				Name:     "Valid Input Num 2",
				Input:    "007007009",
				Expected: true,
				ErrorMsg: "The number passed %s should have been %v",
			},
			{
				Name:     "Incomplete number",
				Input:    "00700700",
				Expected: false,
				ErrorMsg: "The number passed %s should have been %v",
			},
			{
				Name:     "Non-numeric string",
				Input:    "testing",
				Expected: false,
				ErrorMsg: "The string passed %s should have been %v",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				result := numberChecker.ValidateInput(tc.Input)
				assert.True(t, tc.Expected == result, tc.ErrorMsg, tc.Input, tc.Expected)
			})
		}
	})
}
