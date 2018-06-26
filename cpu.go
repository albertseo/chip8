package chip8

import (
	"fmt"
  "io/ioutil"
  "math/rand"
)

type cpu struct {
	// Memory is 4096 bytes
	memory [4096]uint8
	// Graphics is 2048 bits
	graphics [2048]uint8
	// There are 16 registers, each with 8 bits of memory
	reg [16]uint8
	// Memory address register I of 16 bits
	i uint16
	// Program Counter register PC of 16 bits
	pc uint16
	// Timer delay and sound delay
	timerDelay uint8
	soundDelay uint8
	// Stack limited to 16 subcalls
	stack [16]uint16
	// Stack pointer
	sp uint16
  // Array to record which key is held
  key [16]int
}

// Returns a new instance of a Chip-8 CPU
func newCpu() *cpu {
  // New returns a pointer to a new zero'ed out cpu struct
	cp8 := new(cpu)
	// Set stack pointer and program counter
	cp8.pc = 0x200
	cp8.sp = 0xfa0

	return cp8
}

// Emulate one cycle of Chip-8
func (c8 *cpu) emulateOneCycle() {
  // Fetch instruction
  inst := c8.fetchInstruction()
  // Execute instruction
  c8.executeInstruction(inst)
  //update timers
}

// Retrives the current Instruction from memory
func (c8 *cpu) fetchInstruction() uint16 {
  return uint16(c8.memory[c8.pc]) << 8 | uint16(c8.memory[c8.pc + 1])
}

// Decodes and executes instruction
func (c8 *cpu) executeInstruction(inst uint16) {
  // Increment pc
  c8.pc += 2

  // Decode instruction  type
  switch inst & 0xF000 {
  case 0x0000:
    switch inst & 0x0FFF {
    case 0x00E0:
      // Clear screen
    case 0x00EE:
      // Return from subroutine
      c8.sp -= 1
      c8.pc = c8.stack[c8.sp]
    }
  case 0x1000:
    // JUMP to instruction at addresss 0x0NNN
    imm := inst & 0x0FFF
    c8.pc = imm
  case 0x2000:
    // CALL subroutine at address 0xNNN
    imm := inst & 0x0FFF
    c8.stack[c8.sp] = c8.pc
    c8.sp += 1
    c8.pc = imm
  case 0x3000:
    // SKIP next instruction if VX == imm
    regX := int(inst >> 8 & 0x0F)
    imm := inst & 0x00FF
    if uint16(c8.reg[regX]) == imm {
      c8.pc += 2
    }
  case 0x4000:
    // SKIP next instruction if VX != imm
    regX := int(inst >> 8 & 0x0F)
    imm := inst & 0x00FF
    if uint16(c8.reg[regX]) != imm {
      c8.pc += 2
    }
  case 0x5000:
    if inst & 0x000F == 0x0 {
      // SKIP next instruction if VX == VY
      regX := int(inst >> 8 & 0x0F)
      regY := int(inst >> 4 & 0x00F)
      if c8.reg[regX] == c8.reg[regY] {
        c8.pc += 2
      }
    }
  case 0x6000:
    // MOVE immediate into register VX
    regX := int(inst >> 8 & 0x0F)
    imm := inst & 0x00FF
    c8.reg[regX] = uint8(imm)
    // return preamble + fmt.Sprintf("MVI V%X, %02X", uint16(reg), uint16(imm))
  case 0x7000:
    // ADD imm to value at VX
    regX := int(inst >> 8 & 0x0F)
    imm := inst & 0x00FF
    c8.reg[regX] += uint8(imm)
  case 0x8000:
    // Operates with two registers
    regX := int(inst >> 8 & 0x0F)
    regY := int(inst >> 4 & 0x00F)
    switch inst & 0x000F {
    case 0x0:
      // Set VX to the value of VY
      c8.reg[regX] = c8.reg[regY]
    case 0x1:
      // Set VX to the value of VY | VX
      c8.reg[regX] = c8.reg[regY] | c8.reg[regX]
    case 0x2:
      // Set VX to the value of VY & VX
      c8.reg[regX] = c8.reg[regY] & c8.reg[regX]
    case 0x3:
      // Set VX to the value of VY ^ VX
      c8.reg[regX] = c8.reg[regY] ^ c8.reg[regX]
    case 0x4:
      // Set VX to the value of VY + VX, VF set to 1 if there is a carry over
      temp := c8.reg[regY] + c8.reg[regX]
      if temp > 0xFF {
        c8.reg[15] = 1
      } else {
        c8.reg[15] = 0
      }
      c8.reg[regX] = uint8(temp)
    case 0x5:
      // Set VX to the value of VX - VY, VF set to 0 if need to borrow
      if c8.reg[regX] > c8.reg[regY] {
        c8.reg[15] = 1
      } else {
        c8.reg[15] = 0
      }
      c8.reg[regX] = uint8(c8.reg[regX] - c8.reg[regY])
    case 0x6:
      // Set VX to the value of VY >> 1, VF set to least significant digit of VY
      c8.reg[15] = c8.reg[regX] & 1
      c8.reg[regX] = c8.reg[regX] >> 1
    case 0x7:
      // Set VX to the value of VY - VX, VH set to 0 if need to borrow
      if c8.reg[regY] > c8.reg[regX] {
        c8.reg[15] = 1
      } else {
        c8.reg[15] = 0
      }
      c8.reg[regX] = uint8(c8.reg[regY] - c8.reg[regX])
    case 0xE:
      // Set VX and VY to VY << 1, VF set to most significant digit of VY before shift.
      c8.reg[15] = c8.reg[regX] >> 7
      c8.reg[regX] = c8.reg[regX] << 1
    }
  case 0x9000:
    if inst & 0x000F == 0x0 {
      // SKIP next instruction if VX != VY
      regX := (inst >> 8 & 0x0F)
      regY := (inst >> 4 & 0x00F)
      if c8.reg[regX] != c8.reg[regY] {
        c8.pc += 2
      }
    }
  case 0xA000:
    // Set register I to imm
    imm := inst & 0x0FFF
    c8.i = uint16(imm)
  case 0xB000:
    // JUMP to address at imm + value at V0
    imm := inst & 0x0FFF
    c8.pc = uint16(c8.reg[0]) + imm
  case 0xC000:
    // Set register VX to Imm & rand(0,255)
    regX := inst >> 8 & 0x0F
    imm := inst & 0x00ff
    c8.reg[regX] = imm & uint8(Int())
  case 0xD000:
    // Draw stuff to the screen
    // regX := inst >> 8 & 0x0F
    // regY := inst >> 4 & 0x00F
    // N := inst & 0x000F
    // return preamble + fmt.Sprintf("SPRITE V%X, V%X, %X", uint16(regX), uint16(regY), uint16(N))
  case 0xE000:
    switch inst & 0x00FF {
    case 0x9E:
      // SKIP next instruction if key stored in VX is held
      regX := inst >> 8 & 0x0F
      if c8.key[c8.reg[regX]] == 1 {
        c8.pc += 2
      }
    case 0xA1:
      // SKIP next instruction if key stored in VX isn't held
      regX := inst >> 8 & 0x0F
      if c8.key[c8.reg[regX]] != 1 {
        c8.pc += 2
      }
      // return preamble + fmt.Sprintf("SKIP.NKEY V%X", uint16(reg))
    }
  case 0xF000:
    // reg := inst >> 8 & 0x0F
    switch inst & 0x00FF {
    case 0x07:
      // Set VX to value of delay timer
      // return preamble + fmt.Sprintf("MOV V%X DELAY", uint16(reg))
    case 0x0A:
      // Wait for keypress, then store in VX
      // return preamble + fmt.Sprintf("WAITKEY V%X", uint16(reg))
    case 0x15:
      // Set delay timer to value in VX
      // return preamble + fmt.Sprintf("MOV DELAY V%X", uint16(reg))
    case 0x18:
      // Set sound timer to value in VX
      // return preamble + fmt.Sprintf("MOV SOUND V%X", uint16(reg))
    case 0x1E:
      // ADDS VX to I
      // return preamble + fmt.Sprintf("ADD I V%X", uint16(reg))
    case 0x29:
      // Sets I to the location of sprite of character in VX
      // return preamble + fmt.Sprintf("SPRITECHAR I V%X", uint16(reg))
    case 0x33:
      // Stores BCD of VX at I, I+1, I+2
      // return preamble + fmt.Sprintf("MOVBCD V%X", uint16(reg))
    case 0x55:
      // Stores V0 to VX in memory starting at I
      // return preamble + fmt.Sprintf("STORE (I), V0-V%X", uint16(reg))
    case 0x65:
      // Load values at V0 to VX starting at memory address I
      // return preamble + fmt.Sprintf("LOAD V0-V%X, (I)", uint16(reg))
    }
  }
}

func (c8 *cpu) loadFile(fileName string) {
  buffer, err := ioutil.ReadFile(fileName)

  if err != nil {
    log.Fatal(err)
    fmt.Printf("Welp")
  }

  for i := 0; i < len(buffer); i++ {
    c8.memory[512 + i] = buffer[i]
  }
  
}