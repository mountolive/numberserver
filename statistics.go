package main

import (
	"fmt"
	"sync"
)

type Statistics struct {
	sync.Mutex
	Received   int
	Duplicates int
	Total      int
}

// Prints to STDOUT the current statistics of the server,
// regarding received numbers, number of duplicates and
// total number of unique numbers received by the server
func (s *Statistics) PrintCurrent() {
	s.Lock()
	fmt.Printf("Received: %d unique numbers, %d duplicates, "+
		"Unique totals: %d \n", s.Received, s.Duplicates, s.Total)
	s.Unlock()
}

// Updates Total statistics based on the values recv and dups,
// Received and Duplicates. If one of them is negative, it
// silently exits
func (s *Statistics) BulkUpdate(recv, dups int) {
	if recv < 0 || dups < 0 {
		return
	}
	if dups > recv {
		return
	}
	s.Lock()
	defer s.Unlock()
	s.Received = recv
	s.Duplicates = dups
	s.Total = s.Total + (recv - dups)
}

// Increases global duplicate count by 1
func (s *Statistics) IncreaseDups() {
	// TODO Implement
}

// Increases global duplicate count by 1
func (s *Statistics) IncreaseReceived() {
	// TODO Implement
}
