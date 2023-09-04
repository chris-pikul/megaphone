package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gen2brain/beeep"
)

var SEPARATOR = []byte(":")

func initializeActions() error {
	var err error
	iconPath, err = filepath.Abs("./assets/information.png")
	if err != nil {
		return fmt.Errorf("failed to load application icon: %s", err.Error())
	}
	return nil
}

func actionNotify(payload []byte) error {
	parts := bytes.Split(payload, SEPARATOR)
	if len(parts) < 2 {
		return errors.New("malformed notify action, expected at least two parts")
	}

	title := string(parts[0])
	message := string(parts[1])

	return beeep.Notify(title, message, iconPath)
}

func actionAlert(payload []byte) error {
	parts := bytes.Split(payload, SEPARATOR)
	if len(parts) < 2 {
		return errors.New("malformed alert action, expected at least two parts")
	}

	title := string(parts[0])
	message := string(parts[1])

	return beeep.Alert(title, message, iconPath)
}

func actionBeep(payload []byte) error {
	parts := bytes.Split(payload, SEPARATOR)
	if len(parts) >= 2 {
		freq, err := strconv.ParseFloat(string(parts[0]), 32)
		if err != nil {
			return errors.New("malformed frequency for beep messag")
		}

		dur, err := strconv.ParseUint(string(parts[1]), 10, 32)
		if err != nil {
			return errors.New("malformed duration for beep message")
		}
		if dur > 1000 {
			dur = 1000
		}

		return beeep.Beep(freq, int(dur))
	}

	return errors.New("malformed beep message")
}

func actionURI(payload []byte) error {
	var cmd string
	var args []string
	url := string(payload)

	switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}