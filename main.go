package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const HOST string = "255.255.255.255";
const PORT int = 5790

var (
	iconPath string
	quit chan bool
)

func main() {
	fmt.Println("Starting megaphone");

	{
		var err error
		iconPath, err = filepath.Abs("./assets/information.png")
		if err != nil {
			fmt.Println("Unable to resolve default icon for alerts!")
			os.Exit(-1)
			return
		}
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
