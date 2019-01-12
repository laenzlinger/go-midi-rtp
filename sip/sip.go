package sip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// Command defines one of the commands defined by the apple SIP
type Command uint16

const (
	// Invitation message is sent to invite a remote participant to the session.
	Invitation Command = 0x494E
	// InvitationRejected is sent to reject the invitation.
	InvitationRejected Command = 0x4E4F
	// InvitationAccepted is sent to accept the invitation.
	InvitationAccepted Command = 0x4F4B
	// End Message is sent to end the current session.
	End Command = 0x4259
	// Synchronization message is sent to synchronize the timestamps between participants.
	Synchronization Command = 0x434B
	// ReceiverFeedback is sent to update the journal on the remote participant.
	ReceiverFeedback Command = 0x5253
	// BitrateReceiveLimit is currently unused but defined in Wireshark
	// see https://github.com/boundary/wireshark/blob/master/epan/dissectors/packet-applemidi.c
	BitrateReceiveLimit Command = 0x524C
)

const (
	header          = uint16(0xffff)
	protocolVersion = uint32(2)
)

const minimumBufferLengt = 4

// ControlMessage represents the Apple MIDI ControlMessage
//
// see https://en.wikipedia.org/wiki/RTP-MIDI
// see https://developer.apple.com/library/archive/documentation/Audio/Conceptual/MIDINetworkDriverProtocol/MIDI/MIDI.html
type ControlMessage struct {
	Cmd        Command
	Token      uint32
	SSRC       uint32
	Name       string
	Timestamps []uint64
}

// Decode a byte buffer into a ControlMessage
func Decode(buffer []byte) (msg ControlMessage, err error) {
	msg = ControlMessage{}
	if len(buffer) < minimumBufferLengt {
		err = fmt.Errorf("buffer is too small: %d bytes", len(buffer))
		return
	}

	h := binary.BigEndian.Uint16(buffer[0:2])
	if h != header {
		err = fmt.Errorf("invalid header: %x", h)
		return
	}
	msg.Cmd = Command(binary.BigEndian.Uint16(buffer[2:4]))
	switch msg.Cmd {
	case Invitation:
		fallthrough
	case InvitationAccepted:
		fallthrough
	case InvitationRejected:
		fallthrough
	case End:
		version := binary.BigEndian.Uint32(buffer[4:8])
		if version != protocolVersion {
			fmt.Println("Warning: Unsupported protocol version: ", version)
		}
		msg.Token = binary.BigEndian.Uint32(buffer[8:12])
		msg.SSRC = binary.BigEndian.Uint32(buffer[12:16])
		if msg.Cmd != End {
			msg.Name = strings.TrimRight(string(buffer[16:]), "\x00")
		}
	case Synchronization:
		msg.SSRC = binary.BigEndian.Uint32(buffer[4:8])
		count := buffer[8] + 1
		for i := byte(0); i < count; i++ {
			ts := binary.BigEndian.Uint64(buffer[12+i*8 : 20+i*8])
			msg.Timestamps = append(msg.Timestamps, ts)
		}
	case ReceiverFeedback:
		/*
			this.ssrc = buffer.readUInt32BE(4, 8)
			this.sequenceNumber = buffer.readUInt16BE(8)
			break
		*/
	}

	return
}

// Encode the ControlMessage into a byte buffer.
func Encode(m ControlMessage) (buf []byte, err error) {
	b := new(bytes.Buffer)

	binary.Write(b, binary.BigEndian, header)
	binary.Write(b, binary.BigEndian, m.Cmd)

	switch m.Cmd {
	case Invitation:
		fallthrough
	case InvitationAccepted:
		fallthrough
	case InvitationRejected:
		fallthrough
	case End:
		binary.Write(b, binary.BigEndian, protocolVersion)
		binary.Write(b, binary.BigEndian, m.Token)
		binary.Write(b, binary.BigEndian, m.SSRC)
		if m.Cmd != End {
			b.WriteString(m.Name)
			b.WriteByte(0)
		}

	case Synchronization:
		if len(m.Timestamps) < 1 {
			return []byte{}, fmt.Errorf("At least 1 timestamp is expected")
		}
		binary.Write(b, binary.BigEndian, m.SSRC)
		binary.Write(b, binary.BigEndian, byte(len(m.Timestamps)-1))
		binary.Write(b, binary.BigEndian, byte(0x00))
		binary.Write(b, binary.BigEndian, uint16(0x0000))
		for i := 0; i < 3; i++ {
			var ts uint64
			if i < len(m.Timestamps) {
				ts = m.Timestamps[i]
			} else {
                ts = 0
			}
			binary.Write(b, binary.BigEndian, ts)
		}
	case ReceiverFeedback:
		/*
			buffer = new Buffer(12);
			buffer.writeUInt16BE(this.start, 0);
			buffer.writeUInt16BE(commandByte, 2);
			buffer.writeUInt32BE(this.ssrc, 4);
			buffer.writeUInt16BE(this.sequenceNumber, 8);

		*/
	}

	return b.Bytes(), nil
}

func (c Command) String() string {
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(c))
	return string(buffer)
}

func (m ControlMessage) String() string {
	if m.Cmd == Synchronization {
		res := fmt.Sprintf("%v SSRC=%x", m.Cmd, m.SSRC)
		for i, ts := range m.Timestamps {
			res = fmt.Sprintf("%v ts%d=%d", res, i, ts)
		}
		return res
	}
	return fmt.Sprintf("%v token=%x SSRC=%x name=[%v]", m.Cmd, m.Token, m.SSRC, m.Name)
}
