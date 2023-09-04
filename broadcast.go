package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func broadcast(args []string) {
	if len(args) == 0 {
		promptNotice(false)
	} else if len(args) == 1 {
		if err := sendMessage([]byte(fmt.Sprintf("notify:Notice:%s", args[0]))); err != nil {
			fmt.Printf("Failed to send message: %s\n", err.Error())
			os.Exit(-1)
			return
		}
	} else {
		promptAdvanced(args)
	}

	fmt.Println("Sent message")
}

func prompt(message string) (string, error) {
	fmt.Printf("%s\n> ", message)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	} else if len(text) == 0 {
		return "", nil
	} else if strings.ContainsRune(text, ':') {
		fmt.Println("The : characters is reserved")
		return "", nil
	}

	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\n", " ")
	return strings.TrimSpace(text), nil
}

func promptNotice(isAlert bool) {
	PROMPT_TITLE:
	title, err := prompt("Title of your message:")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
	} else if len(title) == 0 {
		goto PROMPT_TITLE
	}

	PROMPT_BODY:
	message, err := prompt("Body of your message:")
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
	} else if len(message) == 0 {
		goto PROMPT_BODY
	}

	action := "notify"
	if isAlert {
		action = "alert"
	}

	if err := sendMessage([]byte(fmt.Sprintf("%s:%s:%s", action, title, message))); err != nil {
		fmt.Printf("Failed to send message: %s\n", err.Error())
		os.Exit(-1)
		return
	}
}

func promptAdvanced(args []string) {
	action := args[0]
	switch strings.ToLower(action) {
	case "notice":
		if len(args) == 2 {
			if err := sendMessage([]byte(fmt.Sprintf("notice:Notice:%s", args[1]))); err != nil {
				fmt.Printf("Failed to send message: %s\n", err.Error())
				os.Exit(-1)
				return
			}
		} else {
			if err := sendMessage([]byte(fmt.Sprintf("notice:%s:%s", args[1], args[2]))); err != nil {
				fmt.Printf("Failed to send message: %s\n", err.Error())
				os.Exit(-1)
				return
			}
		}
	case "alert":
		if len(args) == 2 {
			if err := sendMessage([]byte(fmt.Sprintf("alert:Notice:%s", args[1]))); err != nil {
				fmt.Printf("Failed to send message: %s\n", err.Error())
				os.Exit(-1)
				return
			}
		} else {
			if err := sendMessage([]byte(fmt.Sprintf("alert:%s:%s", args[1], args[2]))); err != nil {
				fmt.Printf("Failed to send message: %s\n", err.Error())
				os.Exit(-1)
				return
			}
		}
	case "beep":
		if len(args) == 2 {
			if err := sendMessage([]byte(fmt.Sprintf("beep:1000:%s", args[1]))); err != nil {
				fmt.Printf("Failed to send message: %s\n", err.Error())
				os.Exit(-1)
				return
			}
		} else {
			if err := sendMessage([]byte(fmt.Sprintf("beep:%s:%s", args[1], args[2]))); err != nil {
				fmt.Printf("Failed to send message: %s\n", err.Error())
				os.Exit(-1)
				return
			}
		}
	default:
		fmt.Printf("Unknown action %s", action)
		os.Exit(-1)
	}
}

func sendMessage(message []byte) error {
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

	if _, err = conn.Write(message); err != nil {
		fmt.Println("Failed to write to UDP broadcast channel")
		return err
	}

	fmt.Println("Message sent")
	return nil
}