package chip8

import (
	"github.com/nsf/termbox-go"
)

type Display struct {
	// Colors for Termbox Display, initially set to the default color
	bg termbox.Attribute
	fg termbox.Attribute

	buffer [2048]uint8
	// (0,0) - - - - (63, 0)
	//  |              |
	//  |              |
	// (0, 31) - - - (63, 31)
}

func newDisplay() (*Display) {
	disp := new(Display)
	// Set foreground and background colors
	disp.bg = termbox.ColorDefault
	disp.fg = termbox.ColorDefault
	return disp
}

func (disp *Display) drawSprite(xStart uint8, yStart uint8, height uint16, memory []uint8) (uint8) {
	flipFlag :=  uint8(0)

	for i := uint8(0); i < uint8(height); i++ {
		graphicsCoord := (yStart + i) * 64 + xStart
		// For each row of the sprite
		currByte := memory[i]
		for j := uint8(0); j < 8; j++ {
			// For each bit in the row of the sprite
			if (currByte >> j) & 1 == 1  { // If the bit in row is 'on'

				if disp.buffer[graphicsCoord] == 1 {
					termbox.SetCell(int(xStart + j), int(yStart + i), ' ', disp.fg, disp.bg)
					flipFlag = 1
				} else {
					termbox.SetCell(int(xStart + j), int(yStart + i), '*', disp.fg, disp.bg)
				}
			}
			graphicsCoord++
		}
	}
	return flipFlag
}

// Draw the buffer to the screen
func (disp *Display) flushDisplay() {
	termbox.Flush()
}

// Clear the display and internal buffer
func (disp *Display) clear() {
	termbox.Clear(disp.fg, disp.bg)
	for i := 0; i < 2048; i++ {
		disp.buffer[i] = 0
	}
}

func initTermbox() {
	// Initialize termbox-go
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	// Hide the termbox cursor
	termbox.HideCursor()
}

func closeTermbox() {
	termbox.Close()
}