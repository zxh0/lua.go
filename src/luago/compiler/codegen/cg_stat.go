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
	self.exp(fcExp, tmp, 0)
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
