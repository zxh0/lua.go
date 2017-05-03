package vm

import "luago/lua"

const HALF_BX = (1<<18 - 1) >> 1 // 131071

/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|opc6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|opc6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|opc6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |opc6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/
type Instruction uint32

func (self Instruction) Execute(vm lua.VM) {
	opcodes[self.Opcode()].action(self, vm)
}

func (self Instruction) Opcode() int {
	return int(self & 0x3F)
}

func (self Instruction) ABC() (a, b, c int) {
	a = int(self >> 6 & 0xFF)
	c = int(self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return
}

func (self Instruction) ABx() (a, bx int) {
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

func (self Instruction) AsBx() (a, sbx int) {
	a, bx := self.ABx()
	return a, bx - HALF_BX
}

func (self Instruction) Ax() int {
	return int(self >> 6)
}

func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}

func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}
