package chip8

import (
	"testing"
	"fmt"
)

func TestDissassemble(t *testing.T) {
	fmt.Printf("Starting Dissassembly... \n")
	c8 := newCpu()
	c8.loadFile("Fishie.ch8")

	for i := 0; i < 10; i++ {
		c8.emulateOneCycle()
	}
}