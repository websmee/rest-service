package local

import "time"

type Object struct {
	ID       int64
	LastSeen time.Time
}
