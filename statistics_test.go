package main

import (
	"math/rand"
	"testing"
	"testing/quick"
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
		// Triggering only a few cases to avoid cluttering of the STDOUT
		if err := quick.Check(asserter, &quick.Config{MaxCount: 15}); err != nil {
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
