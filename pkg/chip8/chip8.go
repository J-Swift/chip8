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
	stack := newStack()

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
		// fmt.Printf("Reading PC [%x]\n", pc)
		b1 = memory.getAddress(pc)
		b2 = memory.getAddress(pc + 1)
		pc += 2

		n1 = (b1 & 0b11110000) >> 4
		n2 = b1 & 0b00001111
		n3 = (b2 & 0b11110000) >> 4
		n4 = b2 & 0b00001111

		handled = false
		switch n1 {
		case 0x0:
			// [00E0] clear screen
			if b1 == 0x00 && b2 == 0xE0 {
				handled = true
				screen.Clear()
			}
		// [1NNN] jump to NNN
		case 0x1:
			handled = true
			combined = (int(n2) << 8) | (int(n3) << 4) | int(n4)
			if combined == pc-2 {
				fmt.Println("\nInfinite loop detected. Exiting....\n")
				return
			}
			pc = combined
		// [2NNN] call subroutine at NNN
		case 0x2:
			handled = true
			combined = (int(n2) << 8) | (int(n3) << 4) | int(n4)
			stack.push(pc)
			pc = combined
		// [3XNN] skip if VX equal to NN
		case 0x3:
			handled = true
			if registers.VariableRegisters[n2] == b2 {
				pc += 2
			}
		// [4XNN] skip if VX not equal to NN
		case 0x4:
			handled = true
			if registers.VariableRegisters[n2] != b2 {
				pc += 2
			}
		// [6XNN] set VX register to NN
		case 0x6:
			handled = true
			registers.VariableRegisters[n2] = b2
		// [7XNN] Add NN to VX register
		case 0x7:
			handled = true
			registers.VariableRegisters[n2] = byte((int(registers.VariableRegisters[n2]) + int(b2)) % 0x1FF)
		// [ANNN] set I register to NNN
		case 0xA:
			handled = true
			combined = (int(n2) << 8) | (int(n3) << 4) | int(n4)
			registers.Index = combined
		// [DXYN] Display N pixels of data at coord X,Y
		case 0xD:
			// fmt.Printf("[%x%x] [%d] [%d] [%d]\n", b1, b2, registers.VariableRegisters[0], registers.VariableRegisters[1], registers.Index)
			handled = true
			x_coord := registers.VariableRegisters[n2]
			y_coord := registers.VariableRegisters[n3]
			spriteData := memory.getAddressMulti(registers.Index, int(n4))
			if screen.Draw(int(x_coord), int(y_coord), spriteData) {
				registers.VariableRegisters[0xF] = 1
			} else {
				registers.VariableRegisters[0xF] = 0
			}
		case 0xF:
			// [FX33] binary-coded decimal conversion
			if b2 == 0x33 {
				handled = true
				vx := registers.VariableRegisters[n2]
				hundreds := vx / 100
				tens := (vx - (hundreds * 100)) / 10
				ones := vx - (hundreds * 100) - (tens * 10)
				memory.setAddress(registers.Index, hundreds)
				memory.setAddress(registers.Index+1, tens)
				memory.setAddress(registers.Index+2, ones)
			}
		}

		if !handled {
			fmt.Printf("[b1 %d] [b2 %d] [n1 %d] [n2 %d] [n3 %d] [n4 %d]\n", b1, b2, n1, n2, n3, n4)
			panic(fmt.Sprintf("Unhandled instruction [%02x%02x] at pc [%04x] adjusted pc [%04x]", b1, b2, pc-2, pc-2-0x200))
		}
	}
}

func loadRom(romPath string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(romPath)
	// _, err := ioutil.ReadFile(romPath)
	if err != nil {
		return nil, err
	}
	// bytes := []byte{
	// 	0x62, 0x0A, 0x63, 0x0C, 0xA2, 0x20, 0xD2, 0x36, 0x12, 0x40, 0xBA, 0x7C, 0xD6, 0xFE, 0x54, 0xAA,
	// }
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
