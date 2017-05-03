package codegen

import . "luago/compiler/ast"
import . "luago/lua/vm"

/*
 repeat
[block]<-.
 until   |jmp
 (exp) --'
*/
func (self *cg) repeatStat(node *RepeatStat) {
	if nilExp, ok := node.Exp.(*NilExp); ok {
		node.Exp = &FalseExp{nilExp.Line}
	}

	pc1 := self.pc()
	self.block(node.Block)
	if !isExpTrue(node.Exp) {
		//self.exp(node.Exp, STAT_REPEAT, 0)

		line := lineOfExp(node.Exp)
		pc2 := self.inst(line, OP_TEST, 0, 0, 0) // todo
		self.inst(line, OP_JMP, 0, pc1-pc2-2, 0) // todo
		self.freeTmp()
	}
}
