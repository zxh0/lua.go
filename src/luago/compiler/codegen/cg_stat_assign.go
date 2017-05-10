package codegen

import . "luago/compiler/ast"

func (self *cg) localAssignStat(node *LocalAssignStat) {
	exps := removeTailNils(node.ExpList)
	nExps := len(exps)
	nNames := len(node.NameList)
	lastExpIsVarargOrFuncCall := false

	for i := 0; i < nNames; i++ {
		a := self.allocTmp()
		if i < nExps {
			if i == nExps-1 && isVarargOrFuncCallExp(exps[i]) {
				lastExpIsVarargOrFuncCall = true
				n := nNames - nExps + 1
				self.exp(exps[i], a, n)
			} else {
				self.exp(exps[i], a, 1)
			}
		}
	}

	if nNames > nExps && !lastExpIsVarargOrFuncCall {
		n := nNames - nExps
		a := self.scope.stackSize - n
		self.loadNil(node.LastLine, a, n)
	} else if nExps > nNames {
		for i := nNames; i < nExps; i++ {
			a := self.allocTmp()
			if i == nExps-1 && isVarargOrFuncCallExp(exps[i]) {
				self.exp(exps[i], a, 0)
			} else {
				self.exp(exps[i], a, 1)
			}
		}
		self.freeTmps(nExps - nNames)
	}

	startPc := self.pc() - 1
	for _, name := range node.NameList {
		self.addLocVar(name, startPc)
	}
}

func (self *cg) assignStat(node *AssignStat) {
	if len(node.VarList) == 1 && len(node.ExpList) == 1 {
		self.assignStat1(node.LastLine,
			node.VarList[0], node.ExpList[0])
	} else {
		self.assignStatN(node)
	}
}

func (self *cg) assignStat1(line int, lhs, rhs Exp) {
	switch x := lhs.(type) {
	case *NameExp: // x = y
		self.assignToName(line, x.Name, rhs)
	case *BracketsExp: // k[v] = x
		self.assignToField(line, x, rhs)
	}
}

func (self *cg) assignToName(line int, name string, rhs Exp) {
	if slot := self.slotOf(name); slot >= 0 {
		iY, tY := self.toOperand0(rhs)
		switch tY {
		case ARG_CONST:
			self.exp(rhs, slot, 1)
		case ARG_REG:
			self.move(line, slot, iY)
		case ARG_UPVAL:
			self.getUpval(line, slot, iY)
		case ARG_GLOBAL: // todo
			envIdx := self.lookupUpval("_ENV")
			self.getTabUp(line, slot, envIdx, iY)
		default:
			// tmp := self.allocTmp()
			// self.exp(rhs, tmp, 1)
			// self.freeTmp()
			switch rhs.(type) {
			case *UnopExp:
				self.exp(rhs, slot, 1)
			case *BinopExp:
				//tmp := self.allocTmp()
				self.exp(rhs, slot, 1)
				//self.freeTmp()
				//self.fixA(len(self.insts)-1, slot)
			case *BracketsExp:
				self.exp(rhs, slot, 1)
			default:
				tmp := self.allocTmp()
				self.exp(rhs, tmp, 1)
				self.freeTmp()
				self.move(line, slot, tmp)
			}
		}
	} else if idx := self.lookupUpval(name); idx >= 0 {
		iY, tY := self.toOperand0(rhs)
		switch tY {
		case ARG_REG:
			self.setUpval(line, iY, idx)
		default:
			tmp := self.allocTmp()
			self.exp(rhs, tmp, 1)
			self.freeTmp()
			self.setUpval(line, tmp, idx)
		}
	} else {
		envIdx := self.lookupUpval("_ENV")
		strIdx := self.indexOf(name)

		iY, tY := self.toOpArg(rhs)
		switch tY {
		case ARG_CONST:
			self.setTabUp(line, envIdx, strIdx, iY)
		case ARG_REG:
			self.setTabUp(line, envIdx, strIdx, iY)
		default:
			tmp := self.allocTmp()
			self.exp(rhs, tmp, 1)
			self.freeTmp()
			self.setTabUp(line, envIdx, strIdx, tmp) // todo
		}
	}
}

func (self *cg) assignToField(line int, lhs *BracketsExp, rhs Exp) {
	var a, b, c int
	nTmps := 0

	iTab, tTab := self.toOpArg(lhs.PrefixExp)
	if tTab == ARG_REG || tTab == ARG_UPVAL {
		a = iTab
	} else {
		a = self.allocTmp()
		nTmps++
		self.exp(lhs.PrefixExp, a, 1)
	}

	iKey, tKey := self.toOpArg(lhs.KeyExp)
	if tKey == ARG_CONST || tKey == ARG_REG {
		b = iKey
	} else {
		b = self.allocTmp()
		nTmps++
		self.exp(lhs.KeyExp, b, 1)
	}

	iVal, tVal := self.toOpArg(rhs)
	if tVal == ARG_CONST || tVal == ARG_REG {
		c = iVal
	} else {
		c = self.allocTmp()
		nTmps++
		self.exp(rhs, c, 1)
	}

	self.freeTmps(nTmps)
	if tTab == ARG_UPVAL {
		self.setTabUp(line, a, b, c)
	} else {
		self.setTable(line, a, b, c)
	}
}

func (self *cg) assignStatN(node *AssignStat) {
	exps := removeTailNils(node.ExpList)
	nExps := len(exps)
	nVars := len(node.VarList)
	nTmps := 0

	operands := make([]int, 0, 16)
	flag := false // todo: rename

	for i, lhs := range node.VarList {
		var iTab, tTab, iKey, tKey int
		iVal := 0 // placeholder
		tVal := 0 // placeholder
		switch x := lhs.(type) {
		case *NameExp:
			if slot := self.slotOf(x.Name); slot >= 0 {
				flag = nExps == nVars && i == nVars-1
				iTab = -1
				tTab = -1
				iKey = slot
				tKey = ARG_REG
			} else if idx := self.lookupUpval(x.Name); idx >= 0 {
				flag = nExps == nVars && i == nVars-1
				iTab = -1
				tTab = -1
				iKey = idx
				tKey = ARG_UPVAL
			} else {
				iTab = self.lookupUpval("_ENV")
				tTab = ARG_GLOBAL
				iKey = self.indexOf(x.Name)
				tKey = ARG_CONST
			}
			if !flag {
				operands = append(operands, iTab, tTab, iKey, tKey, iVal, tVal)
			}
		case *BracketsExp:
			// todo
		}
	}
	for i, rhs := range exps {
		if i < nExps-1 || !flag {
			//iVal, tVal := self.toOpArg(rhs)
			tVal := ARG_REG
			iVal := self.allocTmp()
			nTmps++
			self.exp(rhs, iVal, 1)
			if i*6 < len(operands) {
				operands[i*6+4] = iVal
				operands[i*6+5] = tVal
			}
		}
	}
	if flag {
		self.assignStat1(node.LastLine,
			node.VarList[nVars-1], exps[nExps-1])
	} else if nVars > nExps {
		nNils := nVars - nExps
		tmp := self.allocTmps(nNils)
		nTmps += nNils
		self.loadNil(node.LastLine, tmp, nNils)
		for i := nExps; i < nVars; i++ {
			operands[i*6+4] = tmp + i - nExps
			operands[i*6+5] = ARG_REG
		}
	}

	self.freeTmps(nTmps)
	for i := len(operands) - 1; i > 0; i -= 6 {
		iTab := operands[i-5]
		tTab := operands[i-4]
		iKey := operands[i-3]
		tKey := operands[i-2]
		iVal := operands[i-1]
		// tVal := operands[i]
		if tTab == -1 {
			if tKey == ARG_REG {
				self.move(node.LastLine, iKey, iVal)
			} else if tKey == ARG_UPVAL {
				self.setUpval(node.LastLine, iVal, iKey)
			}
		} else if tTab == ARG_GLOBAL {
			self.setTabUp(node.LastLine, iTab, iKey, iVal)
		}
	}
}

func removeTailNils(exps []Exp) []Exp {
	for {
		if n := len(exps); n > 0 {
			if _, ok := exps[n-1].(*NilExp); ok {
				exps = exps[:n-1]
				continue
			}
		}
		break
	}
	return exps
}
