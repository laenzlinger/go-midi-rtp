package rtp

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_encode_of_message(t *testing.T) {
	// given
	start := time.Now()
	c := MIDICommand{Payload: []byte{0x90, 0x3c, 0x40}}
	mcs := MIDICommands{
		Commands:  []MIDICommand{c},
		Timestamp: start.Add(100 * time.Microsecond),
	}
	m := MIDIMessage{
		SequenceNumber: 0xaabb,
		SSRC:           0xccddeeff,
		Commands:       mcs,
	}

	// when
	b := Encode(m, start)
	// then
	assert.Equal(t, []byte{
		0x80, 0x61, 0xaa, 0xbb, // Header | Sequence Number
		0x00, 0x00, 0x00, 0x01, // Timestamp
		0xcc, 0xdd, 0xee, 0xff, // SRCC
		0x03, 0x90, 0x3c, 0x40, // MIDI Commands
	}, b)
}

func Test_encode_of_empty_commands(t *testing.T) {
	// given
	m := MIDICommands{}
	b := new(bytes.Buffer)
	// when
	m.encode(b, time.Now())
	/* then

	           0                   1                   2                   3
	       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |B|J|Z|P|LEN... |  MIDI list ...                                |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |0|0|0|0|0 0 0 0|
		  +-+-+-+-+-+-+-+-+

	*/
	assert.Equal(t, []byte{0x00}, b.Bytes())
}

func Test_encode_of_empty_command(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	c := MIDICommand{}
	mcs := MIDICommands{Commands: []MIDICommand{c}}
	// when
	mcs.encode(b, time.Now())
	//then
	assert.Equal(t, []byte{0x00}, b.Bytes())
}

func Test_encode_of_single_command_without_delta(t *testing.T) {
	// given
	b := new(bytes.Buffer)
	c := MIDICommand{Payload: []byte{0x90, 0x3c, 0x40}}
	mcs := MIDICommands{Commands: []MIDICommand{c}}
	// when
	mcs.encode(b, time.Now())
	/* then

	           0                   1                   2                   3
	       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |B|J|Z|P|LEN... |  MIDI list ...                                |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |0|0|0|0|0 0 1 1|      0x90          0x3c          0x40         |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	*/
	assert.Equal(t, []byte{
		0x03,             // Header
		0x90, 0x3c, 0x40, // Midi command
	}, b.Bytes())
}

func Test_encode_of_single_command_with_delta(t *testing.T) {
	// given
	now := time.Now()
	b := new(bytes.Buffer)
	c := MIDICommand{
		Payload:   []byte{0x90, 0x3c, 0x40},
		DeltaTime: 10 * time.Millisecond,
	}
	mcs := MIDICommands{
		Commands:  []MIDICommand{c},
		Timestamp: now,
	}
	// when
	mcs.encode(b, now)
	// then
	assert.Equal(t, []byte{
		0x24,             // Header
		0x64,             // Single Byte Delta Time
		0x90, 0x3c, 0x40, // MIDI command
	}, b.Bytes())
}

func Test_encode_of_mulitple_commands(t *testing.T) {
	// given
	now := time.Now()
	b := new(bytes.Buffer)
	mcs := MIDICommands{
		Commands: []MIDICommand{
			{Payload: []byte{0x90, 0x3c, 0x40}},
			{Payload: []byte{0x80, 0x3c, 0x00}, DeltaTime: time.Second},
			{Payload: []byte{0x90, 0x3e, 0x40}},
			{Payload: []byte{0x80, 0x3e, 0x00}, DeltaTime: time.Second},
		},
		Timestamp: now,
	}
	// when
	mcs.encode(b, now)
	/* then
	           0                   1                   2                   3
	       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |B|J|Z|P|LEN... |   LEN (Low)   |    MIDI list ...              |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	      |1|0|0|0|0 0 0 0|     0x11      |     0x90          0x3c        |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	*/
	assert.Equal(t, []byte{
		0x80, 0x11, // Header
		0x90, 0x3c, 0x40, // MIDI command (note on)
		0xce, 0x10, // Delta time (100 ticks)
		0x80, 0x3c, 0x00, // MIDI command (note off)
		0x00,             // Delta time (0 ticks)
		0x90, 0x3e, 0x40, // MIDI command (note on)
		0xce, 0x10, // Delta time (100 ticks)
		0x80, 0x3e, 0x00, // MIDI command (note off)
	}, b.Bytes())
}
