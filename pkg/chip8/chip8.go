package chip8

import (
	"fmt"
	"io/ioutil"
	"os"
)

type quirks struct {
	shiftLoadsYRegister                 bool
	storeAndLoadIncrementsIndexRegister bool
	setOverflowOnAddToIndex             bool
}

type cpu struct {
	memory    *Ram
	screen    *Screen
	registers *Registers
	stack     *Stack
	pc        int
	config    quirks
}

func newCpu(romData []byte) *cpu {
	cpu := cpu{
		memory:    newRam(romData),
		screen:    newScreen(),
		registers: newRegisters(),
		stack:     newStack(),
		pc:        0x200,
		config: quirks{
			shiftLoadsYRegister:                 false,
			storeAndLoadIncrementsIndexRegister: false,
			setOverflowOnAddToIndex:             true,
		},
	}
	return &cpu
}

func (cpu *cpu) tick() {
	// fmt.Printf("Reading PC [%x]\n", pc)
	b1 := byte(cpu.memory.getAddress(cpu.pc))
	b2 := byte(cpu.memory.getAddress(cpu.pc + 1))
	cpu.pc += 2

	n1 := byte((b1 & 0b11110000) >> 4)
	n2 := byte(b1 & 0b00001111)
	n3 := byte((b2 & 0b11110000) >> 4)
	n4 := byte(b2 & 0b00001111)

	handled := false
	switch n1 {
	case 0x0:
		if b1 == 0x00 && b2 == 0xE0 { // [00E0] clear screen
			handled = true
			cpu.screen.Clear()
		} else if b1 == 0x00 && b2 == 0xEE { // [00EE] return from subroutine
			handled = true
			cpu.pc = cpu.stack.pop()
		}
	// [1NNN] jump to NNN
	case 0x1:
		handled = true
		combined := (int(n2) << 8) | (int(n3) << 4) | int(n4)
		if combined == cpu.pc-2 {
			fmt.Printf("\nInfinite loop detected. Exiting....\n\n")
			return
		}
		cpu.pc = combined
	// [2NNN] call subroutine at NNN
	case 0x2:
		handled = true
		combined := (int(n2) << 8) | (int(n3) << 4) | int(n4)
		cpu.stack.push(cpu.pc)
		cpu.pc = combined
	// [3XNN] skip if VX equal to NN
	case 0x3:
		handled = true
		if cpu.registers.VariableRegisters[n2] == b2 {
			cpu.pc += 2
		}
	// [4XNN] skip if VX not equal to NN
	case 0x4:
		handled = true
		if cpu.registers.VariableRegisters[n2] != b2 {
			cpu.pc += 2
		}
	// [5XY0] skip if VX equal to VY
	case 0x5:
		handled = true
		if cpu.registers.VariableRegisters[n2] == cpu.registers.VariableRegisters[n3] {
			cpu.pc += 2
		}
	// [6XNN] set VX register to NN
	case 0x6:
		handled = true
		cpu.registers.VariableRegisters[n2] = b2
	// [7XNN] Add NN to VX register
	case 0x7:
		handled = true
		cpu.registers.VariableRegisters[n2] = byte((int(cpu.registers.VariableRegisters[n2]) + int(b2)) % 0x1FF)
	case 0x8:
		switch n4 {
		// [8XY0] Set VX to VY
		case 0x0:
			handled = true
			cpu.registers.VariableRegisters[n2] = cpu.registers.VariableRegisters[n3]
		// [8XY1] Set VX to binary OR with VY
		case 0x1:
			handled = true
			cpu.registers.VariableRegisters[n2] |= cpu.registers.VariableRegisters[n3]
		// [8XY2] Set VX to binary AND with VY
		case 0x2:
			handled = true
			cpu.registers.VariableRegisters[n2] &= cpu.registers.VariableRegisters[n3]
		// [8XY3] Set VX to binary XOR with VY
		case 0x3:
			handled = true
			cpu.registers.VariableRegisters[n2] ^= cpu.registers.VariableRegisters[n3]
		// [8XY4] Add VX to VY with carry
		case 0x4:
			handled = true
			// check for overflow
			if int(cpu.registers.VariableRegisters[n2])+int(cpu.registers.VariableRegisters[n3]) > 255 {
				cpu.registers.VariableRegisters[0xF] = 1
			} else {
				cpu.registers.VariableRegisters[0xF] = 0
			}
			cpu.registers.VariableRegisters[n2] += cpu.registers.VariableRegisters[n3]
		// [8XY5] Subtract VY from VX with carry
		case 0x5:
			handled = true
			// check for underflow
			if cpu.registers.VariableRegisters[n2] > cpu.registers.VariableRegisters[n3] {
				cpu.registers.VariableRegisters[0xF] = 1
			} else {
				cpu.registers.VariableRegisters[0xF] = 0
			}
			cpu.registers.VariableRegisters[n2] -= cpu.registers.VariableRegisters[n3]
		// [8XY6] Shift VX right with carry
		case 0x6:
			handled = true
			if cpu.config.shiftLoadsYRegister {
				cpu.registers.VariableRegisters[n2] = cpu.registers.VariableRegisters[n3]
			}
			cpu.registers.VariableRegisters[0xF] = cpu.registers.VariableRegisters[n2] & 0b1
			cpu.registers.VariableRegisters[n2] >>= 1
		// [8XY7] Subtract VX from VY with carry
		case 0x7:
			handled = true
			if cpu.registers.VariableRegisters[n3] > cpu.registers.VariableRegisters[n2] {
				cpu.registers.VariableRegisters[0xF] = 1
			} else {
				cpu.registers.VariableRegisters[0xF] = 0
			}
			cpu.registers.VariableRegisters[n2] = cpu.registers.VariableRegisters[n3] - cpu.registers.VariableRegisters[n2]
		// [8XYE] Shift VX left with carry
		case 0xE:
			handled = true
			if cpu.config.shiftLoadsYRegister {
				cpu.registers.VariableRegisters[n2] = cpu.registers.VariableRegisters[n3]
			}
			cpu.registers.VariableRegisters[0xF] = (cpu.registers.VariableRegisters[n2] >> 7) & 0b1
			cpu.registers.VariableRegisters[n2] <<= 1
		}
	// [9XY0] skip if VX not equal to VY
	case 0x9:
		handled = true
		if cpu.registers.VariableRegisters[n2] != cpu.registers.VariableRegisters[n3] {
			cpu.pc += 2
		}
	// [ANNN] set I register to NNN
	case 0xA:
		handled = true
		combined := (int(n2) << 8) | (int(n3) << 4) | int(n4)
		cpu.registers.Index = combined
	// [DXYN] Display N pixels of data at coord X,Y
	case 0xD:
		// fmt.Printf("[%x%x] [%d] [%d] [%d]\n", b1, b2, cpu.registers.VariableRegisters[0], cpu.registers.VariableRegisters[1], cpu.registers.Index)
		handled = true
		x_coord := cpu.registers.VariableRegisters[n2]
		y_coord := cpu.registers.VariableRegisters[n3]
		spriteData := cpu.memory.getAddressMulti(cpu.registers.Index, int(n4))
		if cpu.screen.Draw(int(x_coord), int(y_coord), spriteData) {
			cpu.registers.VariableRegisters[0xF] = 1
		} else {
			cpu.registers.VariableRegisters[0xF] = 0
		}
	case 0xF:
		if b2 == 0x1E { // [FX1E] Add to index
			handled = true
			if cpu.config.setOverflowOnAddToIndex {
				if int(cpu.registers.Index)+int(cpu.registers.VariableRegisters[n2]) > 0xFFF {
					cpu.registers.VariableRegisters[0xF] = 1
				} else {
					cpu.registers.VariableRegisters[0xF] = 0
				}
			}
			cpu.registers.Index = (cpu.registers.Index + int(cpu.registers.VariableRegisters[n2])) % 0x1000
		} else if b2 == 0x33 { // [FX33] binary-coded decimal conversion
			handled = true
			vx := cpu.registers.VariableRegisters[n2]
			hundreds := vx / 100
			tens := (vx - (hundreds * 100)) / 10
			ones := vx - (hundreds * 100) - (tens * 10)
			cpu.memory.setAddress(cpu.registers.Index, hundreds)
			cpu.memory.setAddress(cpu.registers.Index+1, tens)
			cpu.memory.setAddress(cpu.registers.Index+2, ones)
		} else if b2 == 0x55 { // [FX55] store registers in memory
			handled = true
			currentAddress := cpu.registers.Index
			for currentRegister := byte(0); currentRegister <= n2; currentRegister++ {
				cpu.memory.setAddress(currentAddress, cpu.registers.VariableRegisters[currentRegister])
				currentAddress++
				if cpu.config.storeAndLoadIncrementsIndexRegister {
					cpu.registers.Index++
				}
			}
		} else if b2 == 0x65 { // [FX65] load registers from memory
			handled = true
			currentAddress := cpu.registers.Index
			for currentRegister := byte(0); currentRegister <= n2; currentRegister++ {
				cpu.registers.VariableRegisters[currentRegister] = cpu.memory.getAddress(currentAddress)
				currentAddress++
				if cpu.config.storeAndLoadIncrementsIndexRegister {
					cpu.registers.Index++
				}
			}
		}
	}

	if !handled {
		fmt.Printf("[b1 %d] [b2 %d] [n1 %d] [n2 %d] [n3 %d] [n4 %d]\n", b1, b2, n1, n2, n3, n4)
		panic(fmt.Sprintf("Unhandled instruction [%02x%02x] at pc [%04x] adjusted pc [%04x]", b1, b2, cpu.pc-2, cpu.pc-2-0x200))
	}
}

// https://tobiasvl.github.io/blog/write-a-chip-8-emulator

// 00E0 (clear screen)
// 1NNN (jump)
// 6XNN (set register VX)
// 7XNN (add value to register VX)
// ANNN (set index register I)
// DXYN (display/draw)

func runRom(rom []byte) {
	cpu := newCpu(rom)

	for {
		cpu.tick()
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
