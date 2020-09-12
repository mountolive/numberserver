package main

import (
	"context"
	"errors"
	"sync"
)

const LOG_FILE = "./numbers.log"

var BadMaxCapacity = errors.New("Max Capacity for Number Tracker can't be negative")

// Keeps a set of processed numbers
// and statistics book
type NumberTracker struct {
	sync.RWMutex
	KnownNumbers []int
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
	numTracker.KnownNumbers = make([]int, maxCapacity)
	numTracker.Stats = &Statistics{}
	return numTracker, nil
}

// Processes a number and passes it on to a channel
// in a pipelined fashion
func (n *NumberTracker) ProcessNumber(ctx context.Context,
	inputStream <-chan int) <-chan int {
	output := make(chan int)
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
						output <- input
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
