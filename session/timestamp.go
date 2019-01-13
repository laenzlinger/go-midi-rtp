package session

import "time"

// Timestamp is used for control messages or MIDI messages to define the relative session time.
type Timestamp uint64

func (s *MIDINetworkSession) now() Timestamp {
	return Timestamp(time.Since(s.StartTime).Nanoseconds() / int64(100000))
}
// Uint64 returns the long representation of the Timesteamp
func (ts Timestamp) Uint64() uint64 {
	return uint64(ts)
}
