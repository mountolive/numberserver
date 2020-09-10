package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type inputTestCase struct {
	Name     string
	Input    string
	Expected bool
	ErrorMsg string
}

func TestProcessor(t *testing.T) {
	t.Run("Validate Input Num", func(t *testing.T) {
		testCases := []inputTestCase{
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

	t.Run("Check Termination", func(t *testing.T) {
		generalErrorMsg := "The passed word %v should have prompted %v"
		testCases := []inputTestCase{
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
