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
	rom := []byte{0x00, 0xEE}
	cpu := newCpu(rom)
	cpu.stack.push(0x321)
	cpu.tick()
	if cpu.pc != 0x321 {
		t.Errorf("return should set pc to stack pointer but it was [0x%X]", cpu.pc)
	}
}

// 1NNN
func TestJumpToAddress(t *testing.T) {
	rom := []byte{0x1F, 0xAB}
	cpu := newCpu(rom)
	cpu.tick()
	if len(cpu.stack.innerStack) > 0 {
		t.Errorf("jump should not push previous pc on top of stack")
	}
	if cpu.pc != 0xFAB {
		t.Errorf("jump should set pc to 0xFAB but was [0x%X]", cpu.pc)
	}
}

// 2NNN
func TestCallSubroutineAtAddress(t *testing.T) {
	rom := []byte{0x2F, 0xAB}
	cpu := newCpu(rom)
	cpu.tick()
	topOfStack := cpu.stack.pop()
	if topOfStack != 0x202 {
		t.Errorf("CallSubroutine should push previous pc on top of stack but was [0x%X]", topOfStack)
	}
	if cpu.pc != 0xFAB {
		t.Errorf("CallSubroutine should set pc to 0xFAB but was [0x%X]", cpu.pc)
	}
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
	t.Run("AddNumberToVx no carry smoke test", func(t *testing.T) {
		rom := []byte{0x7B, 0x01}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xB] = 0xFF
		expected := byte(0x0)
		cpu.tick()
		if cpu.registers.VariableRegisters[0xB] != expected {
			t.Errorf("AddNumberToVx Index register should have gone to [0x%02X] but it was [0x%02X]", expected, cpu.registers.VariableRegisters[0xB])
		}
		if cpu.registers.VariableRegisters[0xF] != 0 {
			t.Errorf("AddNumberToVx carry register should not have been set when overflow occurred")
		}
	})

	for vx := byte(0x0); vx <= 0xE; vx++ {
		t.Run(fmt.Sprintf("AddNumberToVx register V%X", vx), func(t *testing.T) {
			rom := []byte{0x70 | vx, 0x23}
			cpu := newCpu(rom)
			cpu.registers.VariableRegisters[vx] = 0x33
			expected := cpu.registers.VariableRegisters[vx] + rom[1]
			cpu.tick()
			if cpu.registers.VariableRegisters[vx] != expected {
				t.Errorf("AddNumberToVx register [V%X] should have been set to [0x%02X] but it was [0x%02X]", vx, expected, cpu.registers.VariableRegisters[vx])
			}
		})
	}
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
					t.Errorf("SetRegisterToRegisterBinaryAND register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
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
					t.Errorf("SetRegisterToRegisterBinaryAND register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
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
			t.Errorf("AddRegistersWithCarry register [VA] should have been set to [0x%02X] but it was [0x%02X]", expected, cpu.registers.VariableRegisters[0xA])
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
			t.Errorf("AddRegistersWithCarry register [VA] should have been set to [0x%02X] but it was [0x%02X]", expected, cpu.registers.VariableRegisters[0xA])
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
					t.Errorf("AddRegistersWithCarry register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 0 {
					t.Errorf("AddRegistersWithCarry carry flag should not have been set when adding [0x%02X] and [0x%02X]", vx, vy)
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
					t.Errorf("AddRegistersWithCarry register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 1 {
					t.Errorf("AddRegistersWithCarry carry flag should have been set when adding [0x%02X] and [0x%02X]", vx, vy)
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
			t.Run(fmt.Sprintf("Subtract register V%X from register V%X no borrow", vy, vx), func(t *testing.T) {
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

			t.Run(fmt.Sprintf("Subtract register V%X from register V%X with borrow", vy, vx), func(t *testing.T) {
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
					t.Errorf("SubtractRegistersWithBorrow borrow flag should have been set when subtracting [0x%02X] from [0x%02X]", vy, vx)
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
func TestSubtractVxFromVyWithBorrow(t *testing.T) {
	t.Run("Subtract register VX from register VY no borrow smoketest", func(t *testing.T) {
		rom := []byte{0x8A, 0xB7}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x02
		cpu.registers.VariableRegisters[0xB] = 0x01
		expected := byte(0xFF)
		cpu.tick()
		if cpu.registers.VariableRegisters[0xA] != expected {
			t.Errorf("SubtractRegistersWithBorrowReverse register [VA] should have been set to [0x%02X] but it was [0x%02X]", expected, cpu.registers.VariableRegisters[0xA])
		}
		if cpu.registers.VariableRegisters[0xF] != 0 {
			t.Errorf("SubtractRegistersWithBorrowReverse borrow flag should not have been set when subtracting [0x%02X] from [0x%02X]", 0x01, 0x02)
		}
	})

	t.Run("Subtract register VX from register VY with borrow smoketest", func(t *testing.T) {
		rom := []byte{0x8A, 0xB7}
		cpu := newCpu(rom)
		cpu.registers.VariableRegisters[0xA] = 0x01
		cpu.registers.VariableRegisters[0xB] = 0x02
		expected := byte(0x01)
		cpu.tick()
		if cpu.registers.VariableRegisters[0xA] != expected {
			t.Errorf("SubtractRegistersWithBorrowReverse register [VA] should have been set to [0x%02X] but it was [0x%02X]", expected, cpu.registers.VariableRegisters[0xA])
		}
		if cpu.registers.VariableRegisters[0xF] != 1 {
			t.Errorf("SubtractRegistersWithBorrowReverse borrow flag should have been set when subtracting [0x%02X] from [0x%02X]", 0x01, 0x02)
		}
	})

	// NOTE(jpr): no 0xF because its used for borrow flag
	for vx := byte(0x0); vx <= 0xE; vx++ {
		for vy := byte(0x0); vy <= 0xE; vy++ {
			if vx == vy {
				continue
			}
			t.Run(fmt.Sprintf("Subtract register V%X from register V%X no borrow", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x07 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0xAB
				cpu.registers.VariableRegisters[vy] = 0x33
				expected := cpu.registers.VariableRegisters[vy] - cpu.registers.VariableRegisters[vx]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SubtractRegistersWithBorrowReverse register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 0 {
					t.Errorf("SubtractRegistersWithBorrowReverse borrow flag should not have been set when subtracting [0x%02X] from [0x%02X]", vx, vy)
				}
			})

			t.Run(fmt.Sprintf("Subtract register V%X from register V%X with borrow", vx, vy), func(t *testing.T) {
				rom := []byte{0x80 | vx, 0x07 | (vy << 4)}
				cpu := newCpu(rom)
				cpu.registers.VariableRegisters[vx] = 0x33
				cpu.registers.VariableRegisters[vy] = 0xAB
				expected := cpu.registers.VariableRegisters[vy] - cpu.registers.VariableRegisters[vx]
				cpu.tick()
				if cpu.registers.VariableRegisters[vx] != expected {
					t.Errorf("SubtractRegistersWithBorrowReverse register [V%X] should have been set to [V%X] [0x%02X] but it was [0x%02X]", vx, vy, expected, cpu.registers.VariableRegisters[vx])
				}
				if cpu.registers.VariableRegisters[0xF] != 1 {
					t.Errorf("SubtractRegistersWithBorrowReverse borrow flag should have been set when subtracting [0x%02X] from [0x%02X]", vx, vy)
				}
			})
		}
	}
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
		t.Errorf("SetIndexRegister should have gone to 0xFAB when not equal but it was [0x%03X]", cpu.registers.Index)
	}
}

// DXYN
func TestDrawSprite(t *testing.T) {
	t.Skip("TODO: DXYN")
}

// FX1E
func TestAddVxToIndex(t *testing.T) {
	t.Run("config.setOverflowOnAddToIndex disabled - overflow smoketest", func(t *testing.T) {
		rom := []byte{0xFB, 0x1E}
		cpu := newCpu(rom)
		cpu.config.setOverflowOnAddToIndex = false
		cpu.registers.Index = 0xFFF
		cpu.registers.VariableRegisters[0xB] = 0x1
		expected := 0x0
		cpu.tick()
		if cpu.registers.Index != expected {
			t.Errorf("AddVxToIndex Index register should have gone to [0x%02X] but it was [0x%02X]", expected, cpu.registers.Index)
		}
		if cpu.registers.VariableRegisters[0xF] != 0 {
			t.Errorf("AddVxToIndex carry register should not have been set when setOverflowOnAddToIndex is disabled")
		}
	})

	t.Run("config.setOverflowOnAddToIndex enabled - no overflow smoketest", func(t *testing.T) {
		rom := []byte{0xFB, 0x1E}
		cpu := newCpu(rom)
		cpu.config.setOverflowOnAddToIndex = true
		cpu.registers.Index = 0xFFE
		cpu.registers.VariableRegisters[0xB] = 0x1
		expected := 0xFFF
		cpu.tick()
		if cpu.registers.Index != expected {
			t.Errorf("AddVxToIndex Index register should have gone to [0x%02X] but it was [0x%02X]", expected, cpu.registers.Index)
		}
		if cpu.registers.VariableRegisters[0xF] != 0 {
			t.Errorf("AddVxToIndex carry register should not have been set when setOverflowOnAddToIndex is enabled and no overflow occurred")
		}
	})

	t.Run("config.setOverflowOnAddToIndex enabled - overflow smoketest", func(t *testing.T) {
		rom := []byte{0xFB, 0x1E}
		cpu := newCpu(rom)
		cpu.config.setOverflowOnAddToIndex = true
		cpu.registers.Index = 0xFFF
		cpu.registers.VariableRegisters[0xB] = 0x1
		expected := 0x0
		cpu.tick()
		if cpu.registers.Index != expected {
			t.Errorf("AddVxToIndex Index register should have gone to [0x%02X] but it was [0x%02X]", expected, cpu.registers.Index)
		}
		if cpu.registers.VariableRegisters[0xF] != 1 {
			t.Errorf("AddVxToIndex carry register should have been set when setOverflowOnAddToIndex is enabled and overflow occurred")
		}
	})

	// NOTE(jpr): no 0xF because its used for borrow flag
	for registerIdx := byte(0x0); registerIdx <= 0xE; registerIdx++ {
		t.Run(fmt.Sprintf("Add V%X to Index register", registerIdx), func(t *testing.T) {
			rom := []byte{0xF0 | registerIdx, 0x1E}
			cpu := newCpu(rom)
			cpu.registers.Index = 0x23
			cpu.registers.VariableRegisters[registerIdx] = 0x12
			expected := 0x35
			cpu.tick()
			if cpu.registers.Index != expected {
				t.Errorf("AddVxToIndex Index register should have gone to [0x%02X] but it was [0x%02X]", expected, cpu.registers.Index)
			}
		})
	}
}

// FX33
func TestBinaryCodedDecimalConversion(t *testing.T) {
	t.Run("BinaryCodedDecimalConversion ones", func(t *testing.T) {
		rom := []byte{0xFB, 0x33}
		cpu := newCpu(rom)
		cpu.registers.Index = 0x500
		cpu.registers.VariableRegisters[0xB] = byte(1)
		cpu.tick()

		hundreds := cpu.memory.getAddress(cpu.registers.Index)
		tens := cpu.memory.getAddress(cpu.registers.Index + 1)
		ones := cpu.memory.getAddress(cpu.registers.Index + 2)
		if hundreds != 0 {
			t.Errorf("BCD should have set hundreds to 0 but was [%d]", hundreds)
		}
		if tens != 0 {
			t.Errorf("BCD should have set tens to 0 but was [%d]", tens)
		}
		if ones != 1 {
			t.Errorf("BCD should have set ones to 1 but was [%d]", ones)
		}
	})
	t.Run("BinaryCodedDecimalConversion ones", func(t *testing.T) {
		rom := []byte{0xFB, 0x33}
		cpu := newCpu(rom)
		cpu.registers.Index = 0x500
		cpu.registers.VariableRegisters[0xB] = byte(21)
		cpu.tick()

		hundreds := cpu.memory.getAddress(cpu.registers.Index)
		tens := cpu.memory.getAddress(cpu.registers.Index + 1)
		ones := cpu.memory.getAddress(cpu.registers.Index + 2)
		if hundreds != 0 {
			t.Errorf("BCD should have set hundreds to 0 but was [%d]", hundreds)
		}
		if tens != 2 {
			t.Errorf("BCD should have set tens to 2 but was [%d]", tens)
		}
		if ones != 1 {
			t.Errorf("BCD should have set ones to 1 but was [%d]", ones)
		}
	})
	t.Run("BinaryCodedDecimalConversion ones", func(t *testing.T) {
		rom := []byte{0xFB, 0x33}
		cpu := newCpu(rom)
		cpu.registers.Index = 0x500
		cpu.registers.VariableRegisters[0xB] = byte(213)
		cpu.tick()

		hundreds := cpu.memory.getAddress(cpu.registers.Index)
		tens := cpu.memory.getAddress(cpu.registers.Index + 1)
		ones := cpu.memory.getAddress(cpu.registers.Index + 2)
		if hundreds != 2 {
			t.Errorf("BCD should have set hundreds to 2 but was [%d]", hundreds)
		}
		if tens != 1 {
			t.Errorf("BCD should have set tens to 1 but was [%d]", tens)
		}
		if ones != 3 {
			t.Errorf("BCD should have set ones to 3 but was [%d]", ones)
		}
	})
}

// FX55
func TestStoreRegistersInMemory(t *testing.T) {
	t.Skip("TODO: FX55")
}

// FX65
func TestLoadRegistersFromMemory(t *testing.T) {
	t.Skip("TODO: FX65")
}
