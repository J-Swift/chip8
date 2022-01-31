package chip8

import (
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

func runRom(rom []byte) {
	memory := newRam(rom)
	screen := newScreen()
	registers := newRegisters()

	pc := 0x200

	var b1 int
	var b2 int
	var n1 int
	var n2 int
	var n3 int
	var n4 int
	var combined int

	var handled bool
	for {
		b1 = int(memory.getAddress(pc))
		pc += 1
		b2 = int(memory.getAddress(pc))
		pc += 1

		n1 = (b1 & 0b11110000) >> 4
		n2 = b1 & 0b00001111
		n3 = (b2 & 0b11110000) >> 4
		n4 = b2 & 0b00001111

		handled = false
		switch n1 {
		case 0x0:
			if b1 == 0x00 && b2 == 0xE0 {
				handled = true
				screen.clear()
			}
		case 0xA:
			handled = true
			combined = (n2 << 8) & (n3 << 4) & n4
			registers.Index = combined
		}

		if !handled {
			fmt.Printf("[b1 %d] [b2 %d] [n1 %d] [n2 %d] [n3 %d] [n4 %d]\n", b1, b2, n1, n2, n3, n4)
			panic(fmt.Sprintf("Unhandled instruction [%02x%02x] at pc [%04x] adjusted pc [%04x]", b1, b2, pc-2, pc-2-0x200))
		}
	}
}

func loadRom(romPath string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(romPath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Run(romPath string) {
	fmt.Printf("Running [%s]...\n\n", romPath)

	rom, err := loadRom(romPath)
	if err != nil {
		fmt.Printf("Error loading rom: %s", err.Error())
		os.Exit(1)
	}

	runRom(rom)

	fmt.Println("Done.")
}
