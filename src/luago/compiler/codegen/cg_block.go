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

	if node.RetExps != nil {
		self.cgRetStat(node.RetExps, node.LastLine)
	}
}

func (self *codeGen) cgRetStat(exps []Exp, lastLine int) {
	nExps := len(exps)
	if nExps == 0 {
		self.emitReturn(lastLine, 0, 0)
		return
	}

	if nExps == 1 {
		if nameExp, ok := exps[0].(*NameExp); ok {
			if r := self.indexOfLocVar(nameExp.Name); r >= 0 {
				self.emitReturn(lastLine, r, 1)
				return
			}
		}
		if fcExp, ok := exps[0].(*FuncCallExp); ok {
			r := self.allocReg()
			self.cgTailCallExp(fcExp, r)
			self.freeReg()
			self.emitReturn(lastLine, r, -1)
			return
		}
	}

	multRet := isVarargOrFuncCall(exps[nExps-1])
	for i, exp := range exps {
		r := self.allocReg()
		if i == nExps-1 && multRet {
			self.cgExp(exp, r, -1)
		} else {
			self.cgExp(exp, r, 1)
		}
	}
	self.freeRegs(nExps)

	a := self.usedRegs() // correct?
	if multRet {
		self.emitReturn(lastLine, a, -1)
	} else {
		self.emitReturn(lastLine, a, nExps)
	}
}
