package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	iconPath string
	quit chan bool
)

func startListening() {
	if err := initializeActions(); err != nil {
		fmt.Printf("Could not start listening: %s", err.Error())
		os.Exit(-1)
		return
	}

	quit = make(chan bool, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<- signals
		
		fmt.Println("Picked up close signal")
		Quit()
	}()

	go listen()

	// Wait for close signal
	fmt.Println("Waiting for close signal")
	<- quit
}

func listen() {
	var err error
	var packet net.PacketConn
	var size int
	var addr net.Addr
	buffer := make([]byte, 1024)

	for {
		packet, err = net.ListenPacket("udp4", fmt.Sprintf(":%d", PORT))
		if err != nil {
			fmt.Printf("Error listening for packet on %d: %s\n", PORT, err.Error())
			break
		}

		size, addr, err = packet.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("Failed to read from address: %s", err.Error())
			packet.Close()
			continue
		}
		fmt.Printf("Read %d bytes from %s: %s\n", size, addr.String(), buffer[:size])

		if err = parseMessage(buffer[:size]); err != nil {
			fmt.Printf("Error parsing message: %s", err.Error())
		}

		packet.Close()
	}

	Quit()
}

func parseMessage(message []byte) error {
	sepLoc := bytes.IndexRune(message, ':')
	if sepLoc < 0 || sepLoc >= len(message) {
		return errors.New("malformed message")
	}

	action := string(message[:sepLoc])
	payload := message[sepLoc+1:]

	switch action {
	case "notify":
		return actionNotify(payload)
	case "alert":
		return actionAlert(payload)
	case "beep":
		return actionBeep(payload)
	case "uri":
		return actionURI(payload)
	default:
		return fmt.Errorf("unknown message action %s", action)
	}
}
