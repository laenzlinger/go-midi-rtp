package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laenzlinger/go-midi-rtp/rtp"
	"github.com/laenzlinger/go-midi-rtp/session"

	"github.com/grandcat/zeroconf"
)

func main() {
	port := 6005
	bonjourName := "send-note"
	server, err := zeroconf.Register(bonjourName, "_apple-midi._udp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	s := session.Start(bonjourName, uint16(port))

	msg := make(chan rune, 1)
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println(err)
			}
			msg <- char
		}

	}()

	run := true
	for run {
		select {
		case <-sig:
			run = false
		case <-msg:
			mcs := rtp.MIDICommands{
				Timestamp: time.Now(),
				// Play 4 notes of a C-Major chord
				Commands: []rtp.MIDICommand{
					// On the first command, Delta Time should not be needed but it seems
					// that the Apple Midi Network Driver ignores all delta times
					// in case the Z-flag is not set.
					{Payload: []byte{0x96, 0x3c, 0x4f}, DeltaTime: time.Millisecond},
					{Payload: []byte{0x86, 0x3c, 0x00}, DeltaTime: 500 * time.Millisecond},
					{Payload: []byte{0x96, 0x40, 0x5f}},
					{Payload: []byte{0x86, 0x40, 0x00}, DeltaTime: 500 * time.Millisecond},
					{Payload: []byte{0x96, 0x43, 0x6f}},
					{Payload: []byte{0x86, 0x43, 0x00}, DeltaTime: 500 * time.Millisecond},
					{Payload: []byte{0x96, 0x48, 0x7f}},
					{Payload: []byte{0x86, 0x48, 0x00}, DeltaTime: time.Second},
				},
			}

			s.SendMIDICommands(mcs)
		}
	}

	log.Println("Shutting down.")
	s.End()

}
