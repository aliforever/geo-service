package main

import "time"

type Statistics struct {
	TimeElapsed      time.Duration
	AcceptedEntries  int
	DiscardedEntries int
}
