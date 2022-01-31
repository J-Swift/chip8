package chip8

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

// https://tobiasvl.github.io/blog/write-a-chip-8-emulator

// 00E0 (clear screen)
// 1NNN (jump)
// 6XNN (set register VX)
// 7XNN (add value to register VX)
// ANNN (set index register I)
// DXYN (display/draw)

func readBytes(bytes string, numBytes int) (string, string) {
	if bytes == "" {
		return "", ""
	}
	if numBytes >= len(bytes) {
		return bytes, ""
	}
	return bytes[0:numBytes], bytes[numBytes:]
}

func printRom(rom string) {
	pc := 0
	var chars string
	for len(rom) > 0 {
		fmt.Printf("%08x ", pc)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s ", chars)
		chars, rom = readBytes(rom, 2)
		fmt.Printf("%s\n", chars)
		pc += 16
	}
}

func loadRom(romPath string) (string, error) {
	bytes, err := ioutil.ReadFile(romPath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Run(romPath string) {
	fmt.Printf("Running [%s]...\n", romPath)

	rom, err := loadRom(romPath)
	if err != nil {
		fmt.Printf("Error loading rom: %s", err.Error())
		os.Exit(1)
	}
	fmt.Println()
	printRom(rom)
	fmt.Println()

	fmt.Println("Done.")
}
