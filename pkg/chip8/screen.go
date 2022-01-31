package chip8

import "fmt"

type Screen struct {
	columns      int
	rows         int
	offRune      rune
	onRune       rune
	pixels       []bool
	screenBuffer []rune
}

func newScreen() *Screen {
	screen := Screen{columns: 64, rows: 32, offRune: '⬛', onRune: '🟨'}

	screen.pixels = make([]bool, screen.columns*screen.rows)

	screen.screenBuffer = make([]rune, screen.columns*screen.rows)
	for i := 0; i < len(screen.screenBuffer); i++ {
		screen.screenBuffer[i] = screen.offRune
	}

	return &screen
}

func (s *Screen) Clear() {
	s.pixels = make([]bool, s.columns*s.rows)
}

func (s *Screen) Draw() {
	s.doDraw()
}

func (s *Screen) doDraw() {
	for row := 0; row < s.rows; row++ {
		fmt.Println(string(s.screenBuffer[row*s.columns : row*s.columns+s.columns]))
	}
}
