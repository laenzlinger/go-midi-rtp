package sip

import (
	"encoding/binary"
	"fmt"
)

type command uint16

const (
	invitation          command = 0x494E
	invitationRejected  command = 0x4E4F
	invitationAccepted  command = 0x4F4B
	end                 command = 0x4259
	synchronization     command = 0x434B
	receiverFeedback    command = 0x5253
	bitrateReceiveLimit command = 0x524C
)

const header = 0xffff

const minimumBufferLengt = 4

// ControlMessage represents the Apple MIDI control messages.ControlMessage
//
// see https://en.wikipedia.org/wiki/RTP-MIDI
type ControlMessage struct {
	cmd command
}

// Parse a buffer into a control message
func Parse(buffer []byte) (m ControlMessage, err error) {
	if len(buffer) < minimumBufferLengt {
		return ControlMessage{}, fmt.Errorf("buffer is too small: %d bytes", len(buffer))
	}

	h := binary.BigEndian.Uint16(buffer[0:2])
	if h != header {
		return ControlMessage{}, fmt.Errorf("invalid header: %x", h)
	}

	cmd := command(binary.BigEndian.Uint16(buffer[2:4]))
	switch cmd {
	case invitation:
	case invitationAccepted:
	case invitationRejected:
	case end:
		/*
			this.version = buffer.readUInt32BE(4);
			this.token = buffer.readUInt32BE(8);
			this.ssrc = buffer.readUInt32BE(12);
			this.name = buffer.toString('utf-8', 16);
			break;
		*/
	case synchronization:
		/*
			this.ssrc = buffer.readUInt32BE(4, 8)
			this.count = buffer.readUInt8(8)
			this.padding = (buffer.readUInt8(9) << 0xF0) + buffer.readUInt16BE(10)
			this.timestamp1 = buffer.slice(12, 20) //[buffer.readUInt32BE(12), buffer.readUInt32BE(16)];
			this.timestamp2 = buffer.slice(20, 28) //[buffer.readUInt32BE(20), buffer.readUInt32BE(24)];
			this.timestamp3 = buffer.slice(28, 36) //[buffer.readUInt32BE(28), buffer.readUInt32BE(32)];
			break
		*/
	case receiverFeedback:
		/*
			this.ssrc = buffer.readUInt32BE(4, 8)
			this.sequenceNumber = buffer.readUInt16BE(8)
			break
		*/
	}

	return ControlMessage{cmd: cmd}, nil
}
