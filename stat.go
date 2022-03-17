package geoservice

import "time"

type Statistics struct {
	Elapsed          time.Duration
	ElapsedParsed    time.Duration
	ElapsedAppend    time.Duration
	Duplicates       int
	AcceptedEntries  int
	DiscardedEntries int
}
