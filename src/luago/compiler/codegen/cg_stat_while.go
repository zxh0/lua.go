package codegen

import . "luago/compiler/ast"

/*
        while
    .-> (exp) --.
jmp2|    do     |jmp1
    '--[block]  |
         end <--'
*/
func (self *codeGen) whileStat(node *WhileStat) {
	if nilExp, ok := node.Exp.(*NilExp); ok {
		node.Exp = &FalseExp{nilExp.Line}
	}

	var jmpToEndPcs []int
	startPc := self.pc()
	endless := isExpTrue(node.Exp)

	if !endless {
		jmpToEndPcs = self.testExp(node.Exp, node.Line)
	}

	self.blockWithNewScope(node.Block)

	endPc := self.pc()
	self.emitJmp(node.Block.LastLine, startPc-endPc-1)

	if !endless && jmpToEndPcs != nil {
		for _, pc := range jmpToEndPcs {
			self.fixSbx(pc, endPc-pc+1)
		}
	}
}
