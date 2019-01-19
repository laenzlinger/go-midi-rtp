package rtp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/laenzlinger/go-midi-rtp/timestamp"
)

// Generic RTP constants
const (
	version2Bit  = 0x80
	extensionBit = 0x10
	paddingBit   = 0x20
	markerBit    = 0x80
	ccMask       = 0x0f
	ptMask       = 0x7f
	countMask    = 0x1f
)

// RTP-MIDI constants
const (
	minimumBufferLengt = 12
)

const (
	padding   = 0x00
	extension = 0x00
	ccBits    = 0x00
	firstByte = version2Bit | padding | extension | ccBits
)

const (
	marker      = markerBit
	payloadType = 0x61
	secondByte  = payloadType
)

// MIDIMessage represents a MIDI package exchanged over RTP.
//
// The implementation is tested only with Apple MIDI Network Driver.
//
// see https://en.wikipedia.org/wiki/RTP-MIDI
// see https://developer.apple.com/library/archive/documentation/Audio/Conceptual/MIDINetworkDriverProtocol/MIDI/MIDI.html
// see https://tools.ietf.org/html/rfc6295
/*
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   | V |P|X|  CC   |M|     PT      |        Sequence number        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           Timestamp                           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                             SSRC                              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+


   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                     MIDI command section ...                  |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                       Journal section ...                     |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

*/
type MIDIMessage struct {
	SequenceNumber uint16
	SSRC           uint32
	Commands       MIDICommands
}

// MIDICommands the list of MIDICommand sent inside a MIDIMessage
type MIDICommands struct {
	Timestamp time.Time
	Commands  []MIDICommand
}

// MIDIPayload contains the MIDI payload to be sent.
type MIDIPayload []byte

// MIDICommand represents a single command containing a DeltaTime and the Payload
type MIDICommand struct {
	DeltaTime time.Duration
	Payload   MIDIPayload
}

// Decode a byte buffer into a MIDIMessage
func Decode(buffer []byte) (msg MIDIMessage, err error) {
	msg = MIDIMessage{}
	if len(buffer) < minimumBufferLengt {
		err = fmt.Errorf("buffer is too small: %d bytes", len(buffer))
		return
	}
	// FIXME implement decoder
	return
}

// Encode the MIDIMessage into a byte buffer.
func Encode(m MIDIMessage, start time.Time) []byte {

	b := new(bytes.Buffer)

	b.WriteByte(firstByte)
	b.WriteByte(secondByte)
	binary.Write(b, binary.BigEndian, m.SequenceNumber)
	ts := timestamp.Of(m.Commands.Timestamp, start).Uint32()
	binary.Write(b, binary.BigEndian, uint32(ts))
	binary.Write(b, binary.BigEndian, m.SSRC)

	m.Commands.encode(b, start)

	return b.Bytes()
}

func (m MIDIMessage) String() string {
	return fmt.Sprintf("RM SSRC=0x%x sn=%d", m.SSRC, m.SequenceNumber)
}

/*

0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|B|J|Z|P|LEN... |  MIDI list ...                                |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                  Figure 2 -- MIDI Command Section


+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Delta Time 0     (1-4 octets long, or 0 octets if Z = 0)     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  MIDI Command 0   (1 or more octets long)                     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Delta Time 1     (1-4 octets long)                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  MIDI Command 1   (1 or more octets long)                     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                              ...                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Delta Time N     (1-4 octets long)                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  MIDI Command N   (0 or more octets long)                     |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                Figure 3 -- MIDI List Structure
*/

const (
	emtpyHeader  = byte(0x00)
	bigHeaderBit = 0x80 // Big Header: 2 octets
	journalBit   = 0x40 // Journal persent
	zeroDeltaBit = 0x20 // DeltaTime present for first MIDI command
	phantomBit   = 0x10 // Status byte was not present in original MIDI command
	lenMask      = 0x0f // Mask for the length information
)

func (mcs MIDICommands) encode(w io.Writer, start time.Time) {
	if len(mcs.Commands) == 0 {
		w.Write([]byte{emtpyHeader})
		return
	}
	header := emtpyHeader
	b := new(bytes.Buffer)

	for i, mc := range mcs.Commands {
		if i == 0 && mc.DeltaTime > 0 {
			header = header | zeroDeltaBit
			timestamp.EncodeDeltaTime(mcs.Timestamp, start, mc.DeltaTime, b)
		}
		if i > 0 {
			timestamp.EncodeDeltaTime(mcs.Timestamp, start, mc.DeltaTime, b)
		}
		mc.Payload.encode(b)
	}

	if b.Len() > 4095 {
		// FIXME handle messages with size > 4095 octets (error and crop)
	} else if b.Len() > 15 {
		header = header | bigHeaderBit | (byte(b.Len()>>8) & lenMask)
		count := byte(b.Len())
		w.Write([]byte{header, count})
	} else {
		header = header | (byte(b.Len()) & lenMask)
		w.Write([]byte{header})
	}

	w.Write(b.Bytes())
}

func (p MIDIPayload) encode(w io.Writer) {
	// FIXME maybe this encoding is not correct
	if len(p) == 0 {
		return
	}
	w.Write(p)
}
