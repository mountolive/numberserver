package main

import (
	"fmt"
	"sync"
)

// Bookkeeping struct for input count
type Statistics struct {
	sync.Mutex
	Received   int
	Duplicates int
	Total      int
}

// Prints to STDOUT the current statistics of the server,
// regarding received numbers, number of duplicates and
// total number of unique numbers received (and logged) by the server.
// Resets count of Received and Duplicates after reporting
func (s *Statistics) PrintCurrent() {
	s.Lock()
	defer s.Unlock()
	fmt.Printf("Received %d unique numbers, %d duplicates. "+
		"Unique totals: %d \n", s.Received, s.Duplicates, s.Total)
	s.Received = 0
	s.Duplicates = 0
}

// Increases sessions' duplicate count by 1
func (s *Statistics) IncreaseDups() {
	s.Lock()
	s.Duplicates += 1
	s.Unlock()
}

// Increases session's unique received count by 1
func (s *Statistics) IncreaseReceived() {
	s.Lock()
	defer s.Unlock()
	s.Received += 1
	s.Total += 1
}
