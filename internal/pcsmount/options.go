package pcsmount

import "time"

// Options controls a mount session.
type Options struct {
	RemoteRoot   string
	CacheTTL     time.Duration
	Debug        bool
	SingleThread bool
	FuseOptions  []string
}
