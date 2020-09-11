package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type newNumberTrackerCase struct {
	Name    string
	Cap     int
	Errored bool
}

type processNumberCase struct {
	Name     string
	Canceled bool
	Repeated bool
	Inbound  int
	Outbound int
}

func TestNumberTracker(t *testing.T) {
	t.Run("New Number Tracker", func(t *testing.T) {
		genericError := "Got: %v, Expected: %v"
		testCases := []newNumberTrackerCase{
			{
				Name: "Positive cap",
				Cap:  999999999,
			},
			{
				Name:    "Negative cap",
				Cap:     -123,
				Errored: true,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				tracker, err := NewNumberTracker(tc.Cap)
				errored := err != nil
				// It will stop evaluation if this fails
				require.True(t, tc.Errored == errored, genericError, errored, tc.Errored)
				if errored {
					assert.True(t, errors.Is(err, BadMaxCapacity), "Expected BadMaxCapacity, got %v", err)
				} else {
					maxCap := len(tracker.KnownNumbers)
					assert.True(t, maxCap == tc.Cap, genericError, maxCap, tc.Cap)
				}
			})
		}
	})

	t.Run("Process Number", func(t *testing.T) {
		tracker, err := NewNumberTracker(999999999)
		if err != nil {
			t.Error(err)
		}
		genericError := "Got: %v, Expected: %v"
		testCases := []processNumberCase{
			{
				Name:     "Canceled context",
				Canceled: true,
				Inbound:  100,
			},
			{
				Name:    "Ignored overflown int",
				Inbound: 1000000000,
			},
			{
				Name:    "Ignored negative int",
				Inbound: -200,
			},
			{
				Name:     "Correct new int",
				Inbound:  300,
				Outbound: 300,
			},
			{
				Name:     "Igonored repeated int",
				Inbound:  100,
				Repeated: true,
			},
		}
		// Helper function to be used for checking
		// outbound channel state
		checkChan := func(outbound <-chan int) bool {
			// If more than 2 second passes, return
			ticker := time.Tick(time.Second * 2)
			for {
				select {
				case <-outbound:
					return false
				case <-ticker:
					return true
				}
			}
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				inbound := make(chan int)
				outbound := tracker.ProcessNumber(ctx, inbound)
				inbound <- tc.Inbound
				if tc.Canceled {
					cancel()
					v, ok := <-outbound
					assert.True(t, v == 0, genericError, v, 0)
					assert.False(t, ok, genericError, ok, false)
				} else {
					received := <-outbound
					if !tc.Repeated {
						assert.True(t, tc.Inbound == received, genericError, received, tc.Inbound)
					} else {
						inbound <- tc.Inbound
						timedout := checkChan(outbound)
						assert.True(t, timedout, genericError, timedout, false)
					}
				}
			})
		}
	})
}
