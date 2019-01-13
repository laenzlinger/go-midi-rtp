package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		s.End()
	case <-msg:
		s.SendMIDIMessage([]byte{0x90, 0x3c, 0x40})
		s.End()
	}

	log.Println("Shutting down.")

}
