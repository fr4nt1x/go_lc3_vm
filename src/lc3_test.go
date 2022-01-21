package main

import (
	"testing"
)

func TestGetSignedFromUnsignedSlice(t *testing.T) {
	input := []uint16{0xffff,
		0x0001}
	expected := []int16{-1, 1}

	actual := getSignedFromUnsignedSlice(input)

	for i, v := range expected {
		if actual[i] != v {
			t.Errorf("Failed ! got %x want %x", actual[i], v)
		}
	}

}
func TestGetSignedFromUnsigned(t *testing.T) {
	input := []uint16{0xffff,
		0x0001}
	expected := []int16{-1, 1}

	var actual int16

	for i, v := range input {
		actual = getSignedFromUnsigned(v)
		if actual != expected[i] {
			t.Errorf("Failed ! got %x want %x", actual, expected[i])
		}
	}
}

func TestAddRegisters(t *testing.T) {
	reg[0] = 0x0000
	reg[1] = 0xffff
	reg[2] = 0xffff

	/* add R1 and R2*/
	/*0001    000    001    0    00     010*/
	/*OP      DES    SRC    Mode Unused SRC2*/
	var expected uint16

	for i := uint16(0); i < 0xffff; i++ {
		for j := uint16(0); j < 0xffff; j++ {
			reg[0] = 0x0000
			reg[1] = i
			reg[2] = j
			expected = i + j
			add(uint16(0x1042))
			if reg[0] != expected {
				t.Errorf("Failed ! got %x want %x", reg[0], expected)
			}
		}
	}
	for i := uint16(0); i < 0xff; i++ {
		for j := uint16(0); j < 0xff; j++ {
			reg[0] = 0x0000
			reg[4] = i
			reg[7] = j
			expected = i + j
			add(uint16(0x1F07))
			if reg[7] != expected {
				t.Errorf("Failed ! got %x want %x", reg[7], expected)
			}
		}
	}
	//Check Condition code
}

func TestAddImmediate(t *testing.T) {
	reg[1] = 0xffff
	reg[2] = 0xffff

	/* add R1 and R2*/
	/*0001    000    001    0    00     010*/
	/*OP      DES    SRC    Mode Unused SRC2*/
	var expected uint16
	code := uint16(0x1060)

	//positive second operand
	for i := uint16(0); i < 0xffff; i++ {
		for j := uint16(0); j < 0x000f; j++ {

			reg[0] = 0x0000
			reg[1] = i
			expected = i + j
			add(uint16(code | j))
			if reg[0] != expected {
				t.Errorf("Failed ! got %x want %x", reg[0], expected)
			}
		}
	}
	//negative second operand
	mask := (uint32(0x0000FFFF) << 5) & uint32(0x0000FFFF)
	for i := uint16(0); i < 0xffff; i++ {
		for j := uint16(0x0011); j < 0x0020; j++ {
			//Everything is positive
			reg[0] = 0x0000
			reg[1] = i
			expected = i + (j | uint16(mask))
			add(uint16(code | j))
			if reg[0] != expected {
				t.Errorf("Failed ! got %x want %x", reg[0], expected)
			}
		}
	}
}

func TestSignExtend5(t *testing.T) {
	input := []uint16{0x001F, //-1
		0x0011, //-15
		0x0004} //4

	expected := []uint16{0xFFFF,
		0xFFF1,
		0x0004}
	var actual uint16
	for i, v := range input {
		actual = signExtend(v, 5)
		if actual != expected[i] {
			t.Errorf("Failed ! got %x want %x", actual, expected[i])
		}
	}
}

func TestSignExtend9(t *testing.T) {
	input := []uint16{0x01ff, //-1
		0x0101, //-255
		0x0004} //4

	expected := []uint16{0xFFFF,
		0xFF01,
		0x0004}
	var actual uint16
	for i, v := range input {
		actual = signExtend(v, 9)
		if actual != expected[i] {
			t.Errorf("Failed ! got %x want %x", actual, expected[i])
		}
	}
}
