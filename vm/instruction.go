package vm

import "github.com/zxh0/lua.go/api"

const (
	SIZE_OP = 7
	SIZE_A  = 8
	SIZE_B  = 8
	SIZE_C  = 8
	SIZE_Bx = SIZE_C + SIZE_B + 1
	SIZE_Ax = SIZE_Bx + SIZE_A
	SIZE_sJ = SIZE_Bx + SIZE_A
)

const (
	MAXARG_A   = (1 << SIZE_A) - 1
	MAXARG_B   = (1 << SIZE_B) - 1
	MAXARG_C   = (1 << SIZE_C) - 1
	MAXARG_sC  = MAXARG_C >> 1 // OFFSET_sC
	MAXARG_Bx  = (1 << SIZE_Bx) - 1
	MAXARG_sBx = MAXARG_Bx >> 1 // OFFSET_sBx
	MAXARG_sJ  = ((1 << SIZE_sJ) - 1) >> 1
)

/*
 31      23      15       7      0
  ┌------┐┌------┐-┌------┐┌-----┐
  |  C:8 ||  B:8 |k|  A:8 || Op:7| iABC
  └------┘└------┘-└------┘└-----┘
  ┌---------------┐┌------┐┌-----┐
  |     Bx:17     ||  A:8 || Op:7| iABx
  └---------------┘└------┘└-----┘
  ┌---------------┐┌------┐┌-----┐
  |    sBx:17     ||  A:8 || Op:7| iAsBx
  └---------------┘└------┘└-----┘
  ┌-----------------------┐┌-----┐
  |          Ax:25        || Op:7| iAx
  └-----------------------┘└-----┘
  ┌-----------------------┐┌-----┐
  |          sJ:25        || Op:7| isJ
  └-----------------------┘└-----┘
*/
type Instruction uint32

func (instr Instruction) Execute(vm api.LuaVM) {
	//debugPrint(instr)
	opTable[instr.Opcode()].action(instr, vm)
}

func (instr Instruction) Opcode() int {
	return int(instr & 0x7F)
}

func (instr Instruction) ABC() (a, b, c int) {
	a, b, c, _ = instr.ABCk()
	return
}
func (instr Instruction) AsBk() (a, sB, k int) {
	a, sB, _, k = instr.ABCk()
	sB -= MAXARG_sC
	return
}
func (instr Instruction) ABsC() (a, b, sC int) {
	a, b, sC, _ = instr.ABCk()
	sC -= MAXARG_sC
	return
}

func (instr Instruction) ABCk() (a, b, c, k int) {
	a = int(instr >> 7 & 0xFF)
	k = int(instr >> 15 & 0x01)
	b = int(instr >> 16 & 0xFF)
	c = int(instr >> 24 & 0xFF)
	return
}

func (instr Instruction) ABx() (a, bx int) {
	a = int(instr >> 7 & 0xFF)
	bx = int(instr >> 15)
	return
}

func (instr Instruction) AsBx() (a, sbx int) {
	a, bx := instr.ABx()
	return a, bx - MAXARG_sBx
}

func (instr Instruction) Ax() int {
	return int(instr >> 7)
}
func (instr Instruction) sJ() int {
	return instr.Ax() - MAXARG_sJ
}

func (instr Instruction) OpName() string {
	return opTable[instr.Opcode()].name
}

func (instr Instruction) OpMode() byte {
	return opTable[instr.Opcode()].opMode
}

// debug
func debugPrint(instr Instruction) {
	print(instr.OpName(), " ")
	switch instr.OpMode() {
	case iABC:
		a, b, c := instr.ABC()
		println(a, b, c)
	case iABx:
		a, bx := instr.ABx()
		println(a, bx)
	case iAsBx:
		a, sBx := instr.AsBx()
		println(a, sBx)
	case iAx:
		ax := instr.Ax()
		println(ax)
	case isJ:
		// TODO
	}
}
