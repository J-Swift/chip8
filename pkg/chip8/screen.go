package chip8

import (
	"fmt"
	"time"
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
	screen.resetBuffers()
	return &screen
}

func (s *Screen) resetBuffers() {
	s.pixels = make([]bool, int(s.columns)*int(s.rows))
	s.screenBuffer = make([]rune, int(s.columns)*int(s.rows))
	for i := 0; i < len(s.screenBuffer); i++ {
		s.screenBuffer[i] = s.offRune
	}
}

func (s *Screen) Clear() {
	s.resetBuffers()
	s.doDraw()
}

func (s *Screen) Draw(x_coord int, y_coord int, spriteData []byte) bool {
	didTurnOffPixel := false

	wrapped_x_coord := x_coord % (s.columns - 1)
	wrapped_y_coord := y_coord % (s.rows - 1)

	for row := 0; row < len(spriteData); row++ {
		for col := 0; col < 8; col++ {
			target_x_coord := wrapped_x_coord + col
			target_y_coord := wrapped_y_coord + row
			if target_x_coord >= s.columns || target_y_coord >= s.rows {
				continue
			}
			screenOffset := target_y_coord*s.rows + target_x_coord
			currentSpriteBitIsSet := ((spriteData[row] << col) & 0b10000000) > 0

			if currentSpriteBitIsSet {
				// flipping a pixel from on to off
				if s.pixels[screenOffset] {
					didTurnOffPixel = true
					s.pixels[screenOffset] = false
					s.screenBuffer[screenOffset] = s.offRune
				} else {
					s.pixels[screenOffset] = true
					s.screenBuffer[screenOffset] = s.onRune
				}
			}
		}
	}

	s.doDraw()

	return didTurnOffPixel
}

func (s *Screen) doDraw() {
	fmt.Printf("\033[0;0H")
	for row := 0; row < s.rows; row++ {
		fmt.Println(string(s.screenBuffer[int(row)*int(s.columns) : int(row)*int(s.columns)+int(s.columns)]))
	}
	time.Sleep(50 * time.Millisecond)
}
