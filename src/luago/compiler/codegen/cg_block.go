package codegen

import . "luago/compiler/ast"

func (self *cg) block(node *Block) {
	for _, stat := range node.Stats {
		self.stat(stat)
	}

	if node.RetStat != nil {
		self.retStat(node.RetStat)
	}
}

func (self *cg) retStat(node *RetStat) {
	nExps := len(node.ExpList)
	if nExps == 0 {
		self._return(node.LastLine, 0, 0)
		return
	}

	if nExps == 1 {
		switch exp := node.ExpList[0].(type) {
		case *NameExp:
			if slot := self.slotOf(exp.Name); slot >= 0 {
				self._return(node.LastLine, slot, 1)
				return
			}
		case *FuncCallExp:
			tmp := self.allocTmp()
			self.tailCallExp(exp, tmp)
			self.freeTmp()
			self._return(node.LastLine, tmp, -1)
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
		self._return(node.LastLine, a, -1)
	} else {
		self._return(node.LastLine, a, nExps)
	}
}
