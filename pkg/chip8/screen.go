package chip8

import (
	"fmt"
)

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
	screen.Clear()
	return &screen
}

func (s *Screen) Clear() {
	s.pixels = make([]bool, s.columns*s.rows)
	s.screenBuffer = make([]rune, s.columns*s.rows)
	for i := 0; i < len(s.screenBuffer); i++ {
		s.screenBuffer[i] = s.offRune
	}
}

func (s *Screen) Draw(x_coord byte, y_coord byte, spriteData []byte) bool {
	didTurnOffPixel := false

	x_coord = x_coord % (byte(s.columns) - 1)
	y_coord = y_coord % (byte(s.rows) - 1)

	for currentY := y_coord; currentY < byte(s.rows) && int(currentY) < (int(y_coord)+len(spriteData)); currentY++ {
		offset := y_coord*byte(s.columns) + x_coord
		if spriteData[currentY-y_coord] == 0 {
			// flipping a pixel from on to off
			if s.pixels[offset] {
				didTurnOffPixel = true
			}
			s.pixels[offset] = false
			s.screenBuffer[offset] = s.offRune
		} else {
			s.pixels[offset] = true
			s.screenBuffer[offset] = s.onRune
		}
	}
	s.doDraw()

	return didTurnOffPixel
}

func (s *Screen) doDraw() {
	for row := 0; row < s.rows; row++ {
		fmt.Println(string(s.screenBuffer[row*s.columns : row*s.columns+s.columns]))
	}
}
