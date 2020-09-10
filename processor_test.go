package main

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

type validateInputTestCase struct {
	Name     string
	Input    string
	Expected bool
	ErrorMsg string
}

type statisticsTestCase struct {
	Name       string
	Received   int
	Duplicates int
	Expected   int
}

func TestStatistics(t *testing.T) {
	t.Run("Update Basic Correctness", func(t *testing.T) {
		errMsg := "Got %d, Expected: %d"
		testCases := []statisticsTestCase{
			{
				Name:       "No change",
				Received:   1,
				Duplicates: 1,
				Expected:   100,
			},
			{
				Name:       "Negative Received",
				Received:   -1,
				Duplicates: 30,
				Expected:   100,
			},
			{
				Name:       "Negative Duplicates",
				Received:   20,
				Duplicates: -2,
				Expected:   100,
			},
			{
				Name:       "Changed 1",
				Received:   33,
				Duplicates: 13,
				Expected:   120,
			},
			{
				Name:       "Changed 2",
				Received:   456,
				Duplicates: 56,
				Expected:   500,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				s := &Statistics{Total: 100}
				s.Update(tc.Received, tc.Duplicates)
				assert.True(t, tc.Expected == s.Total, errMsg, s.Total, tc.Expected)
			})
		}
	})

	t.Run("Update Many", func(t *testing.T) {
		s := &Statistics{Total: 100}
		asserter := func(recv, dups int) bool {
			curr := s.Total
			if recv < 0 || dups < 0 {
				return curr == s.Total
			}
			if dups > recv {
				return curr == s.Total
			}
			if s.Total > 999999999 {
				return true
			}
			s.Update(recv, dups)
			return s.Total > 0
		}
		if err := quick.Check(asserter, &quick.Config{MaxCount: 100}); err != nil {
			t.Error(err)
		}
	})
}

func TestValidateInputNum(t *testing.T) {
	t.Run("Validate Input", func(t *testing.T) {
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
				input := fmt.Sprintf("%s%s", tc.Input, LINE_BREAK)
				result := ValidateInputNum(input)
				assert.True(t, tc.Expected == result, tc.ErrorMsg, tc.Input, tc.Expected)
			})
		}
	})

}

func TestCheckTermination(t *testing.T) {
	t.Run("Check Termination", func(t *testing.T) {
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
				input := fmt.Sprintf("%s%s", tc.Input, LINE_BREAK)
				result := CheckTermination(input)
				assert.True(t, tc.Expected == result, generalErrorMsg, tc.Input, tc.Expected)
			})
		}
	})

}
