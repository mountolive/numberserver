package main

import (
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

type bulkUpdateTestCase struct {
	Name       string
	Received   int
	Duplicates int
	Expected   int
}

func TestStatistics(t *testing.T) {
	t.Run("PrintCurrent", func(t *testing.T) {
		s := &Statistics{Total: 100, Received: 12, Duplicates: 32}
		asserter := func() bool {
			limit := rand.Intn(40)
			for i := 0; i < limit; i++ {
				if limit%2 == 0 {
					s.IncreaseDups()
				}
				s.IncreaseReceived()
			}
			s.PrintCurrent()
			return s.Duplicates == 0 && s.Received == 0
		}
		if err := quick.Check(asserter, &quick.Config{MaxCount: 100}); err != nil {
			t.Error(err)
		}

	})
	t.Run("BulkUpdate Basic Correctness", func(t *testing.T) {
		errMsg := "Got %d, Expected: %d"
		testCases := []bulkUpdateTestCase{
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
				s.BulkUpdate(tc.Received, tc.Duplicates)
				assert.True(t, tc.Expected == s.Total, errMsg, s.Total, tc.Expected)
			})
		}
	})

	t.Run("BulkUpdate Many", func(t *testing.T) {
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
			s.BulkUpdate(recv, dups)
			return s.Total > 0
		}
		if err := quick.Check(asserter, &quick.Config{MaxCount: 100}); err != nil {
			t.Error(err)
		}
	})

	t.Run("Increse Duplicates", func(t *testing.T) {
		s := &Statistics{Total: 100}
		asserter := func() bool {
			previous := s.Total
			s.IncreaseDups()
			return s.Total == previous
		}
		if err := quick.Check(asserter, &quick.Config{MaxCount: 100}); err != nil {
			t.Error(err)
		}
	})

	t.Run("Increse Received", func(t *testing.T) {
		s := &Statistics{Total: 100}
		asserter := func() bool {
			expected := s.Total + 1
			s.IncreaseReceived()
			return s.Total == expected
		}
		if err := quick.Check(asserter, &quick.Config{MaxCount: 100}); err != nil {
			t.Error(err)
		}
	})
}
