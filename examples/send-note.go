package main

import (
	"github.com/laenzlinger/go-midi-rtp/sip"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
)

func main() {
	server, err := zeroconf.Register("GoZeroconf", "_apple-midi._udp", "local.", 6005, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ":6005")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	//simple read
	buffer := make([]byte, 1024)
	n, _, err := pc.ReadFrom(buffer)

	dump := hex.Dump(buffer[:n])
	fmt.Println(dump)
	if err != nil {
		log.Fatal(err)
	}

//	00000000  ff ff 49 4e 00 00 00 02  3d 1b 58 ba 02 f8 f5 8e  |..IN....=.X.....|
//	00000010  55 4d 30 30 38 35 36 00                           |UM00856.|


	cmd, err := sip.Parse(buffer[:n])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("received command: %v", cmd)

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

func handleMessage()  {

}
