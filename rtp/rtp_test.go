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
		Timestamp: start.Add(100*time.Microsecond),
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

/*
note on
00000000  80 61 98 78 2f 0f 1d c5  7e 33 a0 24 43 90 3c 40  |.a.x/...~3.$C.<@|
00000010  20 98 72 00 06 08 00 77  08                       | .r....w.|

note off
00000000  80 61 98 79 2f 0f 24 92  7e 33 a0 24 43 80 3c 00  |.a.y/.$.~3.$C.<.|
00000010  20 98 72 00 07 08 81 f1  3c 40                    | .r.....<@|
*/

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
	      |0|0|0|0|0 0 1 1|      0x90          0x3c          0x40         |
		  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	*/
	assert.Equal(t, []byte{0x03, 0x90, 0x3c, 0x40}, b.Bytes())
}
