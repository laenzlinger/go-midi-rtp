package rtp

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_encode_of_message(t *testing.T) {
	// given
	c := MIDICommand{Payload: []byte{0x90, 0x3c, 0x40}}
	mcs := MIDICommands{
		Commands:  []MIDICommand{c},
		Timestamp: time.Now(),
	}
	m := MIDIMessage{
		SequenceNumber: 0xaabb,
		SSRC:           0xccddeeff,
		Commands:       mcs,
	}

	// when
	b := Encode(m)
	// then
	assert.Equal(t, []byte{
		0x80, 0xe1, 0xaa, 0xbb, // Header | Sequence Number
		0x00, 0x00, 0x00, 0x00, // Timestamp
		0xcc, 0xdd, 0xee, 0xff, // SRCC
		0x23, 0x90, 0x3c, 0x40, // MIDI Commands
	}, b)
}

func Test_encode_of_empty_commands(t *testing.T) {
	// given
	m := MIDICommands{}
	b := new(bytes.Buffer)
	// when
	m.encode(b)
	/* then

	           0                   1                   2                   3
	       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |B|J|Z|P|LEN... |  MIDI list ...                                |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |0|0|0|0|0 0 0 0|                                               |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	*/
	assert.Equal(t, []byte{0x00}, b.Bytes())
}

func Test_encode_of_empty_command(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	c := MIDICommand{}
	mcs := MIDICommands{Commands: []MIDICommand{c}}
	// when
	mcs.encode(b)
	//then
	assert.Equal(t, []byte{0x00}, b.Bytes())
}

func Test_encode_of_single_command_without_delta(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	c := MIDICommand{Payload: []byte{0x90, 0x3c, 0x40}}
	mcs := MIDICommands{Commands: []MIDICommand{c}}
	// when
	mcs.encode(b)
	/* then

	           0                   1                   2                   3
	       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |B|J|Z|P|LEN... |  MIDI list ...                                |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |0|0|1|0|0 0 1 1|      0x90          0x3c          0x40         |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	*/
	assert.Equal(t, []byte{0x23, 0x90, 0x3c, 0x40}, b.Bytes())
}
