package session

import (
	"fmt"
	"log"
	"net"

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
	ControlPort net.Addr
	// MidiPort is ised to exchange MIDI payload and synchronisation
	MidiPort    net.Addr
	// MDNSName used to advertise expose the remote session with multicast DNS.
	BonjourName string
}

// MidiNetworkConnection specifies a connection to a MIDI network host.
type MidiNetworkConnection struct {
	Session *MidiNetworkSession
	Host    MidiNetworkHost
	State   state
}

// Create a new connection
func create(msg sip.ControlMessage, session *MidiNetworkSession) *MidiNetworkConnection {
	host := MidiNetworkHost{BonjourName: msg.Name}
	conn := MidiNetworkConnection{
		Session: session,
		Host:    host,
		State:   initial,
	}
	return &conn
}

// HandleControl a sipControlMessage
func (conn *MidiNetworkConnection) HandleControl(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	switch msg.Cmd {
	case sip.Invitation:
		conn.handleInvitation(msg, pc, addr)
	}
}

func (conn MidiNetworkConnection) handleInvitation(msg sip.ControlMessage, pc net.PacketConn, addr net.Addr) {
	switch conn.State {
	case initial:
		conn.Host.ControlPort = addr
		conn.sendInvitationAccepted(msg, addr, pc)
		conn.State = controlChannelEstablished
	case controlChannelEstablished:
		conn.Host.MidiPort = addr
		conn.sendInvitationAccepted(msg, addr, pc)
		conn.State = ready
	case ready:
		// FIXME send NO to control port
	}
}

func (conn MidiNetworkConnection) sendInvitationAccepted(msg sip.ControlMessage, addr net.Addr, pc net.PacketConn) {

	accept := sip.ControlMessage{
		Cmd:   sip.InvitationAccepted,
		Token: msg.Token,
		SSRC:  msg.SSRC, // FIXME use own session token
		Name:  conn.Session.BonjourName,
	}

	_, err := pc.WriteTo(sip.Encode(accept), addr)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("<- outgoing message: %v", accept)
}
