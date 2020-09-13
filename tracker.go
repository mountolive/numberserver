package main

import (
	"context"
	"errors"
	"strconv"
	"sync"
)

var BadMaxCapacity = errors.New("Max Capacity for Number Tracker can't be negative")

// Keeps a set of processed numbers
// and statistics book
type NumberTracker struct {
	sync.RWMutex
	KnownNumbers []byte
	Stats        *Statistics
}

// Creates a new NumberTracker.
// maxCapacity should be the maximum possible number
// that can be present known by the tracker
func NewNumberTracker(maxCapacity int) (*NumberTracker, error) {
	numTracker := &NumberTracker{}
	if maxCapacity < 0 {
		return numTracker, BadMaxCapacity
	}
	// Adding one to account for the maximum number possible
	// 999999999 in the default example (indexing starts at 0)
	numTracker.KnownNumbers = make([]byte, maxCapacity+1)
	numTracker.Stats = &Statistics{}
	return numTracker, nil
}

// Processes a number, validates and passes it on to a channel
// in a pipelined fashion (after converting it to a string)
func (n *NumberTracker) ProcessNumber(ctx context.Context,
	inputStream <-chan int) <-chan string {
	output := make(chan string)
	numSetLength := len(n.KnownNumbers)
	go func() {
		defer close(output)
		for input := range inputStream {
			select {
			case <-ctx.Done():
				return
			default:
				valid := input >= 0 && input < numSetLength
				if valid {
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
