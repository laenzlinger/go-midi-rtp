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
				Commands:  []rtp.MIDICommand{
					{Payload: []byte{0x96, 0x3c, 0x7f}},
					{Payload: []byte{0x86, 0x3c, 0x00}, DeltaTime: 10*time.Millisecond},
				},
			}

			s.SendMIDICommands(mcs)
		}
	}

	log.Println("Shutting down.")
	s.End()

}
