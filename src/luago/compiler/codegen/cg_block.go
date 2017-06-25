package codegen

import . "luago/compiler/ast"

// todo: rename
func (self *codeGen) cgBlockWithNewScope(node *Block, breakable bool) {
	self.enterScope(breakable)
	self.cgBlock(node)
	self.exitScope(self.pc() + 1)
}

func (self *codeGen) cgBlock(node *Block) {
	for _, stat := range node.Stats {
		self.cgStat(stat)
	}

	if node.RetStat != nil {
		self.retStat(node.RetStat)
	}
}

func (self *codeGen) retStat(node *RetStat) {
	nExps := len(node.ExpList)
	if nExps == 0 {
		self.emitReturn(node.LastLine, 0, 0)
		return
	}

	if nExps == 1 {
		switch exp := node.ExpList[0].(type) {
		case *NameExp:
			if slot := self.slotOf(exp.Name); slot >= 0 {
				self.emitReturn(node.LastLine, slot, 1)
				return
			}
		case *FuncCallExp:
			tmp := self.allocTmp()
			self.cgTailCallExp(exp, tmp)
			self.freeTmp()
			self.emitReturn(node.LastLine, tmp, -1)
			return
		}
	}

	lastExpIsVarargOrFuncCall := false
	for i, exp := range node.ExpList {
		tmp := self.allocTmp()
		if i == nExps-1 && isVarargOrFuncCallExp(exp) {
			lastExpIsVarargOrFuncCall = true
			self.exp(exp, tmp, -1)
		} else {
			self.exp(exp, tmp, 1)
		}
	}
	self.freeTmps(nExps)

	a := self.scope.nLocals // correct?
	if lastExpIsVarargOrFuncCall {
		self.emitReturn(node.LastLine, a, -1)
	} else {
		self.emitReturn(node.LastLine, a, nExps)
	}
}
