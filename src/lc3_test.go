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

/* Tests for operations */
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

func TestLDI(t *testing.T) {

	reg[R_R1] = 0xffff
	reg[R_PC] = 0x3000
	addressToStore := uint16(0x5000)
	expectedValue := uint16(42)

	mem[reg[R_PC]+2] = addressToStore
	mem[addressToStore] = expectedValue

	/* ldi*/
	/*1010    001    000000001*/
	/*OP      DES    PCOffset9 =1*/

	code := uint16(0xA200)

	//positive offset
	for pcOffset9 := uint16(0x0000); pcOffset9 <= 0x00FF; pcOffset9++ {
		reg[R_R1] = 0xffff
		reg[R_PC] = 0x3000
		code = uint16(0xA200) | uint16(pcOffset9)
		mem[reg[R_PC]+1+pcOffset9] = addressToStore
		mem[addressToStore] = expectedValue

		ldi(code)
		if reg[R_R1] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R1], expectedValue)
		}
	}

	//negative offset
	for pcOffset9 := uint16(0x0100); pcOffset9 < 0x001FF; pcOffset9++ {
		reg[R_R1] = 0xffff
		reg[R_PC] = 0x3000
		code = uint16(0xA200) | (uint16(pcOffset9))
		mediate_address := reg[R_PC] + 1 + (pcOffset9 | 0xFE00) //SignExtend negative
		mem[mediate_address] = addressToStore
		mem[addressToStore] = expectedValue

		ldi(code)
		if reg[R_R1] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R1], expectedValue)
		}
	}
}

func TestLD(t *testing.T) {

	reg[R_R1] = 0xffff
	reg[R_PC] = 0x3000
	expectedValue := uint16(42)

	/* ld*/
	/*0010    001    000000001*/
	/*OP      DES    PCOffset9 = 1*/

	code := uint16(0xA200)

	//positive offset
	for pcOffset9 := uint16(0x0000); pcOffset9 <= 0x00FF; pcOffset9++ {
		reg[R_R1] = 0xffff
		reg[R_PC] = 0x3000
		code = uint16(0xA200) | uint16(pcOffset9)
		mem[reg[R_PC]+1+pcOffset9] = expectedValue

		ld(code)
		if reg[R_R1] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R1], expectedValue)
		}
	}

	//negative offset
	for pcOffset9 := uint16(0x0100); pcOffset9 < 0x001FF; pcOffset9++ {
		reg[R_R1] = 0xffff
		reg[R_PC] = 0x3000
		code = uint16(0xA200) | (uint16(pcOffset9))
		mediate_address := reg[R_PC] + 1 + (pcOffset9 | 0xFE00) //SignExtend negative
		mem[mediate_address] = expectedValue

		ld(code)
		if reg[R_R1] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R1], expectedValue)
		}
	}
}

func TestLDR(t *testing.T) {

	expectedValue := uint16(42)

	/* ldr*/
	/*0110    001 010    000001*/
	/*OP      DES BaseR  offset*/

	code := uint16(0x6280)

	//positive offset
	for pcOffset6 := uint16(0x0000); pcOffset6 <= 0x001F; pcOffset6++ {
		reg[R_R1] = 0xffff
		reg[R_R2] = 0x3000

		code = uint16(0x6280) | uint16(pcOffset6)
		mem[reg[R_R2]+pcOffset6] = expectedValue

		ldr(code)
		if reg[R_R1] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R1], expectedValue)
		}
	}

	//negative offset
	for pcOffset6 := uint16(0x0021); pcOffset6 <= 0x003F; pcOffset6++ {
		reg[R_R1] = 0xffff
		reg[R_R2] = 0x3000

		code = uint16(0x6280) | uint16(pcOffset6)
		mem[reg[R_R2]+(pcOffset6|0xFFC0)] = expectedValue

		ldr(code)
		if reg[R_R1] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R1], expectedValue)
		}
	}
}

func TestLEA(t *testing.T) {

	reg[R_R7] = 0xffff
	reg[R_PC] = 0x3000

	/* lea*/
	/*1110    111    000000001*/
	/*OP      DES    PCOffset9 = 1*/

	code := uint16(0xEE00)
	expectedValue := uint16(0)

	//positive offset
	for pcOffset9 := uint16(0x0000); pcOffset9 <= 0x00FF; pcOffset9++ {
		reg[R_R7] = 0xffff
		reg[R_PC] = 0x8211
		code = uint16(0xEE00) | uint16(pcOffset9)
		expectedValue = reg[R_PC] + 1 + pcOffset9

		lea(code)
		if reg[R_R7] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R7], expectedValue)
		}
	}

	//negative offset
	for pcOffset9 := uint16(0x0100); pcOffset9 < 0x001FF; pcOffset9++ {
		reg[R_R7] = 0xffff
		reg[R_PC] = 0x8211
		code = uint16(0xEE00) | (uint16(pcOffset9))
		expectedValue = reg[R_PC] + 1 + (pcOffset9 | 0xFE00) //SignExtend negative
		lea(code)
		if reg[R_R7] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", reg[R_R7], expectedValue)
		}
	}
}

func TestAndRegisters(t *testing.T) {

	/* logic and R1 and R2*/
	/*0101    000    001    0    00     010*/
	/*OP      DES    SRC    Mode Unused SRC2*/
	var expected uint16

	for i := uint16(0); i < 0xffff; i++ {
		for j := uint16(0); j < 0xffff; j++ {
			reg[0] = 0x0000
			reg[1] = i
			reg[2] = j
			expected = i & j
			and(uint16(0x5042))
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
			expected = i & j
			and(uint16(0x5F07))
			if reg[7] != expected {
				t.Errorf("Failed ! got %x want %x", reg[7], expected)
			}
		}
	}
	//Check Condition code
}

func TestAndImmediate(t *testing.T) {
	reg[1] = 0xffff
	reg[2] = 0xffff

	/* add R1 and R2*/
	/*0001    000    001    0    00     010*/
	/*OP      DES    SRC    Mode Unused SRC2*/
	var expected uint16
	code := uint16(0x5060)

	//positive second operand
	for i := uint16(0); i < 0xffff; i++ {
		for j := uint16(0); j < 0x000f; j++ {

			reg[0] = 0x0000
			reg[1] = i
			expected = i & j
			and(uint16(code | j))
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
			expected = i & (j | uint16(mask))
			and(uint16(code | j))
			if reg[0] != expected {
				t.Errorf("Failed ! got %x want %x", reg[0], expected)
			}
		}
	}
}

func TestNot(t *testing.T) {

	/* logic not:  R1 = not R2*/
	/*1001    000    001    111111*/
	/*OP      DES    SRC    Unused*/
	var expected uint16

	for i := uint16(0); i < 0xffff; i++ {
		reg[0] = 0x0000
		reg[1] = i
		expected = ^i

		not(uint16(0x907F))
		if reg[0] != expected {
			t.Errorf("Failed ! got %x want %x", reg[0], expected)
		}
	}
	//Check Condition code
}

func TestBR(t *testing.T) {

	reg[R_PC] = 0x3000

	/* br*/
	/*0000    000    000000001*/
	/*OP      CC     PCOffset9 =1*/

	code := uint16(0x0E01)
	var expectedAddress uint16
	//positive offset
	for pcOffset9 := uint16(0x0000); pcOffset9 <= 0x00FF; pcOffset9++ {
		reg[R_PC] = 0x3000
		reg[R_COND] = 0x0001
		code = uint16(0x0E00) | uint16(pcOffset9)
		expectedAddress = reg[R_PC] + 1 + pcOffset9
		br(code)
		if reg[R_PC] != uint16(expectedAddress) {
			t.Errorf("Failed ! got %x want %x", reg[R_COND], expectedAddress)
		}
	}

	//negative offset
	for pcOffset9 := uint16(0x0100); pcOffset9 < 0x001FF; pcOffset9++ {
		reg[R_PC] = 0x3000
		reg[R_COND] = 0x0001

		code = uint16(0x0E00) | (uint16(pcOffset9))

		expectedAddress = reg[R_PC] + 1 + (pcOffset9 | 0xFE00)
		br(code)
		if reg[R_PC] != uint16(expectedAddress) {
			t.Errorf("Failed ! got %x want %x", reg[R_COND], expectedAddress)
		}
	}
}

func TestST(t *testing.T) {

	reg[R_R1] = 0xffff
	reg[R_PC] = 0x3000
	expectedValue := uint16(42)

	/* ld*/
	/*0011    001    000000001*/
	/*OP      DES    PCOffset9 = 1*/

	code := uint16(0x3200)
	var address uint16
	//positive offset
	for pcOffset9 := uint16(0x0000); pcOffset9 <= 0x00FF; pcOffset9++ {
		reg[R_R1] = expectedValue
		reg[R_PC] = 0x3000
		code = uint16(0x3200) | uint16(pcOffset9)

		address = reg[R_PC] + 1 + pcOffset9
		mem[address] = 0xffff

		st(code)
		if mem[address] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", mem[address], expectedValue)
		}
	}

	//negative offset
	for pcOffset9 := uint16(0x0100); pcOffset9 < 0x001FF; pcOffset9++ {
		reg[R_R1] = expectedValue
		reg[R_PC] = 0x3000
		code = uint16(0x3200) | (uint16(pcOffset9))
		address = reg[R_PC] + 1 + (pcOffset9 | 0xFE00) //SignExtend negative
		mem[address] = 0xffff

		st(code)
		if mem[address] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", mem[address], expectedValue)
		}
	}
}

func TestSTI(t *testing.T) {

	reg[R_R1] = 0xffff
	reg[R_PC] = 0x3000
	addressToStore := uint16(0x5000)
	expectedValue := uint16(42)

	mem[reg[R_PC]+2] = addressToStore
	mem[addressToStore] = expectedValue

	/* ldi*/
	/*1011    001    000000001*/
	/*OP      DES    PCOffset9 =1*/

	code := uint16(0xB200)

	//positive offset
	for pcOffset9 := uint16(0x0000); pcOffset9 <= 0x00FF; pcOffset9++ {
		reg[R_R1] = expectedValue
		reg[R_PC] = 0x3000
		code = uint16(0xB200) | uint16(pcOffset9)
		mem[reg[R_PC]+1+pcOffset9] = addressToStore
		mem[addressToStore] = 0xffff

		sti(code)
		if mem[addressToStore] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", mem[addressToStore], expectedValue)
		}
	}

	//negative offset
	for pcOffset9 := uint16(0x0100); pcOffset9 < 0x001FF; pcOffset9++ {
		reg[R_R1] = expectedValue
		reg[R_PC] = 0x3000
		code = uint16(0xB200) | (uint16(pcOffset9))
		address := reg[R_PC] + 1 + (pcOffset9 | 0xFE00) //SignExtend negative
		mem[address] = addressToStore
		mem[addressToStore] = 0xffff

		sti(code)
		if mem[addressToStore] != uint16(expectedValue) {
			t.Errorf("Failed ! got %x want %x", mem[addressToStore], expectedValue)
		}
	}
}
