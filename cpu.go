package chip8

import (
	"fmt"
)

type cpu struct {
	// Memory is 4096 bytes
	memory [4096]uint8
	// Graphics is 2048 bits
	graphics [2048]uint8
	// There are 16 registers, each with 8 bits of memory
	reg [16]uint8
	// Memory address register I of 16 bits
	I uint16
	// Program Counter register PC of 16 bits
	pc uint16
	// Timer delay and sound delay
	timerDelay uint8
	soundDelay uint8
	// Stack limited to 16 subcalls
	stack [16]uint16
	// Stack pointer
	sp uint16
}

// Returns a new instance of a Chip-8 CPU
func newCpu() *cpu {
  // New returns a pointer to a new zero'ed out cpu struct
	cp8 := new(cpu)
	// Set stack pointer and program counter
	cp8.sp = 0x200
	cp8.pc = 0xfa0

	return cp8
}

// Emulate one cycle of Chip-8
func emulateOneCycle() {
  
}
