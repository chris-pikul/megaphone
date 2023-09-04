package main

import (
	"fmt"
	"net"
)

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

	message := []byte("notify:Hello!:How are you?")
	if _, err = conn.Write(message); err != nil {
		fmt.Println("Failed to write to UDP broadcast channel")
		return err
	}

	fmt.Println("Message sent")
	return nil
}