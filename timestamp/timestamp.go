package timestamp

import (
	"io"
	"time"
)

const (
	rate = 100 * time.Microsecond
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
/*
   One-Octet Delta Time:

      Encoded form: 0ddddddd
      Decoded form: 00000000 00000000 00000000 0ddddddd

   Two-Octet Delta Time:

      Encoded form: 1ccccccc 0ddddddd
      Decoded form: 00000000 00000000 00cccccc cddddddd

   Three-Octet Delta Time:

      Encoded form: 1bbbbbbb 1ccccccc 0ddddddd
      Decoded form: 00000000 000bbbbb bbcccccc cddddddd

   Four-Octet Delta Time:

      Encoded form: 1aaaaaaa 1bbbbbbb 1ccccccc 0ddddddd
      Decoded form: 0000aaaa aaabbbbb bbcccccc cddddddd

*/
func EncodeDeltaTime(reference time.Time, start time.Time, delta time.Duration, w io.Writer) {

	ticks := Of(reference.Add(delta), start).Uint32() - Of(reference, start).Uint32()
	if ticks >= 0x10000000 {
		// FIXME pass through the error up to the client
		// send the highest possible value
		w.Write([]byte{0xff, 0xff, 0xff, 0x8f})
	} else if ticks >= 0x200000 {
		low := byte(ticks & 0x7f)
		byte2 := byte((ticks >> 7) | 0x80)
		byte3 := byte((ticks >> 14) | 0x80)
		high := byte((ticks >> 21) | 0x80)
		w.Write([]byte{high, byte3, byte2, low})
	} else if ticks >= 0x4000 {
		low := byte(ticks & 0x7f)
		middle := byte((ticks >> 7) | 0x80)
		high := byte((ticks >> 14) | 0x80)
		w.Write([]byte{high, middle, low})
	} else if ticks >= 0x80 {
		low := byte(ticks & 0x7f)
		high := byte((ticks >> 7) | 0x80)
		w.Write([]byte{high, low})
	} else {
		w.Write([]byte{byte(ticks)})
	}

}

// Uint64 returns the long representation of the Timesteamp
func (ts Timestamp) Uint64() uint64 {
	return uint64(ts)
}

// Uint32 returns the short representation of the Timesteamp
func (ts Timestamp) Uint32() uint32 {
	return uint32(ts) & 0xffffffff
}
