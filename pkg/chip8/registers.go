package chip8

type Registers struct {
	Index int
}

func newRegisters() *Registers {
	registers := Registers{}
	return &registers
}
