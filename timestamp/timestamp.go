package timestamp

import (
	"encoding/binary"
	"io"
	"time"
)

const (
	nanoSecond  = 1
	milliSecond = 1000 * nanoSecond
	rate        = 100 * milliSecond
)

// Timestamp is used for control messages or MIDI messages to define the relative session time.
type Timestamp uint64

// Now returns the Timestamp of now
func Now(start time.Time) Timestamp {
	return Of(time.Now(), start)
}

// Of returns the Timestam of the given time
func Of(t time.Time, start time.Time) Timestamp {
	return Timestamp(t.Sub(start).Nanoseconds() / int64(rate))
}

// EncodeDeltaTime writes the encoded delta time onto the writer
func EncodeDeltaTime(reference time.Time, start time.Time, delta time.Duration, w io.Writer) {

	ticks := Of(reference.Add(delta), start).Uint64() - Of(reference, start).Uint64()
	// FIXME correctly handle higher values
	binary.Write(w, binary.BigEndian, byte(ticks))

}

// Uint64 returns the long representation of the Timesteamp
func (ts Timestamp) Uint64() uint64 {
	return uint64(ts)
}

// Uint32 returns the short representation of the Timesteamp
func (ts Timestamp) Uint32() uint32 {
	return uint32(ts) & 0xffffffff
}
