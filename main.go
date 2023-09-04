package main

import (
	"fmt"
	"os"
)

const HOST string = "255.255.255.255";
const PORT int = 5790

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		command := string(args[0])
		switch command {
		case "listen":
			startListening()
		case "send":
			broadcast(args[1:])
		case "help":
			printHelp()
		default:
			fmt.Printf("Unknown command %s\n\n", command)
			printHelp()
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

func printHelp() {
	fmt.Println(`Megaphone - Broadcast actions and messages across your network

megaphone
megaphone listen
--------------------------------------------------------------------------------
Listens for any broadcasted actions and performs them as necessary.


megaphone send [args...]
--------------------------------------------------------------------------------
Broadcasts a new action for all listeners on the network. Note: there is a
1024 byte limit to any actions broadcasted. Your messages may be truncated if
above that.

If no arguments are made, it will prompt for a message to send, and broadcast it
as a "notice".

If one argument is made, it is interpreted as a message body and a generic
"notice" action is broadcasted with the given argument as it's message body.

If more than one argument is made, it is considered a full action and is parsed
and broadcasted as such:
megaphone send notice [title] <message>
megaphone send alert [title] <message>
megaphone send beep [frequency] <duration>`)
}