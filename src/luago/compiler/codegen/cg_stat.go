package codegen

import . "luago/compiler/ast"

func (self *codeGen) cgStat(node Stat) {
	switch stat := node.(type) {
	case DoStat:
		self.cgBlockWithNewScope(stat, false)
	case FuncCallStat:
		self.cgFuncCallStat(stat)
	case *BreakStat:
		self.cgBreakStat(stat)
	case *RepeatStat:
		self.cgRepeatStat(stat)
	case *WhileStat:
		self.cgWhileStat(stat)
	case *IfStat:
		self.cgIfStat(stat)
	case *ForNumStat:
		self.cgForNumStat(stat)
	case *ForInStat:
		self.cgForInStat(stat)
	case *LocalAssignStat:
		self.cgLocalAssignStat(stat)
	case *AssignStat:
		self.cgAssignStat(stat)
	case *LabelStat, *GotoStat:
		panic("label and goto statements are not supported!")
	}
}

func (self *codeGen) cgFuncCallStat(node FuncCallStat) {
	fcExp := (*FuncCallExp)(node)
	tmp := self.allocTmp()
	self.cgExp(fcExp, tmp, 0)
	self.freeTmp()
}

func (self *codeGen) cgBreakStat(node *BreakStat) {
	pc := self.emitJmp(node.Line, 0)
	self.addBreakJmp(pc)
}

/*
        ______________
       |  false? jmp  |
       V              /
repeat block until exp
*/
func (self *codeGen) cgRepeatStat(node *RepeatStat) {
	self.enterScope(true)
	pcBeforeBlock := self.pc()
	self.cgBlock(node.Block)

	exp := node.Exp
	if !isTrueAtCompileTime(exp) {
		exp = castNilToFalse(exp)
		tmp, _ := self.exp2OpArgX(exp, ARG_REG)
		line := lastLineOfExp(exp)

		self.emitTest(line, tmp, 0)
		self.emitJmp(line, pcBeforeBlock-self.pc()-1)
	} else {
		if strExp, ok := exp.(*StringExp); ok {
			self.indexOfConstant(strExp.Str)
		}
	}

	self.exitScope(self.pc() + 1)
}

/*
           ______________
          /  false? jmp  |
         /               |
while exp do block end <-'
      ^           \
      |___________/
           jmp
*/
func (self *codeGen) cgWhileStat(node *WhileStat) {
	pcBeforeExp := self.pc()
	pcOfJmpToEnd := -1

	exp := node.Exp
	if !isTrueAtCompileTime(exp) {
		exp = castNilToFalse(exp)
		tmp, _ := self.exp2OpArgX(exp, ARG_REG)
		line := lastLineOfExp(exp)

		self.emitTest(line, tmp, 0)
		pcOfJmpToEnd = self.emitJmp(line, 0)
	} else {
		if strExp, ok := exp.(*StringExp); ok {
			self.indexOfConstant(strExp.Str)
		}
	}

	self.enterScope(true)
	self.cgBlock(node.Block)
	self.emitJmp(node.Block.LastLine, pcBeforeExp-self.pc()-1)
	self.exitScope(self.pc())

	if pcOfJmpToEnd >= 0 {
		self.fixSbx(pcOfJmpToEnd, self.pc()-pcOfJmpToEnd)
	}
}

/*
         _________________       _________________       _____________
        / false? jmp      |     / false? jmp      |     / false? jmp  |
       /                  V    /                  V    /              V
if exp1 then block1 elseif exp2 then block2 elseif true then block3 end <-.
                   \                       \                       \      |
                    \_______________________\_______________________\_____|
                    jmp                     jmp                     jmp
*/
func (self *codeGen) cgIfStat(node *IfStat) {
	pcOfJmpToEnds := make([]int, 0, len(node.Exps))
	pcOfJmpToElseif := -1

	for i := 0; i < len(node.Exps); i++ {
		exp := node.Exps[i]
		block := node.Blocks[i]

		if pcOfJmpToElseif >= 0 {
			self.fixSbx(pcOfJmpToElseif, self.pc()-pcOfJmpToElseif)
		}
		if !isTrueAtCompileTime(exp) {
			tmp, _ := self.exp2OpArgX(exp, ARG_REG)
			line := lastLineOfExp(exp)

			self.emitTest(line, tmp, 0)
			pcOfJmpToElseif = self.emitJmp(line, 0)
		} else {
			pcOfJmpToElseif = -1
			if strExp, ok := exp.(*StringExp); ok {
				self.indexOfConstant(strExp.Str)
			}
		}

		self.cgBlockWithNewScope(block, false)
		if i < len(node.Exps)-1 {
			pc := self.emitJmp(block.LastLine, 0)
			pcOfJmpToEnds = append(pcOfJmpToEnds, pc)
		}
	}

	if pcOfJmpToElseif >= 0 {
		self.fixSbx(pcOfJmpToElseif, self.pc()-pcOfJmpToElseif)
	}
	for _, pc := range pcOfJmpToEnds {
		self.fixSbx(pc, self.pc() - pc)
	}
}

func (self *codeGen) cgForNumStat(node *ForNumStat) {
	forIdxVar := "(for index)"
	forLmtVar := "(for limit)"
	forStpVar := "(for step)"

	self.enterScope(true)

	self.cgStat(&LocalAssignStat{
		LastLine: node.LineOfFor,
		NameList: []string{forIdxVar, forLmtVar, forStpVar},
		ExpList:  []Exp{node.InitExp, node.LimitExp, node.StepExp},
	})
	self.addLocVar(node.VarName, self.pc()+2)

	a := self.scope.stackSize - 3
	prepPc := self.emitForPrep(node.LineOfDo, a, 0)
	self.cgBlock(node.Block)
	loopPc := self.emitForLoop(node.LineOfFor, a, 0)

	self.fixSbx(prepPc, loopPc-prepPc-1)
	self.fixSbx(loopPc, prepPc-loopPc)

	self.exitScope(self.pc())
	self.fixEndPc(forIdxVar, 1)
	self.fixEndPc(forLmtVar, 1)
	self.fixEndPc(forStpVar, 1)
}

func (self *codeGen) cgForInStat(node *ForInStat) {
	forGeneratorVar := "(for generator)"
	forStateVar := "(for state)"
	forControlVar := "(for control)"

	self.enterScope(true)

	self.cgStat(&LocalAssignStat{
		//LastLine: 0,
		NameList: []string{forGeneratorVar, forStateVar, forControlVar},
		ExpList:  node.ExpList,
	})
	for _, name := range node.NameList {
		self.addLocVar(name, self.pc()+2)
	}

	// todo: ???
	if len(node.NameList) < 3 {
		n := 3 - len(node.NameList)
		self.allocTmps(n)
		self.freeTmps(n)
	}

	jmpToTFC := self.emitJmp(node.LineOfDo, 0)
	self.cgBlock(node.Block)
	self.fixSbx(jmpToTFC, self.pc()-jmpToTFC)

	line := lineOfExp(node.ExpList[0])
	slotOfGeneratorVar := self.slotOf(forGeneratorVar)
	self.emitTForCall(line, slotOfGeneratorVar, len(node.NameList))
	self.emitTForLoop(line, slotOfGeneratorVar+2, jmpToTFC-self.pc()-1)

	self.exitScope(self.pc() - 1)
	self.fixEndPc(forGeneratorVar, 2)
	self.fixEndPc(forStateVar, 2)
	self.fixEndPc(forControlVar, 2)
}

func (self *codeGen) cgLocalAssignStat(node *LocalAssignStat) {
	if len(node.ExpList) == 1 {
		exp0 := node.ExpList[0]
		if fd, ok := exp0.(*FuncDefExp); ok {
			if !fd.IsAno {
				name := node.NameList[0]
				slot := self.addLocVar(name, self.pc()+2)
				self.cgExp(exp0, slot, 1)
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
				self.cgExp(exps[i], a, n)
			} else {
				self.cgExp(exps[i], a, 1)
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
				self.cgExp(exps[i], a, 0)
			} else {
				self.cgExp(exps[i], a, 1)
			}
		}
		self.freeTmps(nExps - nNames)
	}

	startPc := self.pc() + 1
	for _, name := range node.NameList {
		self.addLocVar(name, startPc)
	}
}

func (self *codeGen) cgAssignStat(node *AssignStat) {
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
		strIdx := self.indexOfConstant(name)

		arg, argKind := self.toOpArg(rhs)
		switch argKind {
		case ARG_CONST, ARG_REG:
			self.emitSetTabUp(line, envIdx, strIdx, arg)
		default:
			tmp := self.allocTmp()
			self.cgExp(rhs, tmp, 1)
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
		self.cgExp(exp, a, 1)
	default:
		tmp := self.allocTmp()
		self.cgExp(exp, tmp, 1)
		self.freeTmp()
		self.emitMove(line, a, tmp)
	}
}

func (self *codeGen) assignToUpval(node Exp, line, idx int) {
	if slot, ok := self.isLocVar(node); ok {
		self.emitSetUpval(line, slot, idx)
	} else {
		tmp := self.allocTmp()
		self.cgExp(node, tmp, 1)
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
		self.cgExp(lhs.PrefixExp, a, 1)
	}

	iKey, tKey := self.toOpArg(lhs.KeyExp)
	if tKey == ARG_CONST || tKey == ARG_REG {
		b = iKey
	} else {
		b = self.allocTmp()
		nTmps++
		self.cgExp(lhs.KeyExp, b, 1)
	}

	iVal, tVal := self.toOpArg(rhs)
	if tVal == ARG_CONST || tVal == ARG_REG {
		c = iVal
	} else {
		c = self.allocTmp()
		nTmps++
		self.cgExp(rhs, c, 1)
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
				self.cgExp(rhs, argVal, 1)
				operands[i*6+4] = argVal
				operands[i*6+5] = kindVal
			} else { // i > nVars-1
				argVal := allocator.allocTmp()
				self.cgExp(rhs, argVal, 0)
			}
		} else { // last exp & (vararg | funccall)
			n := nVars - nExps + 1
			a := allocator.allocTmp()
			self.cgExp(rhs, a, n)
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
