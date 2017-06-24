package codegen

import . "luago/compiler/ast"

func (self *codeGen) cgStat(node Stat) {
	switch stat := node.(type) {
	case DoStat:
		self.cgBlock(stat)
	case FuncCallStat:
		self.cgFuncCallStat(stat)
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
	case *BreakStat:
		// todo
	}
}

func (self *codeGen) cgFuncCallStat(node FuncCallStat) {
	fcExp := (*FuncCallExp)(node)
	tmp := self.allocTmp()
	self.exp(fcExp, tmp, 0)
	self.freeTmp()
}

/*
        ______________
       |  false? jmp  |
       V              /
repeat block until exp
*/
func (self *codeGen) cgRepeatStat(node *RepeatStat) {
	self.enterScope()
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

	self.cgBlockWithNewScope(node.Block)
	self.emitJmp(node.Block.LastLine, pcBeforeExp-self.pc()-1)

	if pcOfJmpToEnd >= 0 {
		self.fixSbx(pcOfJmpToEnd, self.pc()-pcOfJmpToEnd)
	}
}

/*
        _____________           _____________           _____________
       / false? jmp  |         / false? jmp  |         / false? jmp  |
      /              V        /              V        /              V
if exp1 then block1 elseif exp2 then block2 elseif true then block3 end <-.
                  \                       \                       \       |
                   \_______________________\_______________________\______|
                   jmp                     jmp                     jmp
*/
func (self *codeGen) cgIfStat(node *IfStat) {
	jmp2elseIfs := map[int]bool{}
	jmp2ends := map[int]bool{}

	for i := 0; i < len(node.Exps); i++ {
		if i > 0 {
			for pc, _ := range jmp2elseIfs {
				self.fixSbx(pc, self.pc()-pc)
			}
			jmp2elseIfs = map[int]bool{} // clear map
		}

		self.ifExpBlock(node, i, jmp2elseIfs, jmp2ends)
	}

	for pc, _ := range jmp2elseIfs {
		self.fixSbx(pc, self.pc()-pc)
	}
	for pc, _ := range jmp2ends {
		self.fixSbx(pc, self.pc()-pc)
	}
}

// todo: rename
func (self *codeGen) ifExpBlock(node *IfStat, i int,
	jmp2elseIfs, jmp2ends map[int]bool) {

	exp := node.Exps[i]
	block := node.Blocks[i]
	lineOfThen := node.Lines[i]

	if isExpTrue(exp) {
		if strExp, ok := exp.(*StringExp); ok {
			self.indexOfConstant(strExp.Str)
		}
	} else {
		pendingJmps := self.testExp(exp, lineOfThen)
		for _, pc := range pendingJmps {
			jmp2elseIfs[pc] = true
		}
	}

	self.cgBlockWithNewScope(block)
	if i < len(node.Exps)-1 {
		pc := self.emitJmp(block.LastLine, 0)
		jmp2ends[pc] = true
	}
}

func (self *codeGen) cgForNumStat(node *ForNumStat) {
	forIdxVar := "(for index)"
	forLmtVar := "(for limit)"
	forStpVar := "(for step)"

	self.enterScope()

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

	self.enterScope()

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
