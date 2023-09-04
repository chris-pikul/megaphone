package main

import (
	"fmt"
	"os"
)

const HOST string = "255.255.255.255";
const PORT int = 5790

func main() {
	fmt.Println("Starting megaphone");

	
	args := os.Args[1:]

	if len(args) > 0 {
		command := string(args[0])
		switch command {
		case "listen":
			startListening()
		case "send":
			broadcast(args[1:])
		default:
			fmt.Printf("Unknown command %s", command)
			os.Exit(-1)
		}
	} else {
		// Default is listen
		startListening()
	}
}

func Quit() {
	fmt.Println("Requesting close")
	quit <- true
}