package sip

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	End                 Command = 0x4259
	// Synchronization message is sent to synchronize the timestamps between participants.
	Synchronization     Command = 0x434B
	// ReceiverFeedback is sent to update the journal on the remote participant.
	ReceiverFeedback    Command = 0x5253
	// BitrateReceiveLimit is currently unused but defined in Wireshark
	// see https://github.com/boundary/wireshark/blob/master/epan/dissectors/packet-applemidi.c
	BitrateReceiveLimit Command = 0x524C
)

const (
	header          = uint16(0xffff)
	protocolVersion = uint32(2)
)

const minimumBufferLengt = 4

// ControlMessage represents the Apple MIDI control messages.ControlMessage
//
// see https://en.wikipedia.org/wiki/RTP-MIDI
// see https://developer.apple.com/library/archive/documentation/Audio/Conceptual/MIDINetworkDriverProtocol/MIDI/MIDI.html
type ControlMessage struct {
	Cmd     Command
	Token   uint32
	SSRC    uint32
	Name    string
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
	
	version := binary.BigEndian.Uint32(buffer[4:8])
	if (version != protocolVersion) {
		fmt.Println("Warning: Unsupported protocol version: ", version)
	}


	cmd := Command(binary.BigEndian.Uint16(buffer[2:4]))
	message := ControlMessage{Cmd: cmd}
	switch cmd {
	case Invitation:
		fallthrough
	case InvitationAccepted:
		fallthrough
	case InvitationRejected:
		fallthrough
	case End:
		message.Token = binary.BigEndian.Uint32(buffer[8:12])
		message.SSRC = binary.BigEndian.Uint32(buffer[12:16])
		message.Name = string(buffer[16:])
		/*
			this.version = buffer.readUInt32BE(4);
			this.token = buffer.readUInt32BE(8);
			this.ssrc = buffer.readUInt32BE(12);
			this.name = buffer.toString('utf-8', 16);
			break;
		*/
	case Synchronization:
		/*
			this.ssrc = buffer.readUInt32BE(4, 8)
			this.count = buffer.readUInt8(8)
			this.padding = (buffer.readUInt8(9) << 0xF0) + buffer.readUInt16BE(10)
			this.timestamp1 = buffer.slice(12, 20) //[buffer.readUInt32BE(12), buffer.readUInt32BE(16)];
			this.timestamp2 = buffer.slice(20, 28) //[buffer.readUInt32BE(20), buffer.readUInt32BE(24)];
			this.timestamp3 = buffer.slice(28, 36) //[buffer.readUInt32BE(28), buffer.readUInt32BE(32)];
			break
		*/
	case ReceiverFeedback:
		/*
			this.ssrc = buffer.readUInt32BE(4, 8)
			this.sequenceNumber = buffer.readUInt16BE(8)
			break
		*/
	}

	return message, nil
}

// Marshall the ControlMessage into a byte buffer.
func Marshall(m ControlMessage) []byte {
	b := new(bytes.Buffer)

	switch m.Cmd {
	case Invitation:
		fallthrough
	case InvitationAccepted:
		fallthrough
	case InvitationRejected:
		fallthrough
	case End:
		binary.Write(b, binary.BigEndian, header)
		binary.Write(b, binary.BigEndian, m.Cmd)
		binary.Write(b, binary.BigEndian, protocolVersion)
		binary.Write(b, binary.BigEndian, m.Token)
		binary.Write(b, binary.BigEndian, m.SSRC)
		b.WriteString(m.Name)
		if m.Cmd != End {
			b.WriteByte(0)
		}

	case Synchronization:
		/*
			buffer = new Buffer(36);
			buffer.writeUInt16BE(this.start, 0);
			buffer.writeUInt16BE(commandByte, 2);
			buffer.writeUInt32BE(this.ssrc, 4);
			buffer.writeUInt8(this.count, 8);
			buffer.writeUInt8(this.padding >>> 0xF0, 9);
			buffer.writeUInt16BE(this.padding & 0x00FFFF, 10);

			this.timestamp1.copy(buffer, 12);
			this.timestamp2.copy(buffer, 20);
			this.timestamp3.copy(buffer, 28);
		*/

	case ReceiverFeedback:
		/*
			buffer = new Buffer(12);
			buffer.writeUInt16BE(this.start, 0);
			buffer.writeUInt16BE(commandByte, 2);
			buffer.writeUInt32BE(this.ssrc, 4);
			buffer.writeUInt16BE(this.sequenceNumber, 8);

		*/
	default:
		//		assert.fail('Not a valid command: "' + this.command + '"');

	}

	return b.Bytes()
}

func (c Command) String() string {
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(c))
	return string(buffer)
}

func (m ControlMessage) String() string {
	return fmt.Sprintf("%v name=%v token=%x SSRC=%x", m.Cmd, m.Name, m.Token, m.SSRC)
}
