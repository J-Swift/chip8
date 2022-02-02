package chip8

import (
	"fmt"
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
	t.Run("when equal", func(t *testing.T) {
		rom := []byte{0x3A, 0x43}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x43
		cpu.tick()
		if cpu.pc != 0x204 {
			t.Errorf("SkipWhenEqual should have gone to 0x204 when equal but it was [0x%X]", cpu.pc)
		}
	})

	t.Run("when not equal", func(t *testing.T) {
		rom := []byte{0x3A, 0x43}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x42
		cpu.tick()
		if cpu.pc != 0x202 {
			t.Errorf("SkipWhenEqual should have gone to 0x202 when not equal but it was [0x%X]", cpu.pc)
		}
	})
}

// 4XNN
func TestSkipIfVxNotEqualToNumber(t *testing.T) {
	t.Run("when equal", func(t *testing.T) {
		rom := []byte{0x4A, 0x43}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x43
		cpu.tick()
		if cpu.pc != 0x202 {
			t.Errorf("SkipWhenNotEqual should have gone to 0x202 when equal but it was [0x%X]", cpu.pc)
		}
	})

	t.Run("when not equal", func(t *testing.T) {
		rom := []byte{0x4A, 0x43}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x42
		cpu.tick()
		if cpu.pc != 0x204 {
			t.Errorf("SkipWhenNotEqual should have gone to 0x204 when not equal but it was [0x%X]", cpu.pc)
		}
	})
}

// 5XY0
func TestSkipIfVxEqualToVy(t *testing.T) {
	t.Run("when equal", func(t *testing.T) {
		rom := []byte{0x5A, 0xB0}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x43
		cpu.registers.VariableRegisters[0xB] = 0x43
		cpu.tick()
		if cpu.pc != 0x204 {
			t.Errorf("SkipWhenEqualRegisters should have gone to 0x204 when equal but it was [0x%X]", cpu.pc)
		}
	})

	t.Run("when not equal", func(t *testing.T) {
		rom := []byte{0x5A, 0xB0}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x43
		cpu.registers.VariableRegisters[0xB] = 0x42
		cpu.tick()
		if cpu.pc != 0x202 {
			t.Errorf("SkipWhenEqualRegisters should have gone to 0x202 when not equal but it was [0x%X]", cpu.pc)
		}
	})
}

// 6XNN
func TestSetVxToNumber(t *testing.T) {
	for registerIdx := byte(0x0); registerIdx <= 0xF; registerIdx++ {
		t.Run(fmt.Sprintf("Set register V%X to number", registerIdx), func(t *testing.T) {
			rom := []byte{0x60 | registerIdx, 0xAB}
			cpu := newCpu(rom)
			cpu.tick()
			if cpu.registers.VariableRegisters[registerIdx] != 0xAB {
				t.Errorf("SetRegisterToNumber register [V%X] should have gone to 0xAB when not equal but it was [0x%X]", registerIdx, cpu.registers.VariableRegisters[registerIdx])
			}
		})
	}
}

// 7XNN
func TestAddNumberToVx(t *testing.T) {
	t.Skip("TODO: 7XNN")
}

// 8XY0
func TestSetVxToVyDirect(t *testing.T) {
	for vx := byte(0x0); vx <= 0xF; vx++ {
		for vy := byte(0x0); vy <= 0xF; vy++ {
			t.Run(fmt.Sprintf("Set register V%X to register V%X", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x00 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0x33
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SetRegisterToRegisterDirect register [V%X] should have been set to [V%X] [0x%X] but it was [0x%X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
			})
		}
	}
}

// 8XY1
func TestSetVxToVyBinaryOR(t *testing.T) {
	for vx := byte(0x0); vx <= 0xF; vx++ {
		for vy := byte(0x0); vy <= 0xF; vy++ {
			t.Run(fmt.Sprintf("BinaryOR register V%X with register V%X", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x01 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0x33
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vx] | cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SetRegisterToRegisterBinaryOR register [V%X] should have been set to [V%X] [0x%X] but it was [0x%X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
			})
		}
	}
}

// 8XY2
func TestSetVxToVyBinaryAND(t *testing.T) {
	for vx := byte(0x0); vx <= 0xF; vx++ {
		for vy := byte(0x0); vy <= 0xF; vy++ {
			t.Run(fmt.Sprintf("BinaryAND register V%X with register V%X", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x02 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0x33
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vx] & cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SetRegisterToRegisterBinaryAND register [V%X] should have been set to [V%X] [0x%X] but it was [0x%X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
			})
		}
	}
}

// 8XY3
func TestSetVxToVyBinaryXOR(t *testing.T) {
	for vx := byte(0x0); vx <= 0xF; vx++ {
		for vy := byte(0x0); vy <= 0xF; vy++ {
			t.Run(fmt.Sprintf("BinaryXOR register V%X with register V%X", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x03 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0x33
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vx] ^ cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SetRegisterToRegisterBinaryAND register [V%X] should have been set to [V%X] [0x%X] but it was [0x%X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
			})
		}
	}
}

// 8XY4
func TestAddVyToVxWithCarry(t *testing.T) {
	t.Run("Add register VX with register VY no overflow smoketest", func(t *testing.T) {
		rom := []byte{0x8A, 0xB4}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0xFE
		cpu.registers.VariableRegisters[0xB] = 0x01
		expected := byte(0xFF)
		cpu.tick()
		if cpu.registers.VariableRegisters[0xA] != expected {
			t.Errorf("AddRegistersWithCarry register [VA] should have been set to [255] but it was [0x%X]", cpu.registers.VariableRegisters[0xA])
		}
		if cpu.registers.VariableRegisters[0xF] != 0 {
			t.Errorf("AddRegistersWithCarry carry flag should not have been set when adding [0x%02X] and [0x%02X]", 0xFE, 0x01)
		}
	})

	t.Run("Add register VX with register VY with overflow smoketest", func(t *testing.T) {
		rom := []byte{0x8A, 0xB4}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0xFE
		cpu.registers.VariableRegisters[0xB] = 0x02
		expected := byte(0x0)
		cpu.tick()
		if cpu.registers.VariableRegisters[0xA] != expected {
			t.Errorf("AddRegistersWithCarry register [VA] should have been set to [0] but it was [0x%X]", cpu.registers.VariableRegisters[0xA])
		}
		if cpu.registers.VariableRegisters[0xF] != 1 {
			t.Errorf("AddRegistersWithCarry carry flag should have been set when adding [0x%02X] and [0x%02X]", 0xFE, 0x02)
		}
	})

	// NOTE(jpr): no 0xF because its used for carry flag
	for vx := byte(0x0); vx <= 0xE; vx++ {
		for vy := byte(0x0); vy <= 0xE; vy++ {
			if vx == vy {
				continue
			}
			t.Run(fmt.Sprintf("Add register V%X with register V%X no overflow", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x04 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0x33
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vx] + cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("AddRegistersWithCarry register [V%X] should have been set to [V%X] [0x%X] but it was [0x%X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 0 {
					t.Errorf("AddRegistersWithCarry carry flag should not have been set when adding [%X] and [%X]", vx, vy)
				}
			})

			t.Run(fmt.Sprintf("Add register V%X with register V%X with overflow", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x04 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0xFB
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vx] + cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("AddRegistersWithCarry register [V%X] should have been set to [V%X] [0x%X] but it was [0x%X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 1 {
					t.Errorf("AddRegistersWithCarry carry flag should have been set when adding [%X] and [%X]", vx, vy)
				}
			})
		}
	}
}

// 8XY5
func TestSubtractVyFromVxWithBorrow(t *testing.T) {
	t.Run("Subtract register VY from register VX no borrow smoketest", func(t *testing.T) {
		rom := []byte{0x8A, 0xB5}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x01
		cpu.registers.VariableRegisters[0xB] = 0x02
		expected := byte(0xFF)
		cpu.tick()
		if cpu.registers.VariableRegisters[0xA] != expected {
			t.Errorf("SubtractRegistersWithBorrow register [VA] should have been set to [0x%02X] but it was [0x%02X]", expected, cpu.registers.VariableRegisters[0xA])
		}
		if cpu.registers.VariableRegisters[0xF] != 0 {
			t.Errorf("SubtractRegistersWithBorrow borrow flag should not have been set when subtracting [0x%02X] from [0x%02X]", 0x02, 0x01)
		}
	})

	t.Run("Subtract register VY from register VX with borrow smoketest", func(t *testing.T) {
		rom := []byte{0x8A, 0xB5}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x02
		cpu.registers.VariableRegisters[0xB] = 0x01
		expected := byte(0x01)
		cpu.tick()
		if cpu.registers.VariableRegisters[0xA] != expected {
			t.Errorf("SubtractRegistersWithBorrow register [VA] should have been set to [0x%02X] but it was [0x%02X]", expected, cpu.registers.VariableRegisters[0xA])
		}
		if cpu.registers.VariableRegisters[0xF] != 1 {
			t.Errorf("SubtractRegistersWithBorrow borrow flag should have been set when subtracting [0x%02X] from [0x%02X]", 0x01, 0x02)
		}
	})

	// NOTE(jpr): no 0xF because its used for borrow flag
	for vx := byte(0x0); vx <= 0xE; vx++ {
		for vy := byte(0x0); vy <= 0xE; vy++ {
			if vx == vy {
				continue
			}
			t.Run(fmt.Sprintf("Add register V%X with register V%X no borrow", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x05 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0x33
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vx] - cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SubtractRegistersWithBorrow register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 0 {
					t.Errorf("SubtractRegistersWithBorrow borrow flag should not have been set when subtracting [0x%02X] from [0x%02X]", vy, vx)
				}
			})

			t.Run(fmt.Sprintf("Add register V%X with register V%X with borrow", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x05 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0xAB
				cpu.registers.VariableRegisters[vy] = 0x33
				expected := cpu.registers.VariableRegisters[vx] - cpu.registers.VariableRegisters[vy]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SubtractRegistersWithBorrow register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 1 {
					t.Errorf("SubtractRegistersWithBorrow borrow flag should have been set when subtracting [0x%02X] and [0x%02X]", vy, vx)
				}
			})
		}
	}
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
	t.Run("when equal", func(t *testing.T) {
		rom := []byte{0x9A, 0xB0}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x43
		cpu.registers.VariableRegisters[0xB] = 0x43
		cpu.tick()
		if cpu.pc != 0x202 {
			t.Errorf("SkipWhenNotEqualRegisters should have gone to 0x202 when equal but it was [0x%X]", cpu.pc)
		}
	})

	t.Run("when not equal", func(t *testing.T) {
		rom := []byte{0x9A, 0xB0}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x43
		cpu.registers.VariableRegisters[0xB] = 0x42
		cpu.tick()
		if cpu.pc != 0x204 {
			t.Errorf("SkipWhenNotEqualRegisters should have gone to 0x204 when not equal but it was [0x%X]", cpu.pc)
		}
	})
}

// ANNN
func TestSetIndexToNumber(t *testing.T) {
	rom := []byte{0xAF, 0xAB}
	cpu := newCpu(rom)
	cpu.tick()
	if cpu.registers.Index != 0xFAB {
		t.Errorf("SetIndexRegister should have gone to 0xFAB when not equal but it was [0x%X]", cpu.registers.Index)
	}
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
