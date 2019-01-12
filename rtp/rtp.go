package rtp

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	padding            = 0x00
	extension          = 0x00
	cc                 = 0x00
	marker             = markerBit
	payloadType        = 0x61
)

// MIDIMessage represents a MIDI package exchanged over RTP.
// 
// The implementation is tested only with Apple MIDI Network Driver.
// 
// see https://en.wikipedia.org/wiki/RTP-MIDI
// see https://developer.apple.com/library/archive/documentation/Audio/Conceptual/MIDINetworkDriverProtocol/MIDI/MIDI.html
// see https://tools.ietf.org/html/rfc6295
type MIDIMessage struct {
	SequenceNumber uint32
	SSRC           uint32
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

	// FIXME implement encoder
	binary.Write(b, binary.BigEndian, 0x00)

	return b.Bytes()
}

func (m MIDIMessage) String() string {
	return fmt.Sprintf("sn=%x SSRC=%x", m.SequenceNumber, m.SSRC)
}
