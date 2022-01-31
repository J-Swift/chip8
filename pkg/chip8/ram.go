package chip8

import "fmt"

type Ram struct {
	bytes []byte
}

func newRam(bytes []byte) *Ram {
	memspace := make([]byte, 4096)
	for i := 0; i < len(memspace); i++ {
		memspace[i] = 0
	}

	for i := 0; i < len(bytes); i++ {
		memspace[0x200+i] = bytes[i]
	}

	ram := Ram{bytes: memspace}
	return &ram
}

func (r *Ram) getAddress(address int) byte {
	if !(0 <= address && address <= len(r.bytes)-1) {
		panic(fmt.Sprintf("[%d] Invalid address [%d]", len(r.bytes), address))
	}

	return r.bytes[address]
}
