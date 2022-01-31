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

	// Instruction helpers
	var b1 byte
	var b2 byte
	var n1 byte
	var n2 byte
	var n3 byte
	var n4 byte

	// Conditoinally used helpers
	var combined int

	var handled bool
	for {
		b1 = memory.getAddress(pc)
		pc += 1
		b2 = memory.getAddress(pc)
		pc += 1

		n1 = (b1 & 0b11110000) >> 4
		n2 = b1 & 0b00001111
		n3 = (b2 & 0b11110000) >> 4
		n4 = b2 & 0b00001111

		handled = false
		switch n1 {
		// clear screen
		case 0x0:
			if b1 == 0x00 && b2 == 0xE0 {
				handled = true
				screen.Clear()
			}
		// set VX register
		case 0x6:
			handled = true
			combined = (int(n2) << 8) & (int(n3) << 4) & int(n4)
			registers.Index = combined
		// set I register
		case 0xA:
			handled = true
			registers.VariableRegisters[n2] = b2
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
