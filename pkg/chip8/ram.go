package chip8

import "fmt"

type Ram struct {
	bytes            []byte
	fontStoredAt     int
	bytesPerFontChar int
}

var font []byte

func init() {
	font = []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
}

func newRam(bytes []byte) *Ram {
	ram := Ram{fontStoredAt: 0x050, bytesPerFontChar: 5}

	memspace := make([]byte, 4096)

	// load font into memory
	for i := 0; i < len(font); i++ {
		memspace[ram.fontStoredAt+i] = font[i]
	}

	// load rom into memory
	for i := 0; i < len(bytes); i++ {
		memspace[0x200+i] = bytes[i]
	}

	ram.bytes = memspace
	return &ram
}

func (r *Ram) getAddressForFontChar(c byte) int {
	return r.fontStoredAt + int(c*byte(r.bytesPerFontChar))
}

func (r *Ram) getAddress(address int) byte {
	if !(0 <= address && address <= len(r.bytes)-1) {
		panic(fmt.Sprintf("[%d] Invalid address [%d]", len(r.bytes), address))
	}

	return r.bytes[address]
}

func (r *Ram) getAddressMulti(address int, count int) []byte {
	if !(0 <= address && (address+count) <= len(r.bytes)-1) {
		panic(fmt.Sprintf("[%d] Invalid address [%d] count [%d]", len(r.bytes), address, count))
	}

	return r.bytes[address : address+count]
}

func (r *Ram) setAddress(address int, value byte) {
	r.bytes[address] = value
}
