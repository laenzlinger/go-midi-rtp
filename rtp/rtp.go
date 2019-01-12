package rtp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
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
	secondByte  = marker | payloadType
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
func Encode(m MIDIMessage) []byte {

	b := new(bytes.Buffer)

	b.WriteByte(firstByte)
	b.WriteByte(secondByte)
	binary.Write(b, binary.BigEndian, m.SequenceNumber)
	// FIXME encode timestamp
	binary.Write(b, binary.BigEndian, uint32(0))
	binary.Write(b, binary.BigEndian, m.SSRC)

	m.Commands.encode(b)

	return b.Bytes()
}

func (m MIDIMessage) String() string {
	return fmt.Sprintf("sn=%x SSRC=%x", m.SequenceNumber, m.SSRC)
}

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |B|J|Z|P|LEN... |  MIDI list ...                                |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
const (
	emtpyHeader   = byte(0x00)
	bigMessageBit = 0x80
	journalBit    = 0x40
	zeroDeltaBit  = 0x20
	phantomBit    = 0x10
	lenMask       = 0x0f
)

func (mcs MIDICommands) encode(w io.Writer) {
	if len(mcs.Commands) == 0 {
		w.Write([]byte{emtpyHeader})
		return
	}
	header := emtpyHeader
	b := new(bytes.Buffer)
	if len(mcs.Commands) == 1 {
		mc := mcs.Commands[0]
		if mc.DeltaTime == 0 && len(mc.Payload) > 0 {
			header = header | zeroDeltaBit
			mc.Payload.encode(b)
		} 

		// FIXME handle message with delta time
		
	}

	// FIXME handle multiple commands

	// FIXME handle messages with size > 15 octets
	header = header | (byte(b.Len()) & lenMask)

	binary.Write(w, binary.BigEndian, header)
	w.Write(b.Bytes())
}

func (p MIDIPayload) encode(w io.Writer) {
	// FIXME maybe this encoding is not correct
	if len(p) == 0 {
		return
	}
	w.Write(p)
}
