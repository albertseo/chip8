package chip8

import (
	"testing"
  "fmt"
)

func TestString(t *testing.T) {
	// 0x0XXX
	inst := instruction(0x00E0)
	result := fmt.Sprintf("%v", inst)
	expected := "00E0 CLS"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x00EE)
	result = fmt.Sprintf("%v", inst)
	expected = "00EE RTS"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x01EE)
	result = fmt.Sprintf("%v", inst)
	expected = "Instruction not recognized"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x1XXX JUMP
	inst = instruction(0x1AAB)
	result = fmt.Sprintf("%v", inst)
	expected = "1AAB JUMP 0x0AAB"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x1001)
	result = fmt.Sprintf("%v", inst)
	expected = "1001 JUMP 0x0001"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x2XXX CALL
	inst = instruction(0x2AFD)
	result = fmt.Sprintf("%v", inst)
	expected = "2AFD CALL 0x0AFD"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x3XXX SKIP.EQ
	inst = instruction(0x3AFD)
	result = fmt.Sprintf("%v", inst)
	expected = "3AFD SKIP.EQ VA, FD"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x4XXX SKIP.NEQ
	inst = instruction(0x4AFD)
	result = fmt.Sprintf("%v", inst)
	expected = "4AFD SKIP.NEQ VA, FD"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x5XXX SKIP.EQ
	inst = instruction(0x5A10)
	result = fmt.Sprintf("%v", inst)
	expected = "5A10 SKIP.EQ VA, V1"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x5A11)
	result = fmt.Sprintf("%v", inst)
	expected = "Instruction not recognized"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x6XXX MOV
	inst = instruction(0x6A10)
	result = fmt.Sprintf("%v", inst)
	expected = "6A10 MVI VA, 10"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x7XXX
	inst = instruction(0x7210)
	result = fmt.Sprintf("%v", inst)
	expected = "7210 ADD V2, 10"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x8XXX
	inst = instruction(0x8210)
	result = fmt.Sprintf("%v", inst)
	expected = "8210 MOV V2, V1"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x85E1)
	result = fmt.Sprintf("%v", inst)
	expected = "85E1 OR V5, VE"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x8422)
	result = fmt.Sprintf("%v", inst)
	expected = "8422 AND V4, V2"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x8163)
	result = fmt.Sprintf("%v", inst)
	expected = "8163 XOR V1, V6"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x8964)
	result = fmt.Sprintf("%v", inst)
	expected = "8964 ADD. V9, V6"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x8F25)
	result = fmt.Sprintf("%v", inst)
	expected = "8F25 SUB. VF, V2"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x8A86)
	result = fmt.Sprintf("%v", inst)
	expected = "8A86 SHR. VA, V8"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x8A77)
	result = fmt.Sprintf("%v", inst)
	expected = "8A77 SUBB. VA, V7"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x82AE)
	result = fmt.Sprintf("%v", inst)
	expected = "82AE SHL. V2, VA"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0x9XXX
	inst = instruction(0x92A0)
	result = fmt.Sprintf("%v", inst)
	expected = "92A0 SKIP.NEQ V2, VA"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}
	inst = instruction(0x92A1)
	result = fmt.Sprintf("%v", inst)
	expected = "Instruction not recognized"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

  // 0xAXXX
	inst = instruction(0xA2A1)
	result = fmt.Sprintf("%v", inst)
	expected = "A2A1 MVI I 0x2A1"
	if result != expected {
		t.Errorf("toString was incorrect, got %v, expected %v", result, expected)
	}

}
