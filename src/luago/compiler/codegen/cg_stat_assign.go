package codegen

import . "luago/compiler/ast"

func (self *codeGen) localAssignStat(node *LocalAssignStat) {
	if len(node.ExpList) == 1 {
		exp0 := node.ExpList[0]
		if fd, ok := exp0.(*FuncDefExp); ok {
			if !fd.IsAno {
				name := node.NameList[0]
				slot := self.addLocVar(name, self.pc()+2)
				self.exp(exp0, slot, 1)
				return
			}
		}
	}

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
		self.emitLoadNil(node.LastLine, a, n)
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

	startPc := self.pc() + 1
	for _, name := range node.NameList {
		self.addLocVar(name, startPc)
	}
}

func (self *codeGen) assignStat(node *AssignStat) {
	if len(node.VarList) == 1 && len(node.ExpList) == 1 {
		self.assignStat1(node.LastLine,
			node.VarList[0], node.ExpList[0])
	} else {
		self.assignStatN(node)
	}
}

func (self *codeGen) assignStat1(line int, lhs, rhs Exp) {
	switch x := lhs.(type) {
	case *NameExp: // x = y
		self.assignToName(line, x.Name, rhs)
	case *BracketsExp: // k[v] = x
		self.assignToField(line, x, rhs)
	}
}

func (self *codeGen) assignToName(line int, name string, rhs Exp) {
	if slot := self.slotOf(name); slot >= 0 {
		self.assignToLocalVar(rhs, line, slot)
	} else if idx := self.lookupUpval(name); idx >= 0 {
		self.assignToUpval(rhs, line, idx)
	} else {
		envIdx := self.lookupUpval("_ENV")
		strIdx := self.indexOf(name)

		arg, argKind := self.toOpArg(rhs)
		switch argKind {
		case ARG_CONST, ARG_REG:
			self.emitSetTabUp(line, envIdx, strIdx, arg)
		default:
			tmp := self.allocTmp()
			self.exp(rhs, tmp, 1)
			self.freeTmp()
			self.emitSetTabUp(line, envIdx, strIdx, tmp) // todo
		}
	}
}

func (self *codeGen) assignToLocalVar(node Exp, line, a int) {
	switch exp := node.(type) {
	case *NilExp, *FalseExp, *TrueExp,
		*IntegerExp, *FloatExp, *StringExp,
		*VarargExp,
		*BinopExp, *UnopExp,
		*NameExp, *BracketsExp:
		self.exp(exp, a, 1)
	default:
		tmp := self.allocTmp()
		self.exp(exp, tmp, 1)
		self.freeTmp()
		self.emitMove(line, a, tmp)
	}
}

func (self *codeGen) assignToUpval(node Exp, line, idx int) {
	if slot, ok := self.isLocVar(node); ok {
		self.emitSetUpval(line, slot, idx)
	} else {
		tmp := self.allocTmp()
		self.exp(node, tmp, 1)
		self.freeTmp()
		self.emitSetUpval(line, tmp, idx)
	}
}

// func (self *codeGen) assignToGlobalVar(node Exp, line, envIdx, nameIdx int) {

// }

func (self *codeGen) assignToField(line int, lhs *BracketsExp, rhs Exp) {
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
		self.emitSetTabUp(line, a, b, c)
	} else {
		self.emitSetTable(line, a, b, c)
	}
}

func (self *codeGen) assignStatN(node *AssignStat) {
	exps := removeTailNils(node.ExpList)
	nExps := len(exps)
	nVars := len(node.VarList)

	operands := make([]int, 0, 6*nVars)
	allocator := self.newTmpAllocator(-1)
	lastExpIsVarargOrFuncCall := nExps > 0 &&
		isVarargOrFuncCallExp(exps[nExps-1])
	// lastVarIsNameExp := isNameExp(node.VarList[nVars-1])

	for _, lhs := range node.VarList {
		// if i == nVars-1 && nExps == nVars && lastVarIsNameExp {
		// 	name := lhs.(*NameExp).Name
		// 	self.isGlobalVar(name) // todo
		// 	break
		// }

		var argTab, kindTab, argKey, kindKey, argVal, kindVal int
		switch x := lhs.(type) {
		case *NameExp:
			if envIdx, nameIdx, ok := self.isGlobalVar(x.Name); ok {
				argTab, kindTab = envIdx, ARG_UPVAL
				argKey, kindKey = nameIdx, ARG_CONST
			} else {
				argKey, kindKey = self.exp2OpArg(x, ARG_RU, allocator)
			}
		case *BracketsExp:
			argTab, kindTab = self.exp2OpArg(x.PrefixExp, ARG_RU, allocator)
			argKey, kindKey = self.exp2OpArg(x.KeyExp, ARG_RUK, allocator)
		}
		operands = append(operands,
			argTab, kindTab, argKey, kindKey, argVal, kindVal)
	}
	for i, rhs := range exps {
		// if i == nExps-1 && nExps == nVars && lastVarIsNameExp {
		// 	break
		// }
		if i == nExps-1 && nExps == nVars {
			lhs := node.VarList[nExps-1]
			if isNameExp(lhs) {
				name := lhs.(*NameExp).Name
				self.assignToName(node.LastLine, name, rhs)
			} else {
				argTab := operands[i*6+0]
				//kindTab := operands[i*6+1]
				argKey := operands[i*6+2]
				//kindKey := operands[i*6+3]
				argVal, _ := self.exp2OpArg(rhs, ARG_RK, allocator)
				self.emitSetTable(node.LastLine, argTab, argKey, argVal)
			}

			break
		}
		if i < nExps-1 || nExps >= nVars || !lastExpIsVarargOrFuncCall {
			if i <= nVars-1 {
				kindVal := ARG_REG
				argVal := allocator.allocTmp()
				self.exp(rhs, argVal, 1)
				operands[i*6+4] = argVal
				operands[i*6+5] = kindVal
			} else { // i > nVars-1
				argVal := allocator.allocTmp()
				self.exp(rhs, argVal, 0)
			}
		} else { // last exp & (vararg | funccall)
			n := nVars - nExps + 1
			a := allocator.allocTmp()
			self.exp(rhs, a, n)
			for j := nExps - 1; j < nVars; j++ {
				operands[j*6+4] = a + j - nExps + 1
				operands[j*6+5] = ARG_REG
			}
		}
	}
	if nExps < nVars && !lastExpIsVarargOrFuncCall {
		n := nVars - nExps
		a := allocator.allocTmps(n)
		self.emitLoadNil(node.LastLine, a, n)
		for j := nExps; j < nVars; j++ {
			operands[j*6+4] = a + j - nExps
			operands[j*6+5] = ARG_REG
		}
	}
	// if nExps == nVars && lastVarIsNameExp {
	// 	self.assignStat1(node.LastLine, node.VarList[nVars-1], exps[nExps-1])
	// }

	allocator.freeAll()

	for i := len(operands) - 1; i > 0; i -= 6 {
		argTab := operands[i-5]
		kindTab := operands[i-4]
		argKey := operands[i-3]
		kindKey := operands[i-2]
		argVal := operands[i-1]
		kindVal := operands[i]

		switch kindTab {
		case ARG_REG:
			if kindVal == ARG_REG {
				self.emitSetTable(node.LastLine, argTab, argKey, argVal)
			}
		case ARG_UPVAL:
			if kindVal == ARG_REG {
				self.emitSetTabUp(node.LastLine, argTab, argKey, argVal)
			}
		default:
			if kindKey == ARG_UPVAL {
				if kindVal == ARG_REG {
					self.emitSetUpval(node.LastLine, argVal, argKey)
				}
			} else {
				if kindVal == ARG_REG {
					self.emitMove(node.LastLine, argKey, argVal)
				}
			}
		}

		// if kindTab == -1 {
		// 	if kindKey == ARG_REG {
		// 		self.emitMove(node.LastLine, argKey, argVal)
		// 	} else if kindKey == ARG_UPVAL {
		// 		self.emitSetUpval(node.LastLine, argVal, argKey)
		// 	}
		// } else if kindTab == ARG_GLOBAL {
		// 	self.emitSetTabUp(node.LastLine, argTab, argKey, argVal)
		// }
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
