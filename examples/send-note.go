package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/laenzlinger/go-midi-rtp/session"

	"github.com/grandcat/zeroconf"
)

func main() {
	port := 6005
	bonjourName := "GoZeroconf"
	server, err := zeroconf.Register(bonjourName, "_apple-midi._udp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	session.Start(bonjourName, port)

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// Exit by user
	case <-time.After(time.Second * 120):
		// Exit by timeout
	}

	log.Println("Shutting down.")

}
