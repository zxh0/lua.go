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
	if nExps == 1 {
		if fcExp, ok := node.ExpList[0].(*FuncCallExp); ok {
			tmp := self.allocTmp()
			self.tailCallExp(fcExp, tmp)
			self.freeTmp()
			self._return(node.LastLine, tmp, -1)
			return
		}
	}

	for _, exp := range node.ExpList {
		tmp := self.allocTmp()
		self.exp(exp, tmp, 1)
	}
	self.freeTmps(nExps)
	a := self.scope.nLocals // correct?
	self._return(node.LastLine, a, nExps)
}
