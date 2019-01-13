package session

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/laenzlinger/go-midi-rtp/rtp"
	"github.com/laenzlinger/go-midi-rtp/sip"
)

// MIDINetworkSession can offer or accept streams.
type MIDINetworkSession struct {
	LocalNaame     string
	BonjourName    string
	Port           uint16
	SSRC           uint32
	SequenceNumber uint16
	StartTime      time.Time
	connections    sync.Map
}

// Start is starting a new session
func Start(bonjourName string, port uint16) (s *MIDINetworkSession) {
	session := MIDINetworkSession{
		BonjourName:    bonjourName,
		SSRC:           rand.Uint32(),
		Port:           port,
		StartTime:      time.Now(),
		SequenceNumber: uint16(rand.Int()),
	}

	go messageLoop(port, &session)

	go messageLoop(port+1, &session)

	return &session
}

// End is ending a session
func (s *MIDINetworkSession) End() {
	s.connections.Range(func(k, v interface{}) bool {
		v.(*MIDINetworkConnection).End()
		return true
	})
}

// SendMIDIMessage sends the MIDI payload immediately to all MIDINetworkConnections
func (s *MIDINetworkSession) SendMIDIMessage(payload []byte) {
	s.SequenceNumber++
	m := rtp.MIDIMessage{
		SequenceNumber: s.SequenceNumber,
		SSRC:           s.SSRC,
		Commands: rtp.MIDICommands{
			Timestamp: time.Now(),
			Commands:  []rtp.MIDICommand{{Payload: payload}},
		},
	}
	s.connections.Range(func(k, v interface{}) bool {
		v.(*MIDINetworkConnection).SendMIDIMessage(m)
		return true
	})
}

func messageLoop(port uint16, s *MIDINetworkSession) {
	pc, mcErr := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if mcErr != nil {
		panic(mcErr)
	}
	defer pc.Close()
	buffer := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		msg, err := sip.Decode(buffer[:n])
		if err != nil {
			fmt.Println(err)
			fmt.Println(hex.Dump(buffer[:n]))
			continue
		}
		log.Printf("-> incoming message: %v", msg)

		conn, found := s.getConnection(msg)
		if found {
			conn.handleControl(msg, pc, addr)
		}
	}
}

func (s *MIDINetworkSession) getConnection(msg sip.ControlMessage) (c *MIDINetworkConnection, found bool) {
	if msg.Cmd == sip.Invitation {
		log.Printf("New connection requested from remote participant SSRC [%x]", msg.SSRC)
		conn, found := s.connections.LoadOrStore(msg.SSRC, s.createConnection(msg))
		if found {
			log.Printf("Connections was already established to SSRC [%x]", msg.SSRC)		
		}
		return conn.(*MIDINetworkConnection), true
	}
	conn, found := s.connections.Load(msg.SSRC)
	if !found {
		log.Printf("Connection to SSRC [%x] not found", msg.SSRC)		
		return nil, false
	}
	return conn.(*MIDINetworkConnection), found
}

func (s *MIDINetworkSession) removeConnection(conn *MIDINetworkConnection) {
	log.Printf("Connection ended by remote participant SSRC [%x]", conn.RemoteSSRC)
	s.connections.Delete(conn.RemoteSSRC)
}

func (s *MIDINetworkSession) createConnection(msg sip.ControlMessage) *MIDINetworkConnection {
	host := MIDINetworkHost{BonjourName: msg.Name}
	conn := MIDINetworkConnection{
		Session:    s,
		Host:       host,
		RemoteSSRC: msg.SSRC,
		State:      initial,
	}
	return &conn
}
