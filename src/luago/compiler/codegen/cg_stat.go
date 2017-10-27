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
	a := self.allocReg()
	self.cgExp(fcExp, a, 0)
	self.freeReg()
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

	a := self.allocReg()
	self.cgExp(node.Exp, a, 1)
	self.freeReg()

	line := lastLineOfExp(node.Exp)
	self.emitTest(line, a, 0)
	self.emitJmp(line, pcBeforeBlock-self.pc()-1)

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

	a := self.allocReg()
	self.cgExp(node.Exp, a, 1)
	self.freeReg()

	line := lastLineOfExp(node.Exp)
	self.emitTest(line, a, 0)
	pcJmpToEnd := self.emitJmp(line, 0)

	self.enterScope(true)
	self.cgBlock(node.Block)
	self.emitJmp(node.Block.LastLine, pcBeforeExp-self.pc()-1)
	self.exitScope(self.pc())

	self.fixSbx(pcJmpToEnd, self.pc()-pcJmpToEnd)
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
	pcJmpToEnds := make([]int, len(node.Exps))
	pcJmpToElseif := -1

	for i, exp := range node.Exps {
		if pcJmpToElseif >= 0 {
			self.fixSbx(pcJmpToElseif, self.pc()-pcJmpToElseif)
		}

		a := self.allocReg()
		self.cgExp(exp, a, 1)
		self.freeReg()

		line := lastLineOfExp(exp)
		self.emitTest(line, a, 0)
		pcJmpToElseif = self.emitJmp(line, 0)

		block := node.Blocks[i]
		self.cgBlockWithNewScope(block, false)
		if i < len(node.Exps)-1 {
			pcJmpToEnds[i] = self.emitJmp(block.LastLine, 0)
		} else {
			pcJmpToEnds[i] = pcJmpToElseif
		}
	}

	for _, pc := range pcJmpToEnds {
		self.fixSbx(pc, self.pc() - pc)
	}
}

func (self *codeGen) cgForNumStat(node *ForNumStat) {
	forIndexVar := "(for index)"
	forLimitVar := "(for limit)"
	forStepVar := "(for step)"

	self.enterScope(true)

	self.cgStat(&LocalAssignStat{
		LastLine: node.LineOfFor,
		NameList: []string{forIndexVar, forLimitVar, forStepVar},
		ExpList:  []Exp{node.InitExp, node.LimitExp, node.StepExp},
	})
	self.addLocVar(node.VarName, self.pc()+2)

	a := self.usedRegs() - 4
	prepPc := self.emitForPrep(node.LineOfDo, a, 0)
	self.cgBlock(node.Block)
	loopPc := self.emitForLoop(node.LineOfFor, a, 0)

	self.fixSbx(prepPc, loopPc-prepPc-1)
	self.fixSbx(loopPc, prepPc-loopPc)

	self.exitScope(self.pc())
	self.fixEndPc(forIndexVar, 1)
	self.fixEndPc(forLimitVar, 1)
	self.fixEndPc(forStepVar, 1)
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
		self.allocRegs(n)
		self.freeRegs(n)
	}

	jmpToTFC := self.emitJmp(node.LineOfDo, 0)
	self.cgBlock(node.Block)
	self.fixSbx(jmpToTFC, self.pc()-jmpToTFC)

	line := lineOfExp(node.ExpList[0])
	rGenerator := self.indexOfLocVar(forGeneratorVar)
	self.emitTForCall(line, rGenerator, len(node.NameList))
	self.emitTForLoop(line, rGenerator+2, jmpToTFC-self.pc()-1)

	self.exitScope(self.pc() - 1)
	self.fixEndPc(forGeneratorVar, 2)
	self.fixEndPc(forStateVar, 2)
	self.fixEndPc(forControlVar, 2)
}

func (self *codeGen) cgLocalAssignStat(node *LocalAssignStat) {
	if len(node.ExpList) == 1 {
		exp := node.ExpList[0]
		if fd, ok := exp.(*FuncDefExp); ok && !fd.IsAno {
			name := node.NameList[0]
			slot := self.addLocVar(name, self.pc()+2)
			self.cgExp(exp, slot, 1)
			return
		}
	}

	exps := removeTailNils(node.ExpList)
	nExps := len(exps)
	nNames := len(node.NameList)

	if nExps == nNames {
		for _, exp := range exps {
			a := self.allocReg()
			self.cgExp(exp, a, 1)
		}
	} else if nExps > nNames {
		for i, exp := range exps {
			a := self.allocReg()
			if i == nExps-1 && isVarargOrFuncCallExp(exp) {
				self.cgExp(exp, a, 0)
			} else {
				self.cgExp(exp, a, 1)
			}
		}
		self.freeRegs(nExps - nNames)
	} else { // nNames > nExps
		multRet := false
		for i, exp := range exps {
			a := self.allocReg()
			if i == nExps-1 && isVarargOrFuncCallExp(exp) {
				multRet = true
				n := nNames - nExps + 1
				self.cgExp(exp, a, n)
				self.allocRegs(n-1)
			} else {
				self.cgExp(exp, a, 1)
			}
		}
		if !multRet {
			n := nNames - nExps
			a := self.allocRegs(n)
			self.emitLoadNil(node.LastLine, a, n)
		}
	}

	startPc := self.pc() + 1
	for _, name := range node.NameList {
		self.addLocVar(name, startPc)
	}
}

func (self *codeGen) cgAssignStat(node *AssignStat) {
	exps := removeTailNils(node.ExpList)
	nExps := len(exps)
	nVars := len(node.VarList)

	ts := make([]int, nVars)
	ks := make([]int, nVars)
	vs := make([]int, nVars)
	oldRegs := self.usedRegs()

	for i, exp := range node.VarList {
		if bexp, ok := exp.(*BracketsExp); ok {
			ts[i] = self.allocReg()
			self.cgExp(bexp.PrefixExp, ts[i], 1)
			ks[i] = self.allocReg()
			self.cgExp(bexp.KeyExp, ks[i], 1)
		}
	}
	for i := 0; i < nVars; i++ {
		vs[i] = self.usedRegs() + i
	}

	if nExps >= nVars {
		for i, exp := range exps {
			a := self.allocReg()
			if i >= nVars && i == nExps-1 && isVarargOrFuncCallExp(exp) {
				self.cgExp(exp, a, 0)
			} else {
				self.cgExp(exp, a, 1)
			}
		}
	} else { // nVars > nExps
		multRet := false
		for i, exp := range exps {
			a := self.allocReg()
			if i == nExps-1 && isVarargOrFuncCallExp(exp) {
				multRet = true
				n := nVars - nExps + 1
				self.cgExp(exp, a, n)
				self.allocRegs(n-1)
			} else {
				self.cgExp(exp, a, 1)
			}
		}
		if !multRet {
			n := nVars - nExps
			a := self.allocRegs(n)
			self.emitLoadNil(node.LastLine, a, n)
		}
	}

	for i, exp := range node.VarList {
		if nameExp, ok := exp.(*NameExp); ok {
			varName := nameExp.Name
			if a := self.indexOfLocVar(varName); a >= 0 {
				self.emitMove(0, a, vs[i])
			} else if a := self.indexOfUpval(varName); a >= 0 {
				self.emitSetUpval(0, a, vs[i])
			} else {
				envIdx := self.indexOfUpval("_ENV")
				strIdx := self.indexOfConstant(varName)
				self.emitSetTabUp(0, envIdx, strIdx, vs[i])
			}
		} else {
			self.emitSetTable(0, ts[i], ks[i], vs[i])
		}
	}

	// todo
	self.freeRegs(self.usedRegs() - oldRegs)
}
