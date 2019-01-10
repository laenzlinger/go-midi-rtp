package session

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/laenzlinger/go-midi-rtp/sip"
)

type state uint8

const (
	initial state = iota
	controlChannelEstablished
	ready
)

// MidiNetworkHost represents information about the remote
type MidiNetworkHost struct {
	// ControlPort is used to exchange session control messages (IN, OK, NO, BY...)
	ControlAddr net.Addr
	ControlPc   net.PacketConn
	// MidiPort is ised to exchange MIDI payload and synchronisation
	MidiAddr net.Addr
	MidiPc   net.PacketConn
	// MDNSName used to advertise expose the remote session with multicast DNS.
	BonjourName string
}

// MidiNetworkConnection specifies a connection to a MIDI network host.
type MidiNetworkConnection struct {
	Session        *MidiNetworkSession
	Host           MidiNetworkHost
	RemoteSSRC     uint32
	State          state
}

// HandleControl a sipControlMessage
func (conn *MidiNetworkConnection) HandleControl(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	switch msg.Cmd {
	case sip.Invitation:
		conn.handleInvitation(msg, pc, addr)
	case sip.End:
		conn.handleEnd()
	case sip.Synchronization:
		conn.handleSynchonization(msg, pc, addr)
	}
}

// End the session
func (conn *MidiNetworkConnection) End() {
	// FIXME what to do now?
	log.Println("Ending connedtion")
	conn.sendConnectionEnd(conn.Host.ControlAddr, conn.Host.ControlPc)
}

func (conn *MidiNetworkConnection) handleInvitation(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	switch conn.State {
	case initial:
		conn.Host.ControlAddr = addr
		conn.Host.ControlPc = pc
		conn.sendInvitationAccepted(msg, addr, pc)
		conn.State = controlChannelEstablished
	case controlChannelEstablished:
		conn.Host.MidiAddr = addr
		conn.Host.MidiPc = pc
		conn.sendInvitationAccepted(msg, addr, pc)
		conn.State = ready
	case ready:
		// FIXME send NO to control port
	}
}

func (conn *MidiNetworkConnection) handleEnd() {
	conn.Session.removeConnection(conn)
}

func (conn *MidiNetworkConnection) sendConnectionEnd(addr net.Addr, pc net.PacketConn) {

	end := sip.ControlMessage{
		Cmd:  sip.End,
		SSRC: conn.Session.SSRC,
	}

	conn.sendMessage(end, addr, pc)
}

func (conn *MidiNetworkConnection) sendInvitationAccepted(msg sip.ControlMessage, addr net.Addr, pc net.PacketConn) {

	accept := sip.ControlMessage{
		Cmd:   sip.InvitationAccepted,
		Token: msg.Token,
		SSRC:  conn.Session.SSRC,
		Name:  conn.Session.BonjourName,
	}

	conn.sendMessage(accept, addr, pc)
}

func (conn *MidiNetworkConnection) handleSynchonization(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	if conn.State == ready {
		switch len(msg.Timestamps) {
		case 1:
			fallthrough
		case 2:
			ts := uint64(time.Since(conn.Session.StartTime).Nanoseconds() / int64(100000))
			newTs := append(msg.Timestamps, ts)

			sync := sip.ControlMessage{
				Cmd:        sip.Synchronization,
				SSRC:       conn.Session.SSRC,
				Timestamps: newTs,
			}
			conn.sendMessage(sync, addr, pc)
		case 3:
			// FIXME calculate offset_estimate = ((timestamp3 + timestamp1) / 2) - timestamp2
		}
	}
}

func (conn *MidiNetworkConnection) sendMessage(msg sip.ControlMessage, addr net.Addr, pc net.PacketConn) {

	_, err := pc.WriteTo(sip.Encode(msg), addr)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("<- outgoing message: %v", msg)
}
