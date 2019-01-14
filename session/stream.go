package session

import (
	"fmt"
	"log"
	"net"

	"github.com/laenzlinger/go-midi-rtp/rtp"
	"github.com/laenzlinger/go-midi-rtp/sip"
	"github.com/laenzlinger/go-midi-rtp/timestamp"
)

type state uint8

const (
	initial state = iota
	controlChannelEstablished
	ready
)

// MIDINetworkHost represents information about the remote
type MIDINetworkHost struct {
	// ControlPort is used to exchange session control messages (IN, OK, NO, BY...)
	ControlAddr net.Addr
	ControlPc   net.PacketConn
	// MIDIPort is ised to exchange MIDI payload and synchronisation
	MIDIAddr net.Addr
	MIDIPc   net.PacketConn
	// MDNSName used to advertise expose the remote session with multicast DNS.
	BonjourName string
}

// MIDINetworkStream specifies a connection to a MIDI network host.
type MIDINetworkStream struct {
	Session    *MIDINetworkSession
	Host       MIDINetworkHost
	RemoteSSRC uint32
	State      state
}

// End the session
func (conn *MIDINetworkStream) End() {
	log.Println("Ending connedtion")
	conn.sendConnectionEnd(conn.Host.ControlAddr, conn.Host.ControlPc)
}

// SendMIDIMessage sends to given MIDIMessage over the RTP-MIDI data port.
func (conn *MIDINetworkStream) SendMIDIMessage(msg rtp.MIDIMessage) {
	buff := rtp.Encode(msg, conn.Session.StartTime)

	_, err := conn.Host.MIDIPc.WriteTo(buff, conn.Host.MIDIAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("<- outgoing payload: %v", msg)
}

// HandleControl a sipControlMessage
func (conn *MIDINetworkStream) handleControl(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	switch msg.Cmd {
	case sip.Invitation:
		conn.handleInvitation(msg, pc, addr)
	case sip.End:
		conn.handleEnd()
	case sip.Synchronization:
		conn.handleSynchonization(msg, pc, addr)
	}
}

func (conn *MIDINetworkStream) handleInvitation(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	switch conn.State {
	case initial:
		conn.Host.ControlAddr = addr
		conn.Host.ControlPc = pc
		conn.sendInvitationAccepted(msg, addr, pc)
		conn.State = controlChannelEstablished
	case controlChannelEstablished:
		conn.Host.MIDIAddr = addr
		conn.Host.MIDIPc = pc
		conn.sendInvitationAccepted(msg, addr, pc)
		conn.State = ready
	case ready:
		// FIXME send NO to control port
	}
}

func (conn *MIDINetworkStream) handleEnd() {
	conn.Session.removeConnection(conn)
}

func (conn *MIDINetworkStream) sendConnectionEnd(addr net.Addr, pc net.PacketConn) {

	end := sip.ControlMessage{
		Cmd:  sip.End,
		SSRC: conn.Session.SSRC,
	}

	conn.sendControlMessage(end, addr, pc)
}

func (conn *MIDINetworkStream) sendInvitationAccepted(msg sip.ControlMessage, addr net.Addr, pc net.PacketConn) {

	accept := sip.ControlMessage{
		Cmd:   sip.InvitationAccepted,
		Token: msg.Token,
		SSRC:  conn.Session.SSRC,
		Name:  conn.Session.BonjourName,
	}

	conn.sendControlMessage(accept, addr, pc)
}

func (conn *MIDINetworkStream) handleSynchonization(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	if conn.State == ready {
		switch len(msg.Timestamps) {
		case 1:
			fallthrough
		case 2:
			ts := timestamp.Now(conn.Session.StartTime).Uint64()
			newTs := append(msg.Timestamps, ts)

			sync := sip.ControlMessage{
				Cmd:        sip.Synchronization,
				SSRC:       conn.Session.SSRC,
				Timestamps: newTs,
			}
			conn.sendControlMessage(sync, addr, pc)
		case 3:
			// FIXME calculate offset_estimate = ((timestamp3 + timestamp1) / 2) - timestamp2
		}
	}
}

func (conn *MIDINetworkStream) sendControlMessage(msg sip.ControlMessage, addr net.Addr, pc net.PacketConn) {
	buff, err := sip.Encode(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = pc.WriteTo(buff, addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("<- outgoing message: %v", msg)
}
