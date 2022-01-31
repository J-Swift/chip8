package chip8

type Screen struct {
	maxX   int
	maxY   int
	pixels []bool
}

func newScreen() *Screen {
	screen := Screen{maxX: 64, maxY: 32}
	screen.pixels = make([]bool, screen.maxX*screen.maxY)
	return &screen
}

func (s *Screen) clear() {
	s.pixels = make([]bool, s.maxX*s.maxY)
}
