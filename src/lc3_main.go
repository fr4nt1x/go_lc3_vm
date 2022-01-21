package main

import (
	"fmt"
	"math"
	"time"
)

const (
	OP_BR    uint16 = iota /* branch */
	OP_ADD                 /* add  */
	OP_LD                  /* load */
	OP_ST                  /* store */
	OP_JSR                 /* jump register */
	OP_AND                 /* bitwise and */
	OP_LDR                 /* load register */
	OP_STR                 /* store register */
	OP_RTI                 /* unused */
	OP_NOT                 /* bitwise not */
	OP_LDI                 /* load indirect */
	OP_STI                 /* store indirect */
	OP_JMP                 /* jump */
	OP_RES                 /* reserved (unused) */
	OP_LEA                 /* load effective address */
	OP_TRAP                /* execute trap */
	OP_COUNT               /*count of ops*/
)

const (
	R_R0 uint16 = iota
	R_R1
	R_R2
	R_R3
	R_R4
	R_R5
	R_R6
	R_R7
	R_PC /* program counter */
	R_COND
	R_COUNT
)

const (
	FL_POS = 1 << 0 /* P */
	FL_ZRO = 1 << 1 /* Z */
	FL_NEG = 1 << 2 /* N */
)

var mem [math.MaxUint16]uint16
var reg [R_COUNT]uint16

//helper
func getSignedFromUnsignedSlice(unsigned []uint16) []int16 {
	res := make([]int16, len(unsigned))
	for i, v := range unsigned {
		res[i] = getSignedFromUnsigned(v)
	}
	return res
}

func getSignedFromUnsigned(unsigned uint16) int16 {

	res := int16(unsigned & 0x7FFF)
	if unsigned>>15 != 0 {
		unsigned = unsigned - 1
		unsigned = ^unsigned
		res = -1 * int16(unsigned&0x7FFF)
	}
	return res
}

func isBitnSet(input uint16, n uint8) bool {
	//n counts from zero -> 0 = first bit, 8 = 9nth bit
	mask := uint16(0x0001) << n
	return (input & mask) != 0
}

func signExtend(x uint16, nBits uint8) uint16 {
	if isBitnSet(x, nBits-1) {
		mask := (uint32(0x0000FFFF) << nBits) & uint32(0x0000FFFF)
		x |= uint16(mask) // set all bits except last nBits to ones
	}
	return x
}

func update_flags(r uint16) {
	if reg[r] == 0 {
		reg[R_COND] = FL_ZRO
	} else if reg[r]>>15 != 0 {
		reg[R_COND] = FL_NEG
	} else {
		reg[R_COND] = FL_POS
	}
}

//instructions
func br(instr uint16) {
}

func add(instr uint16) {
	destReg := (instr >> 9) & 0x7
	src1 := (instr >> 6) & 0x7

	if !isBitnSet(instr, 5) {
		//add1
		src2 := instr & 0x7
		reg[destReg] = reg[src1] + reg[src2]
	} else {
		//add2
		value := signExtend(instr&0x001F, 5)
		reg[destReg] = reg[src1] + value
	}
	update_flags(destReg)
}

func ld(instr uint16) {

}
func st(instr uint16)   {}
func jsr(instr uint16)  {}
func and(instr uint16)  {}
func ldr(instr uint16)  {}
func str(instr uint16)  {}
func rti(instr uint16)  {}
func not(instr uint16)  {}
func ldi(instr uint16)  {}
func sti(instr uint16)  {}
func jmp(instr uint16)  {}
func res(instr uint16)  {}
func lea(instr uint16)  {}
func trap(instr uint16) {}

//array of op functions
var instr_funcs [OP_COUNT]func(uint16)

//memory read
func mr(address uint16) uint16 {
	return mem[address]
}

//memory write
func mw(address uint16, val uint16) {
	mem[address] = val
}

func main() {
	PC_START := uint16(0x3000)
	mem = [math.MaxUint16]uint16{0}
	reg[R_PC] = PC_START
	running := true
	debug := true
	instr_funcs = [OP_COUNT]func(uint16){br, add, ld, st, jsr, and, ldr, str, rti, not, ldi, sti, jmp, res, lea, trap}
	testinstr := uint16(0x1042) //add r0 = r1 +r2
	mw(PC_START, testinstr)
	reg[1] = 0xffff
	reg[2] = 0xffff
	for running {
		instruction := mr(reg[R_PC])
		op_code := instruction >> 12

		if debug {
			fmt.Println(getSignedFromUnsignedSlice(reg[:]))
			time.Sleep(500 * time.Millisecond)

			fmt.Println("PCOUNT:", reg[R_PC])
			fmt.Println("Instruction:", instruction)
		}
		instr_funcs[op_code](instruction)
		reg[R_PC]++
	}

	fmt.Println(mr(16), PC_START)
}
