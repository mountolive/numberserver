package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type processNumberCase struct {
	Name     string
	Canceled bool
	Repeated bool
	Ignored  bool
	Inbound  int
	Outbound string
}

func TestNumberTracker(t *testing.T) {
	t.Run("Process Number", func(t *testing.T) {
		tracker := NewNumberTracker()
		genericError := "Got: %v, Expected: %v"
		testCases := []processNumberCase{
			{
				Name:     "Canceled context",
				Canceled: true,
				Inbound:  100,
			},
			{
				Name:    "Ignored negative int",
				Inbound: -200,
				Ignored: true,
			},
			{
				Name:     "Correct new int",
				Inbound:  300,
				Outbound: "300",
			},
			{
				Name:     "Igonored repeated int",
				Inbound:  100,
				Repeated: true,
			},
		}
		// Helper function to be used for checking
		// outbound channel state
		checkChan := func(outbound <-chan string) (bool, string) {
			// Hard limit
			// If more than 200 millisecond passes, return
			ticker := time.Tick(time.Millisecond * 200)
			for {
				select {
				case value := <-outbound:
					return false, value
				case <-ticker:
					return true, ""
				}
			}
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				inbound := make(chan int)
				defer close(inbound)
				outbound := tracker.ProcessNumber(ctx, inbound)
				if tc.Canceled {
					cancel()
					inbound <- tc.Inbound
					v, ok := <-outbound
					assert.True(t, v == "", genericError, v, "")
					assert.False(t, ok, genericError, ok, false)
				} else {
					inbound <- tc.Inbound
					if !tc.Repeated {
						timedout, received := checkChan(outbound)
						if tc.Ignored {
							require.True(t, timedout, "Should have timedout out (ignored int)")
							assert.True(t, received == "", genericError, received, 0)
						} else {
							assert.True(t, tc.Outbound == received, genericError, received, tc.Inbound)
						}
					} else {
						<-outbound
						inbound <- tc.Inbound
						timedout, _ := checkChan(outbound)
						assert.True(t, timedout, genericError, timedout, false)
					}
				}
			})
		}
	})
}
