package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type validateInputNumCase struct {
	Name     string
	Input    string
	Expected bool
	ErrorMsg string
}

func TestProcessor(t *testing.T) {
	t.Run("Validate Input", func(t *testing.T) {
		testCases := []validateInputNumCase{
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
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				result := ValidateInputNum(tc.Input)
				assert.True(t, tc.Expected == result, tc.ErrorMsg, tc.Input, tc.Expected)
			})
		}
	})
}
