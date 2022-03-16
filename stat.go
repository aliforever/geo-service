package geoservice

import "time"

type Statistics struct {
	TimeElapsed      time.Duration
	AcceptedEntries  int
	DiscardedEntries int
}
