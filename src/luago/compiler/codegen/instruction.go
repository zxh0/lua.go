package codegen

import . "luago/lua/vm"

// type opArg struct {
// 	val  int
// 	kind int
// }

type instruction struct {
	line   int
	opcode int
	a      int // a or ax
	b      int // b or bx or sbx
	c      int // c
}

type insts []instruction

func (self insts) getLineNumTable() []uint32 {
	lineNums := make([]uint32, len(self))
	for i, inst := range self {
		lineNums[i] = uint32(inst.line)
	}
	return lineNums
}

func (self insts) encode() []uint32 {
	insts := make([]uint32, len(self))
	for i, inst := range self {
		insts[i] = encodeInst(inst)
	}
	return insts
}

func encodeInst(i instruction) uint32 {
	if i.opcode == OP_LOADK {
		i.b -= 0x100
		if i.b >= 0x100 {
			panic("todo!") // OP_LOADKX
		}
	}

	opmode := Instruction(i.opcode).OpMode() // todo
	switch opmode {
	case IABC:
		return uint32(i.opcode | i.a<<6 | i.c<<14 | i.b<<23)
	case IABx:
		return uint32(i.opcode | i.a<<6 | i.b<<14)
	case IAsBx:
		return uint32(i.opcode | i.a<<6 | (i.b+HALF_BX)<<14)
	default: // IAx
		return uint32(i.opcode | i.a<<6)
	}
}
