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
		if nameExp, ok := node.ExpList[0].(*NameExp); ok {
			if r := self.indexOfLocVar(nameExp.Name); r >= 0 {
				self.emitReturn(node.LastLine, r, 1)
				return
			}
		}
		if fcExp, ok := node.ExpList[0].(*FuncCallExp); ok {
			r := self.allocReg()
			self.cgTailCallExp(fcExp, r)
			self.freeReg()
			self.emitReturn(node.LastLine, r, -1)
			return
		}
	}

	multRet := isVarargOrFuncCallExp(node.ExpList[nExps-1])
	for i, exp := range node.ExpList {
		r := self.allocReg()
		if i == nExps-1 && multRet {
			self.cgExp(exp, r, -1)
		} else {
			self.cgExp(exp, r, 1)
		}
	}
	self.freeRegs(nExps)

	a := self.scope.nLocals // correct?
	if multRet {
		self.emitReturn(node.LastLine, a, -1)
	} else {
		self.emitReturn(node.LastLine, a, nExps)
	}
}
