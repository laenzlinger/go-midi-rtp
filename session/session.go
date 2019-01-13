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
	LocalNaame  string
	BonjourName string
	Port        uint16
	SSRC        uint32
	StartTime   time.Time
	connections sync.Map
}

// Start is starting a new session
func Start(bonjourName string, port uint16) (s *MIDINetworkSession) {
	ssrc := rand.Uint32()
	startTime := time.Now()
	session := MIDINetworkSession{
		BonjourName: bonjourName,
		SSRC:        ssrc,
		Port:        port,
		StartTime:   startTime,
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
	m := rtp.MIDIMessage{
		SequenceNumber: 1, // FIXME use random and increase for each message
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
		fmt.Println(hex.Dump(buffer[:n]))
		if err != nil {
			fmt.Println(err)
			continue
		}

		msg, err := sip.Decode(buffer[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Printf("-> incoming message: %v", msg)

		s.getConnection(msg).handleControl(msg, pc, addr)
	}
}

func (s *MIDINetworkSession) getConnection(msg sip.ControlMessage) *MIDINetworkConnection {
	// FIXME optimize to only create a session for IN message
	conn, found := s.connections.LoadOrStore(msg.SSRC, s.createConnection(msg))
	if !found {
		log.Printf("New connection requested from remote participant SSRC [%x]", msg.SSRC)
	}
	return conn.(*MIDINetworkConnection)
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
