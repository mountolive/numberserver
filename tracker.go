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
	KnownNumbers map[int]byte
	Stats        *Statistics
}

// Creates a new NumberTracker.
// It contains a set-book for known, found numbers,
// and a Statistics tracker.
func NewNumberTracker() *NumberTracker {
	return &NumberTracker{KnownNumbers: make(map[int]byte), Stats: &Statistics{}}
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
					if n.checkUniqueness(input) {
						// Marking it as seen
						n.registerNumber(input)
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

func (n *NumberTracker) registerNumber(input int) {
	// Locking reading for consistency
	// Any subsequent read will have the proper state
	n.RLock()
	defer n.RUnlock()
	n.KnownNumbers[input] = 1
}

func (n *NumberTracker) checkUniqueness(input int) bool {
	// Locking writing for consistency
	n.Lock()
	defer n.Unlock()
	return n.KnownNumbers[input] == 0
}
