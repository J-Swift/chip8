package chip8

type Registers struct {
	Index int

	VariableRegisters []byte
}

func newRegisters() *Registers {
	registers := Registers{}
	registers.VariableRegisters = make([]byte, 16)
	return &registers
}
