package chip8

import (
	"testing"
)

func TestDissassemble(t *testing.T) {
	c8 := newCpu()
	c8.loadFile("Fishie.ch8")
	for i := 0; i < 10; i++ {
		c8.emulateOneCycle()
	}
}

func TestInit(t *testing.T) {
	c8 := newCpu()
	if c8.pc != 0x200 {
		t.Errorf("Expected pc of 0x200, got %X instead", c8.pc)
	}
}

func TestInstruction1(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x1FFF)
	c8.executeInstruction(inst)
	if c8.pc != 0x0FFF {
		t.Errorf("Expected pc of 0x0FFF, got %X instead", c8.pc)
	}
}

func TestInstruction2(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x20AA)
	c8.executeInstruction(inst)

	if c8.sp != 0x1 {
		t.Errorf("Expected sp of 0x1 , got %X instead", c8.sp)
	}
	if c8.stack[c8.sp - 1] != 0x200 {
		t.Errorf("Expected sp of 0x200 , got %X instead", c8.stack[c8.sp - 1])
	}
	if c8.pc != 0x00AA {
		t.Errorf("Expected pc of 0x00AA , got %X instead", c8.pc)
	}
}

func TestInstruction3(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x30AA)
	c8.reg[0] = 0xAA
	c8.executeInstruction(inst)
	if c8.pc != 0x204 {
		t.Errorf("Expected pc of 0x204 , got 0x%04X instead", c8.pc)
	}

	c8 = newCpu()
	inst = uint16(0x30AA)
	c8.reg[0] = 0xAB
	c8.executeInstruction(inst)
	if c8.pc != 0x202 {
		t.Errorf("Expected pc of 0x202 , got 0x%04X instead", c8.pc)
	}
}

func TestInstruction4(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x40AA)
	c8.reg[0] = 0xAB
	c8.executeInstruction(inst)

	if c8.pc != 0x0204 {
		t.Errorf("Expected pc of 0x204 , got 0x%04X instead", c8.pc)
	}

	c8 = newCpu()
	inst = uint16(0x40AA)
	c8.reg[0] = 0xAA
	c8.executeInstruction(inst)

	if c8.pc != 0x0202 {
		t.Errorf("Expected pc of 0x202 , got 0x%04X instead", c8.pc)
	}
}

func TestInstruction5(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x50A0)
	c8.reg[0] = 0xAA
	c8.reg[10] = 0xAA
	c8.executeInstruction(inst)

	if c8.pc != 0x0204 {
		t.Errorf("Expected pc of 0x204 , got 0x%04X instead", c8.pc)
	}

	c8 = newCpu()
	inst = uint16(0x50A0)
	c8.reg[0] = 0xAA
	c8.reg[10] = 0xAB
	c8.executeInstruction(inst)

	if c8.pc != 0x0202 {
		t.Errorf("Expected pc of 0x202 , got 0x%04X instead", c8.pc)
	}
}

func TestInstruction6(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x6DAA)
	c8.executeInstruction(inst)

	if c8.reg[0xD] != 0x00AA {
		t.Errorf("Expected value of 0x00AA , got 0x%04X instead", c8.pc)
	}
}

func TestInstruction7(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x7401)
	c8.executeInstruction(inst)
	if c8.reg[4] != 0x01 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.reg[4])
	}
	c8.executeInstruction(inst)
	if c8.reg[4] != 0x02 {
		t.Errorf("Expected value of 0x02 , got 0x%02X instead", c8.reg[4])
	}
}

func TestInstruction8(t *testing.T) {
	// 8XY0
	c8 := newCpu()
	inst := uint16(0x80A0)
	c8.reg[10] = 0xAD
	c8.executeInstruction(inst)
	if c8.reg[0] != 0xAD{
		t.Errorf("Expected value of 0xAD , got 0x%02X instead", c8.reg[0])
	}

	// 8XY1
	c8 = newCpu()
	inst = uint16(0x80A1)
	c8.reg[0] = 0xAA
	c8.reg[0xA] = 0x55
	c8.executeInstruction(inst)
	if c8.reg[0] != 0xFF{
		t.Errorf("Expected value of 0xFF , got 0x%02X instead", c8.reg[0])
	}

	// 8XY2
	c8 = newCpu()
	inst = uint16(0x80A2)
	c8.reg[0] = 0xAA
	c8.reg[0xA] = 0x55
	c8.executeInstruction(inst)
	if c8.reg[0] != (0x00){
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[0])
	}

	// 8XY3
	c8 = newCpu()
	inst = uint16(0x80A3)
	c8.reg[0] = 0xAA
	c8.reg[0xA] = 0x55
	c8.executeInstruction(inst)
	if c8.reg[0] != (0xAA ^ 0x55){
		t.Errorf("Expected value of 0xFF , got 0x%02X instead", c8.reg[0])
	}

	// 8XY4
	c8 = newCpu()
	inst = uint16(0x80A4)
	c8.reg[0] = 0xAA
	c8.reg[0xA] = 0x55
	c8.executeInstruction(inst)
	if c8.reg[0] != (0xAA + 0x55){
		t.Errorf("Expected value of 0xFF , got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[0xF] != 0 {
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[0xF])
	}
	inst = uint16(0x81B4)
	c8.reg[0x1] = 0xFF
	c8.reg[0xB] = 0xFF
	c8.executeInstruction(inst)
	if c8.reg[1] != (0xFE){
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[1])
	}
	if c8.reg[0xF] != 1 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.reg[0xF])
	}

	// 8XY5
	c8 = newCpu()
	inst = uint16(0x80A5)
	c8.reg[0] = 0xFF
	c8.reg[0xA] = 0xAA
	c8.executeInstruction(inst)
	if c8.reg[0] != (0xFF - 0xAA){
		t.Errorf("Expected value of 0x55 , got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[0xF] != 1 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.reg[0xF])
	}
	inst = uint16(0x81B5)
	c8.reg[0x1] = 0x00
	c8.reg[0xB] = 0x01
	c8.executeInstruction(inst)
	if c8.reg[1] != (0xFF){
		t.Errorf("Expected value of 0xFF , got 0x%02X instead", c8.reg[1])
	}
	if c8.reg[0xF] != 0 {
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[0xF])
	}

	// 8XY6
	c8 = newCpu()
	inst = uint16(0x80A6)
	c8.reg[0xA] = 0xFF
	c8.executeInstruction(inst)
	if c8.reg[0] != (0xFF >> 1){
		t.Errorf("Expected value of 0x7F , got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[0xF] != 1 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.reg[0xF])
	}
	inst = uint16(0x80A6)
	c8.reg[0xA] = 0xF0
	c8.executeInstruction(inst)
	if c8.reg[0] != (0xF0 >> 1){
		t.Errorf("Expected value of 0x7F , got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[0xF] != 0 {
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[0xF])
	}

	// 8XY7
	c8 = newCpu()
	inst = uint16(0x80A7)
	c8.reg[0x0] = 0xFF
	c8.reg[0xA] = 0xFF
	c8.executeInstruction(inst)
	if c8.reg[0] != (0x00){
		t.Errorf("Expected value of 0x00, got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[0xF] != 1 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.reg[0xF])
	}
	inst = uint16(0x80A7)
	c8.reg[0x0] = 0x01
	c8.reg[0xA] = 0x00
	c8.executeInstruction(inst)
	if c8.reg[0] != (0xFF){
		t.Errorf("Expected value of 0xFF, got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[0xF] != 0 {
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[0xF])
	}

	// 8XYE
	c8 = newCpu()
	inst = uint16(0x80AE)
	c8.reg[0xA] = 0xFF
	c8.executeInstruction(inst)
	if c8.reg[0] != (0xFE){
		t.Errorf("Expected value of 0xFE , got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[0xF] != 1 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.reg[0xF])
	}
}

// 9XY0
func TestInstruction9(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x90A0)
	c8.reg[0] = 0x0A
	c8.reg[0xA] = 0x0A
	c8.executeInstruction(inst)
	if c8.pc != 0x0202 {
		t.Errorf("Expected pc of 0x202 , got 0x%04X instead", c8.pc)
	}
	inst = uint16(0x90A0)
	c8.reg[0] = 0x0A
	c8.reg[0xA] = 0xDD
	c8.executeInstruction(inst)

	if c8.pc != 0x0206 {
		t.Errorf("Expected pc of 0x204 , got 0x%04X instead", c8.pc)
	}
}

// ANNN
func TestInstructionA(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0xABAA)
	c8.executeInstruction(inst)

	if c8.i != 0x0BAA {
		t.Errorf("Expected i of 0xBAA , got 0x%04X instead", c8.i)
	}
}

// BNNN
func TestInstructionB(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0xBAAA)
	c8.executeInstruction(inst)
	if c8.pc != 0xAAA {
		t.Errorf("Expected pc of 0xAAA , got 0x%04X instead", c8.pc)
	}
	c8.reg[0] = 0xAA
	c8.executeInstruction(inst)
	if c8.pc != 0xB54 {
		t.Errorf("Expected pc of 0xB54 , got 0x%04X instead", c8.pc)
	}
}

// CXNN
func TestInstructionC(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0xC000)
	c8.executeInstruction(inst)

	if c8.reg[0] != 0x00 {
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[0])
	}
}

func TestInstructionE(t *testing.T) {
	// EX9E
	c8 := newCpu()
	inst := uint16(0xED9E)
	c8.executeInstruction(inst)

	if c8.pc != 0x0202 {
		t.Errorf("Expected pc of 0x202 , got 0x%04X instead", c8.pc)
	}
}

func TestInstructionF(t *testing.T) {
	// FX07
	c8 := newCpu()
	inst := uint16(0xF207)
	c8.timerDelay = 0x0D
	c8.executeInstruction(inst)
	if c8.reg[2] != 0x0D {
		t.Errorf("Expected value of 0x0D , got 0x%02X instead", c8.reg[2])
	}

	// FX15
	c8.reg[3] = 0x12
	inst = uint16(0xF315)
	c8.executeInstruction(inst)
	if c8.timerDelay != 0x12 {
		t.Errorf("Expected value of 0x15 , got 0x%02X instead", c8.timerDelay)
	}

	// FX18
	c8.reg[3] = 0x12
	inst = uint16(0xF318)
	c8.executeInstruction(inst)
	if c8.soundDelay != 0x12 {
		t.Errorf("Expected value of 0x15 , got 0x%02X instead", c8.soundDelay)
	}

	// FX1E
	c8.reg[3] = 0x12
	c8.i = 0x00F1
	inst = uint16(0xF31E)
	c8.executeInstruction(inst)
	if c8.i != 0x0012 + 0x00F1 {
		t.Errorf("Expected value of 0x0103 , got 0x%04X instead", c8.i)
	}

	// FX29
	c8.reg[3] = 0x00
	inst = uint16(0xF329)
	c8.executeInstruction(inst)
	if c8.i != 0x0000 {
		t.Errorf("Expected value of 0x0000 , got 0x%04X instead", c8.i)
	}

	c8.reg[3] = 0x01
	inst = uint16(0xF329)
	c8.executeInstruction(inst)
	if c8.i != 0x0005 {
		t.Errorf("Expected value of 0x0000 , got 0x%04X instead", c8.i)
	}

	// FX33
	c8.reg[4] = 0x7B
	inst = uint16(0xF433)
	c8.i = 0x0400
	c8.executeInstruction(inst)
	if c8.memory[c8.i] != 0x0001 {
		t.Errorf("Expected value of 0x0001 , got 0x%04X instead", c8.i)
	}
	if c8.memory[c8.i + 1] != 0x0002 {
		t.Errorf("Expected value of 0x0002 , got 0x%04X instead", c8.i)
	}
	if c8.memory[c8.i + 2] != 0x0003 {
		t.Errorf("Expected value of 0x0003 , got 0x%04X instead", c8.i)
	}

	// FX55
	c8 = newCpu()
	c8.reg[0] = 0x00
	c8.reg[1] = 0x01
	c8.reg[2] = 0x02
	c8.reg[3] = 0x03
	c8.i = 0x0123
	inst = uint16(0xF355)
	c8.executeInstruction(inst)
	if c8.memory[0x0123] != 0x00 {
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.memory[0x0123])
	}
	if c8.memory[0x0124] != 0x01 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.memory[0x0124])
	}
	if c8.memory[0x0125] != 0x02 {
		t.Errorf("Expected value of 0x02 , got 0x%02X instead", c8.memory[0x0125])
	}
	if c8.memory[0x0126] != 0x03 {
		t.Errorf("Expected value of 0x03 , got 0x%02X instead", c8.memory[0x0126])
	}
	// FX65
	c8.reg[0] = 0x00
	c8.reg[1] = 0x00
	c8.reg[2] = 0x00
	c8.reg[3] = 0x00
	c8.i = 0x0123
	inst = uint16(0xF365)
	c8.executeInstruction(inst)
	if c8.reg[0] != 0x00 {
		t.Errorf("Expected value of 0x00 , got 0x%02X instead", c8.reg[0])
	}
	if c8.reg[1] != 0x01 {
		t.Errorf("Expected value of 0x01 , got 0x%02X instead", c8.reg[1])
	}
	if c8.reg[2] != 0x02 {
		t.Errorf("Expected value of 0x02 , got 0x%02X instead", c8.reg[2])
	}
	if c8.reg[3] != 0x03 {
		t.Errorf("Expected value of 0x03 , got 0x%02X instead", c8.reg[3])
	}
}

func TestInstructionX(t *testing.T) {
	c8 := newCpu()
	inst := uint16(0x20AA)
	c8.executeInstruction(inst)

	if c8.pc != 0x0202 {
		//t.Errorf("Expected pc of 0x202 , got 0x%04X instead", c8.pc)
	}
}
