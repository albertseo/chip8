package chip8

import (
	"fmt"
  "io/ioutil"
  "math/rand"

  "github.com/nsf/termbox-go"
)

type cpu struct {
	// Memory is 4096 bytes
	memory [4096]uint8
  // Graphics is 2048 bits
  graphics Display
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
  // Flag to see if need to flush graphics to screen
  drawFlag bool
}

// Preloaded fonts for the memory starting at 0x000 in the memory
var fontSprite = []uint8{ 
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
  0xF0, 0x80, 0xF0, 0x80, 0x80}

// Returns a new instance of a Chip-8 CPU
func newCpu() *cpu {
  // New returns a pointer to a new zero'ed out cpu struct
	cp8 := new(cpu)
	// Set stack pointer and program counter
	cp8.pc = 0x200
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
      c8.graphics.clear()
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
    c8.stack[c8.sp] = c8.pc - 2
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
      temp := uint16(c8.reg[regY]) + uint16(c8.reg[regX])
      if temp > 0xFF {
        c8.reg[15] = 1
      } else {
        c8.reg[15] = 0
      }
      c8.reg[regX] = uint8(temp)
    case 0x5:
      // Set VX to the value of VX - VY, VF set to 0 if need to borrow
      if c8.reg[regX] >= c8.reg[regY] {
        c8.reg[15] = 1
      } else {
        c8.reg[15] = 0
      }
      c8.reg[regX] = c8.reg[regX] - c8.reg[regY]
    case 0x6:
      // Set VX to the value of VY >> 1, VF set to least significant digit of VY
      c8.reg[15] = c8.reg[regY] & 1
      c8.reg[regX] = c8.reg[regY] >> 1
    case 0x7:
      // Set VX to the value of VY - VX, VH set to 0 if need to borrow
      if c8.reg[regY] >= c8.reg[regX] {
        c8.reg[15] = 1
      } else {
        c8.reg[15] = 0
      }
      c8.reg[regX] = c8.reg[regY] - c8.reg[regX]
    case 0xE:
      // Set VX and VY to VY << 1, VF set to most significant digit of VY before shift.
      c8.reg[15] = c8.reg[regY] >> 7
      c8.reg[regX] = c8.reg[regY] << 1
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
    c8.reg[regX] = uint8(imm) & uint8(rand.Int())
  case 0xD000:
    // Draw stuff to the screen
    xCord := c8.reg[inst >> 8 & 0x0F]
    yCord := c8.reg[inst >> 4 & 0x00F]
    height := inst & 0x000F
    c8.reg[15] = c8.graphics.drawSprite(xCord, yCord, height, c8.memory[c8.i:c8.i+height])
    // Maybe Flush to screen here?
    c8.drawFlag = true
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
    }
  case 0xF000:
    regX := inst >> 8 & 0x0F
    switch inst & 0x00FF {
    case 0x07:
      // Set VX to value of delay timer
      c8.reg[regX] = c8.timerDelay
    case 0x0A:
      // Wait for keypress, then store in VX
      c8.reg[regX] = c8.getKey()
    case 0x15:
      // Set delay timer to value in VX
      c8.timerDelay = c8.reg[regX]
    case 0x18:
      // Set sound timer to value in VX
      c8.soundDelay = c8.reg[regX]
    case 0x1E:
      // ADDS VX to I
      c8.i += uint16(c8.reg[regX])
    case 0x29:
      // Sets I to the location of sprite of character in VX
      c8.i = uint16(c8.reg[regX])
    case 0x33:
      // Stores BCD of VX at I, I+1, I+2
      value := c8.reg[regX]
      c8.memory[c8.i + 2] = uint8(value % 10)
      value = value / 10
      c8.memory[c8.i + 1] = uint8(value % 10)
      value = value / 10
      c8.memory[c8.i] = uint8(value % 10)
    case 0x55:
      // Stores V0 to VX in memory starting at I
      for j := 0; j < int(regX); j++ {
        c8.memory[c8.i + uint16(j)] = c8.reg[j]
      }
    case 0x65:
      // Load values at V0 to VX starting at memory address I
      for j := 0; j < int(regX); j++ {
        c8.reg[j] = c8.memory[c8.i + uint16(j)]
      }
    }
  }
}

func (c8 *cpu) loadFile(fileName string) {
  buffer, err := ioutil.ReadFile(fileName)

  if err != nil {
    fmt.Printf("Welp")
  }

  for i := 0; i < len(buffer); i++ {
    c8.memory[512 + i] = buffer[i]
  }
}

var keyMap = map[rune]byte{
  '1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
  'q': 0x04, 'w': 0x05, 'e': 0x06, 'r': 0x0D,
  'a': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
  'z': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
}

func (c8 *cpu) getKey() (uint8) {
  key := termbox.PollEvent()
  return keyMap[key.Ch]
}

func (c8 *cpu) loadSprites() {
  for i := 0; i < 80; i++ {
    c8.memory[i] = fontSprite[i]
  }
}