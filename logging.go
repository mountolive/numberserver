package main

import (
	"context"
	"errors"
)

const LOG_FILE = "./numbers.log"

var BadMaxCapacity = errors.New("Max Capacity for Number Tracker can't be negative")

// Keeps a set of processed numbers
// and statistics book
type NumberTracker struct {
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
	go func() {
		defer close(output)
		for input := range inputStream {
			select {
			case <-ctx.Done():
				return
			default:
				valid := input >= 0 && input < len(n.KnownNumbers)
				if valid && n.KnownNumbers[input] == 0 {
					// Marking it as seen
					n.KnownNumbers[input] = 1
					// passing it on
					output <- input
				}
			}
		}
	}()
	return output
}
