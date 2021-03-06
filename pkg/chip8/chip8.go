package chip8

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type quirks struct {
	shiftLoadsYRegister                 bool
	storeAndLoadIncrementsIndexRegister bool
	setOverflowOnAddToIndex             bool
}

type cpu struct {
	memory     *Ram
	screen     *Screen
	registers  *Registers
	stack      *Stack
	pc         int
	delayTimer byte

	config quirks
	// clock cycles per second
	cpuHz int
	// sound/delay timer decay per second
	timerHz int
	// display updates per seconds
	displayHz int
	random    *rand.Rand
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
		cpuHz:     500,
		timerHz:   60,
		displayHz: 60,
		random:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return &cpu
}

func (cpu *cpu) tick() bool {
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
			return false
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
			if cpu.registers.VariableRegisters[n2] > (0xFF - cpu.registers.VariableRegisters[n3]) {
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
			cpu.registers.VariableRegisters[n2] = cpu.registers.VariableRegisters[n2] - cpu.registers.VariableRegisters[n3]
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
	// [CXNN] Set register to random number masked by NN
	case 0xC:
		handled = true
		cpu.registers.VariableRegisters[n2] = byte(cpu.random.Int()) & b2
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
		if b2 == 0x07 { // [FX07] Load delay timer
			handled = true
			cpu.registers.VariableRegisters[n2] = cpu.delayTimer
		} else if b2 == 0x15 { // [FX15] Set delay timer
			handled = true
			cpu.delayTimer = cpu.registers.VariableRegisters[n2]
		} else if b2 == 0x1E { // [FX1E] Add to index
			handled = true
			if cpu.config.setOverflowOnAddToIndex {
				if int(cpu.registers.Index)+int(cpu.registers.VariableRegisters[n2]) > 0xFFF {
					cpu.registers.VariableRegisters[0xF] = 1
				} else {
					cpu.registers.VariableRegisters[0xF] = 0
				}
			}
			cpu.registers.Index = (cpu.registers.Index + int(cpu.registers.VariableRegisters[n2])) % 0x1000
		} else if b2 == 0x29 { // [FX29] load address of font char
			handled = true
			cpu.registers.Index = cpu.memory.getAddressForFontChar(n2)
		} else if b2 == 0x33 { // [FX33] binary-coded decimal conversion
			handled = true
			vx := cpu.registers.VariableRegisters[n2]
			hundreds := vx / 100
			tens := (vx % 100) / 10
			ones := (vx % 10)
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
		fmt.Printf("[b1 0x%02X] [b2 0x%02X] [n1 0x%X] [n2 0x%X] [n3 0x%X] [n4 0x%X]\n", b1, b2, n1, n2, n3, n4)
		panic(fmt.Sprintf("Unhandled instruction [0x%02X%02X] at pc [0x%04X] adjusted pc [0x%04X]", b1, b2, cpu.pc-2, cpu.pc-2-0x200))
	}

	return true
}

func max(a byte, b byte) byte {
	if a > b {
		return a
	} else {
		return b
	}
}

// https://tobiasvl.github.io/blog/write-a-chip-8-emulator

func runRom(rom []byte) {
	cpu := newCpu(rom)

	cpuTickEveryMs := 1000 / cpu.cpuHz
	delayTickEveryMs := 1000 / cpu.timerHz
	displayTickEveryMs := 1000 / cpu.displayHz

	cpuTimer := cpuTickEveryMs
	delayTimer := delayTickEveryMs
	displayTimer := displayTickEveryMs

	lastTick := time.Now()

gameloop:
	for {
		frameStart := time.Now()
		deltaT := int(time.Since(lastTick).Milliseconds())

		cpuTimer += deltaT
		delayTimer += deltaT
		displayTimer += deltaT

		for cpuTimer >= cpuTickEveryMs {
			cpuTimer -= cpuTickEveryMs
			if !cpu.tick() {
				cpu.screen.doDraw()
				break gameloop
			}
		}

		if displayTimer >= displayTickEveryMs {
			displayTimer = 0
			cpu.screen.doDraw()
		}

		for delayTimer >= delayTickEveryMs {
			delayTimer -= delayTickEveryMs
			cpu.delayTimer = max(0, cpu.delayTimer-1)
		}

		lastTick = time.Now()
		frameElapsed := int(time.Since(frameStart).Milliseconds())

		time.Sleep(time.Duration(delayTickEveryMs-frameElapsed) * time.Millisecond)
	}
	// TODO(jpr): stop sound
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
