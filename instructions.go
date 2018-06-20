package chip8

import (
	"fmt"
)

// New type to represent a chip8 instruction
type instruction uint16

// Retruns a human readable description of a chip8 instruction
func (inst instruction) String() string {
	// Formats instruction as two words, separated by a space
	preamble := fmt.Sprintf("%04X ", uint16(inst))

	switch inst & 0xF000 {
	case 0x0000:
		switch inst & 0x0FFF {
		case 0x00E0:
			// Clear screen
			return preamble + "CLS"
		case 0x00EE:
			// Return from subroutine
			return preamble + "RTS"
		default:
			return "Instruction not recognized"
		}
	case 0x1000:
		// JUMP to instruction at addresss 0x0NNN
		imm := inst & 0x0FFF
		return preamble + fmt.Sprintf("JUMP 0x%04X", uint16(imm))
	case 0x2000:
		// CALL subroutine at address 0xNNN
		imm := inst & 0x0FFF
		return preamble + fmt.Sprintf("CALL 0x%04X", uint16(imm))
	case 0x3000:
		// SKIP next instruction if VX == imm
		reg := inst >> 8 & 0x0F
		imm := inst & 0x00FF
		return preamble + fmt.Sprintf("SKIP.EQ V%X, %02X", uint16(reg), uint16(imm))
	case 0x4000:
		// SKIP next instruction if VX != imm
		reg := inst >> 8 & 0x0F
		imm := inst & 0x00FF
		return preamble + fmt.Sprintf("SKIP.NEQ V%X, %02X", uint16(reg), uint16(imm))
	case 0x5000:
    if inst & 0x000F == 0x0 {
		// SKIP next instruction if VX == VY
		regX := inst >> 8 & 0x0F
		regY := inst >> 4 & 0x00F
		return preamble + fmt.Sprintf("SKIP.EQ V%X, V%X", uint16(regX), uint16(regY))
    } else {
      return "Instruction not recognized"
    }
	case 0x6000:
		// MOVE immediate into register VX
		reg := inst >> 8 & 0x0F
		imm := inst & 0x00FF
		return preamble + fmt.Sprintf("MVI V%X, %02X", uint16(reg), uint16(imm))
	case 0x7000:
		// ADD imm to value at VX
		reg := inst >> 8 & 0x0F
		imm := inst & 0x00FF
		return preamble + fmt.Sprintf("ADD V%X, %02X", uint16(reg), uint16(imm))
	case 0x8000:
		// Operates with two registers
		regX := inst >> 8 & 0x0F
		regY := inst >> 4 & 0x00F
		switch inst & 0x000F {
		case 0x0:
      // Set VX to the value of VY
		  return preamble + fmt.Sprintf("MOV V%X, V%X", uint16(regX), uint16(regY))
		case 0x1:
      // Set VX to the value of VY | VX
		  return preamble + fmt.Sprintf("OR V%X, V%X", uint16(regX), uint16(regY))
		case 0x2:
      // Set VX to the value of VY & VX
		  return preamble + fmt.Sprintf("AND V%X, V%X", uint16(regX), uint16(regY))
		case 0x3:
      // Set VX to the value of VY ^ VX
		  return preamble + fmt.Sprintf("XOR V%X, V%X", uint16(regX), uint16(regY))
		case 0x4:
      // Set VX to the value of VY + VX, VF set to 1 if there is a carry over
		  return preamble + fmt.Sprintf("ADD. V%X, V%X", uint16(regX), uint16(regY))
		case 0x5:
      // Set VX to the value of VX - VY, VF set to 1 if need to borrow
		  return preamble + fmt.Sprintf("SUB. V%X, V%X", uint16(regX), uint16(regY))
		case 0x6:
      // Set VX to the value of VY >> 1, VF set to least significant digit of VY
		  return preamble + fmt.Sprintf("SHR. V%X, V%X", uint16(regX), uint16(regY))
		case 0x7:
      // Set VX to the value of VY - VX, VH set to 1 if need to borrow
		  return preamble + fmt.Sprintf("SUBB. V%X, V%X", uint16(regX), uint16(regY))
		case 0xE:
      // Set VX and VY to VY << 1, VF set to most significant digit of VY before shift.
		  return preamble + fmt.Sprintf("SHL. V%X, V%X", uint16(regX), uint16(regY))
		}
	case 0x9000:
    if inst & 0x000F == 0x0 {
		// SKIP next instruction if VX != VY
		regX := inst >> 8 & 0x0F
		regY := inst >> 4 & 0x00F
		return preamble + fmt.Sprintf("SKIP.NEQ V%X, V%X", uint16(regX), uint16(regY))
    } else {
      return "Instruction not recognized"
    }
	case 0xA000:
    // Set register I to imm
		imm := inst & 0x0fff
		return preamble + fmt.Sprintf("MVI I 0x%04X", uint16(imm))
	case 0xB000:
		// JUMP to address at imm + value at V0
		imm := inst & 0x0fff
		return preamble + fmt.Sprintf("JUMP 0x%04X(V0)", uint16(imm))
	case 0xC000:
    // Set register VX to Imm & rand(0,255)
		reg := inst >> 8 & 0x0F
		imm := inst & 0x0fff
		return preamble + fmt.Sprintf("RAND V%X 0x%04X", uint16(reg), uint16(imm))
	case 0xD000:
		regX := inst >> 8 & 0x0F
		regY := inst >> 4 & 0x00F
    N := inst & 0x000F
		return preamble + fmt.Sprintf("SPRITE V%X, V%X, %X", uint16(regX), uint16(regY), uint16(N))
	case 0xE000:
    switch inst & 0x00FF {
    case 0x9E:
      // SKIP next instruction if key stored in VX is held
		  reg := inst >> 8 & 0x0F
      return preamble + fmt.Sprintf("SKIP.KEY V%X", uint16(reg))
    case 0xA1:
      // SKIP next instruction if key stored in VX isn't held
		  reg := inst >> 8 & 0x0F
      return preamble + fmt.Sprintf("SKIP.NKEY V%X", uint16(reg))
    default:
      return "Instruction not recognized"
    }
	case 0xF000:
		reg := inst >> 8 & 0x0F
    switch inst & 0x00FF {
    case 0x07:
      // Set VX to value of delay timer
      return preamble + fmt.Sprintf("MOV V%X DELAY", uint16(reg))
    case 0x0A:
      // Wait for keypress, then store in VX
      return preamble + fmt.Sprintf("WAITKEY V%X", uint16(reg))
    case 0x15:
      // Set delay timer to value in VX
      return preamble + fmt.Sprintf("MOV DELAY V%X", uint16(reg))
    case 0x18:
      // Set sound timer to value in VX
      return preamble + fmt.Sprintf("MOV SOUND V%X", uint16(reg))
    case 0x1E:
      // ADDS VX to I
      return preamble + fmt.Sprintf("ADD I V%X", uint16(reg))
    case 0x29:
      // Sets I to the location of sprite of character in VX
      return preamble + fmt.Sprintf("SPRITECHAR I V%X", uint16(reg))
    case 0x33:
      // Stores BCD of VX at I, I+1, I+2
      return preamble + fmt.Sprintf("MOVBCD V%X", uint16(reg))
    case 0x55:
      // Stores V0 to VX in memory starting at I
      return preamble + fmt.Sprintf("STORE (I) V0-V%X", uint16(reg))
    case 0x65:
      // Load values at V0 to VX starting at memory address I
      return preamble + fmt.Sprintf("LOAD V0-V%X (I)", uint16(reg))
    }
	default:
		return "Instruction not recognized"
	}
	return "Instruction not recognized"
}

func main() {
	test_instruction := instruction(0x3A00)

	fmt.Println(test_instruction)
}
