package codegen

import . "luago/compiler/ast"
import . "luago/lua/vm"

/*
        while
    .-> (exp) --.
jmp2|    do     |jmp1
    '--[block]  |
         end <--'
*/
func (self *cg) whileStat(node *WhileStat) {
	if nilExp, ok := node.Exp.(*NilExp); ok {
		node.Exp = &FalseExp{nilExp.Line}
	}

	jmp1Pc := 0
	startPc := self.pc()
	endless := isExpTrue(node.Exp)

	if !endless {
		//self.exp(node.Exp, STAT_WHILE, 0)

		self.inst(node.Line, OP_TEST, 0, 0, 0)         // todo
		jmp1Pc = self.inst(node.Line, OP_JMP, 0, 0, 0) // todo
		self.freeTmp()
	}

	self.block(node.Block)

	endPc := self.pc()
	self.inst(node.Line, OP_JMP, 0, startPc-endPc-1, 0) // todo

	if !endless {
		self.fixSbx(jmp1Pc, endPc-jmp1Pc-1)
	}
}
