package chip8

import (
	"testing"
)

func TestSanityCheck(t *testing.T) {
	rom := []byte{0x00, 0xE0}
	cpu := newCpu(rom)
	cpu.screen.drawingEnabled = false

	if cpu.pc != 0x200 {
		t.Errorf("pc should have been 0x200 but was [%X]", cpu.pc)
	}
	if cpu.registers.Index != 0 {
		t.Errorf("index register should be zeroed but was [%d]", cpu.registers.Index)
	}
	for i := 0; i < len(cpu.registers.VariableRegisters); i++ {
		if cpu.registers.VariableRegisters[i] != 0 {
			t.Errorf("registers should be zeroed but [V%X] was [%d]", i, cpu.registers.VariableRegisters[i])
		}
	}
	for i := 0; i < len(cpu.screen.pixels); i++ {
		if cpu.screen.pixels[i] {
			y := i / cpu.screen.columns
			x := (i - y*cpu.screen.columns)
			t.Errorf("screen should be blanked but [%dx%d] was lit", x, y)
		}
	}

	if len(cpu.stack.innerStack) > 0 {
		t.Errorf("stack should be empty")
	}

	cpu.tick()
	if cpu.pc != 0x202 {
		t.Errorf("pc should have advanced to 0x202 but its now [0x%X]", cpu.pc)
	}
}

// 00E0
func TestClearScreen(t *testing.T) {
	rom := []byte{0x00, 0xE0}
	cpu := newCpu(rom)
	cpu.screen.drawingEnabled = false
	for i := 0; i < len(cpu.screen.pixels); i++ {
		cpu.screen.pixels[i] = true
	}
	cpu.tick()
	for i := 0; i < len(cpu.screen.pixels); i++ {
		if cpu.screen.pixels[i] {
			y := i / cpu.screen.columns
			x := (i - y*cpu.screen.columns)
			t.Errorf("screen should be blanked but [%dx%d] was lit", x, y)
		}
	}
}

// 00EE
func TestReturnFromSubroutine(t *testing.T) {
	t.Skip("TODO: 00EE")
}

// 1NNN
func TestJumpToAddress(t *testing.T) {
	rom := []byte{0x1F, 0xAB}
	cpu := newCpu(rom)
	cpu.tick()
	if cpu.pc != 0xFAB {
		t.Errorf("jump should set pc to 0xFAB but was [0x%X]", cpu.pc)
	}
}

// 2NNN
func TestCallSubroutineAtAddress(t *testing.T) {
	t.Skip("TODO: 2NNN")
}

// 3XNN
func TestSkipIfVxEqualToNumber(t *testing.T) {
	t.Skip("TODO: 3XNN")
}

// 4XNN
func TestSkipIfVxNotEqualToNumber(t *testing.T) {
	t.Skip("TODO: 4XNN")
}

// 5XY0
func TestSkipIfVxEqualToVy(t *testing.T) {
	t.Skip("TODO: 5XY0")
}

// 6XNN
func TestSetVxToNumber(t *testing.T) {
	t.Skip("TODO: 6XNN")
}

// 7XNN
func TestAddNumberToVx(t *testing.T) {
	t.Skip("TODO: 7XNN")
}

// 8XY0
func TestSetVxToVyDirect(t *testing.T) {
	t.Skip("TODO: 8XY0")
}

// 8XY1
func TestSetVxToVyBinaryOR(t *testing.T) {
	t.Skip("TODO: 8XY1")
}

// 8XY2
func TestSetVxToVyBinaryAND(t *testing.T) {
	t.Skip("TODO: 8XY2")
}

// 8XY3
func TestSetVxToVyBinaryXOR(t *testing.T) {
	t.Skip("TODO: 8XY3")
}

// 8XY4
func TestAddVyToVxWithCarry(t *testing.T) {
	t.Skip("TODO: 8XY4")
}

// 8XY5
func TestSubtractVyFromVxWithCarry(t *testing.T) {
	t.Skip("TODO: 8XY5")
}

// 8XY6
func TestShiftVxRightWithCarry(t *testing.T) {
	t.Skip("TODO: 8XY6")
}

// 8XY7
func TestSubtractVxFromVyWithCarry(t *testing.T) {
	t.Skip("TODO: 8XY7")
}

// 8XYE
func TestShiftVxLeftWithCarry(t *testing.T) {
	t.Skip("TODO: 8XYE")
}

// 9XY0
func TestSkipIfVxNotEqualToVy(t *testing.T) {
	t.Skip("TODO: 9XY0")
}

// ANNN
func TestSetIndexToNumber(t *testing.T) {
	t.Skip("TODO: ANNN")
}

// DXYN
func TestDrawSprite(t *testing.T) {
	t.Skip("TODO: DXYN")
}

// FX1E
func TestAddVxToIndex(t *testing.T) {
	t.Skip("TODO: FX1E")
}

// FX33
func TestBinaryCodedDecimalConversion(t *testing.T) {
	t.Skip("TODO: FX33")
}

// FX55
func TestStoreRegistersInMemory(t *testing.T) {
	t.Skip("TODO: FX55")
}

// FX65
func TestLoadRegistersFromMemory(t *testing.T) {
	t.Skip("TODO: FX65")
}
