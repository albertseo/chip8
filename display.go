package chip8

import (
	"github.com/nsf/termbox-go"
)

type Display struct {
	// Colors for Termbox Display
	bg termbox.Attribute
	fg termbox.Attribute
}

func newDisplay() (*Display) {
	disp := new(Display)
	// Set foreground and background colors
	disp.bg = termbox.ColorDefault
	disp.fg = termbox.ColorDefault
	return disp
}

func (disp *Display) drawArray(arr [2048]uint8) {
	for i := 0; i < 2048; i++ {
		if arr[i] == 1 {
			// This cell will be filled
			termbox.SetCell(i % 64, i / 32, '*', disp.fg, disp.bg)
		} else {
			// This cell will not be filled
			termbox.SetCell(i % 64, i / 32, ' ', disp.fg, disp.bg)
		}
	}
}

func (disp *Display) drawSprite(xCord uint8, yCord uint8, height uint16, memory []uint8) {
	termbox.HideCursor()
}

// Draw the buffer to the screen
func (disp *Display) flushDisplay() {
	termbox.Flush()
}

// Clear the display
func (disp *Display) clear() {
	termbox.Clear(disp.fg, disp.bg)
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