package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/J-Swift/chip8/pkg/chip8"
)

func ensureRomExits(romPath string) {
	error := ""

	if romPath == "" {
		error = fmt.Sprintf("ROM does not exists [%s]", romPath)
	} else if _, err := os.Stat(romPath); errors.Is(err, os.ErrNotExist) {
		error = fmt.Sprintf("ROM does not exists [%s]", romPath)
	}

	if error != "" {
		fmt.Printf("ERROR: %s\n", error)
		os.Exit(1)
	}
}

func main() {
	romPtr := flag.String("rom", "", "Path to ROM")

	flag.Parse()

	ensureRomExits(*romPtr)

	chip8.Run(*romPtr)
}
