package main

import (
	"context"
	"strconv"
	"sync"
)

// Keeps a set of processed numbers
// and statistics book
type NumberTracker struct {
	sync.RWMutex
	// Max value of Uint32 = 4294967295
	KnownNumbers map[uint32]bool
	Stats        *Statistics
}

// Creates a new NumberTracker.
// It contains a set-book for known, found numbers,
// and a Statistics tracker.
func NewNumberTracker() *NumberTracker {
	return &NumberTracker{KnownNumbers: make(map[uint32]bool), Stats: &Statistics{}}
}

// Processes a number, validates and passes it on to a channel
// in a pipelined fashion (after converting it to a string)
func (n *NumberTracker) ProcessNumber(ctx context.Context,
	inputStream <-chan int) <-chan string {
	output := make(chan string)
	go func() {
		defer close(output)
		for input := range inputStream {
			select {
			case <-ctx.Done():
				return
			default:
				if input >= 0 {
					uintVal := uint32(input)
					if n.checkUniqueness(uintVal) {
						// Marking it as seen
						n.registerNumber(uintVal)
						// passing it on
						output <- strconv.Itoa(input)
						// Increasing unique received count
						n.Stats.IncreaseReceived()
					} else {
						n.Stats.IncreaseDups()
					}
				}
			}
		}
	}()
	return output
}

// Printing current statistics' state
func (n *NumberTracker) PrintStatistics() {
	n.Stats.PrintCurrent()
}

func (n *NumberTracker) registerNumber(input uint32) {
	// Locking reading for consistency
	// Any subsequent read will have the proper state
	n.RLock()
	defer n.RUnlock()
	n.KnownNumbers[input] = true
}

func (n *NumberTracker) checkUniqueness(input uint32) bool {
	// Locking writing for consistency
	n.Lock()
	defer n.Unlock()
	return !n.KnownNumbers[input]
}
