package session

import (
	"fmt"
	"github.com/laenzlinger/go-midi-rtp/sip"
	"net"
)

type state uint8

const (
	initial state = iota
	controlChannelEstablished
	ready
)

// MidiNetworkHost represents information about the remote
type MidiNetworkHost struct {
	ControlPort net.Addr
	MidiPort net.Addr
	// MDNSName used to advertise expose the remote session with multicast DNS.
	MDNSName string
}

// MidiNetworkConnection specifies a connection to a MIDI network host.
type MidiNetworkConnection struct {
	Host MidiNetworkHost
	State state
}

// Create a new connection
func create(msg sip.ControlMessage) *MidiNetworkConnection {
	conn := MidiNetworkConnection{
		Host : MidiNetworkHost{MDNSName: msg.Name},
		State: initial,
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
		// send NO to control port
	}
}


func (conn MidiNetworkConnection) sendInvitationAccepted(msg sip.ControlMessage, addr net.Addr, pc net.PacketConn) {

	accept := sip.ControlMessage{
		Cmd: sip.InvitationAccepted,
		Version: 2,
		Token: msg.Token,
		SSRC: msg.SSRC, // FIXME use own session token
		Name: "GoZeroconf", // FIXME session name
	}

	buf := sip.Marshall(accept)


	_, err := pc.WriteTo(buf, addr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("sent: %v\n", accept)
};
