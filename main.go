package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const HOST string = "255.255.255.255";
const PORT int = 5790

var (
	quit chan bool
)

func main() {
	fmt.Println("Starting megaphone");

	quit = make(chan bool, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<- signals
		
		fmt.Println("Picked up close signal")
		Quit()
	}()

	go listen()

	time.Sleep(time.Second * 5)
	broadcast()

	// Wait for close signal
	fmt.Println("Waiting for close signal")
	<- quit

	fmt.Println("Closing megaphone")
}

func Quit() {
	fmt.Println("Requesting close")
	quit <- true
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

		packet.Close()
	}

	Quit()
}

func broadcast() error {
	fmt.Println("Attempting broadcast...")

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		fmt.Println("Failed to resolve broadcast address")
		return err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		fmt.Println("Failed to dial local UDP broadcast address")
		return err
	}
	defer conn.Close()

	message := []byte("Hello, World!")
	if _, err = conn.Write(message); err != nil {
		fmt.Println("Failed to write to UDP broadcast channel")
		return err
	}

	fmt.Println("Message sent")
	return nil
}
