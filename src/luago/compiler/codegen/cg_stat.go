package codegen

import . "luago/compiler/ast"

func cgStat(fi *funcInfo, node Stat) {
	switch stat := node.(type) {
	case *DoStat:
		cgBlockWithNewScope(fi, stat.Block, false)
	case *FuncCallStat:
		cgFuncCallStat(fi, stat)
	case *BreakStat:
		cgBreakStat(fi, stat)
	case *RepeatStat:
		cgRepeatStat(fi, stat)
	case *WhileStat:
		cgWhileStat(fi, stat)
	case *IfStat:
		cgIfStat(fi, stat)
	case *ForNumStat:
		cgForNumStat(fi, stat)
	case *ForInStat:
		cgForInStat(fi, stat)
	case *AssignStat:
		cgAssignStat(fi, stat)
	case *LocalAssignStat:
		cgLocalAssignStat(fi, stat)
	case *LocalFuncDefStat:
		cgLocalFuncDefStat(fi, stat)
	case *LabelStat, *GotoStat:
		panic("label and goto statements are not supported!")
	}
}

func cgLocalFuncDefStat(fi *funcInfo, node *LocalFuncDefStat) {
	r := fi.addLocVar(node.Name, fi.pc()+2)
	cgFuncDefExp(fi, node.Exp, r)
}

func cgFuncCallStat(fi *funcInfo, node *FuncCallStat) {
	a := fi.allocReg()
	cgExp(fi, node, a, 0)
	fi.freeReg()
}

func cgBreakStat(fi *funcInfo, node *BreakStat) {
	pc := fi.emitJmp(node.Line, 0)
	fi.addBreakJmp(pc)
}

/*
        ______________
       |  false? jmp  |
       V              /
repeat block until exp
*/
func cgRepeatStat(fi *funcInfo, node *RepeatStat) {
	fi.enterScope(true)

	pcBeforeBlock := fi.pc()
	cgBlock(fi, node.Block)

	oldRegs := fi.usedRegs
	a, _ := expToOpArg(fi, node.Exp, ARG_REG)
	fi.resetRegs(oldRegs)

	line := lastLineOf(node.Exp)
	fi.emitTest(line, a, 0)
	fi.emitJmp(line, pcBeforeBlock-fi.pc()-1)

	fi.exitScope(fi.pc() + 1)
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
func cgWhileStat(fi *funcInfo, node *WhileStat) {
	pcBeforeExp := fi.pc()

	oldRegs := fi.usedRegs
	a, _ := expToOpArg(fi, node.Exp, ARG_REG)
	fi.resetRegs(oldRegs)

	line := lastLineOf(node.Exp)
	fi.emitTest(line, a, 0)
	pcJmpToEnd := fi.emitJmp(line, 0)

	fi.enterScope(true)
	cgBlock(fi, node.Block)
	fi.emitJmp(node.Block.LastLine, pcBeforeExp-fi.pc()-1)
	fi.exitScope(fi.pc())

	fi.fixSbx(pcJmpToEnd, fi.pc()-pcJmpToEnd)
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
func cgIfStat(fi *funcInfo, node *IfStat) {
	pcJmpToEnds := make([]int, len(node.Exps))
	pcJmpToElseif := -1

	for i, exp := range node.Exps {
		if pcJmpToElseif >= 0 {
			fi.fixSbx(pcJmpToElseif, fi.pc()-pcJmpToElseif)
		}

		oldRegs := fi.usedRegs
		a, _ := expToOpArg(fi, exp, ARG_REG)
		fi.resetRegs(oldRegs)

		line := lastLineOf(exp)
		fi.emitTest(line, a, 0)
		pcJmpToElseif = fi.emitJmp(line, 0)

		block := node.Blocks[i]
		cgBlockWithNewScope(fi, block, false)
		if i < len(node.Exps)-1 {
			pcJmpToEnds[i] = fi.emitJmp(block.LastLine, 0)
		} else {
			pcJmpToEnds[i] = pcJmpToElseif
		}
	}

	for _, pc := range pcJmpToEnds {
		fi.fixSbx(pc, fi.pc()-pc)
	}
}

func cgForNumStat(fi *funcInfo, node *ForNumStat) {
	forIndexVar := "(for index)"
	forLimitVar := "(for limit)"
	forStepVar := "(for step)"

	fi.enterScope(true)

	cgStat(fi, &LocalAssignStat{
		NameList: []string{forIndexVar, forLimitVar, forStepVar},
		ExpList:  []Exp{node.InitExp, node.LimitExp, node.StepExp},
	})
	fi.addLocVar(node.VarName, fi.pc()+2)

	a := fi.usedRegs - 4
	prepPc := fi.emitForPrep(node.LineOfDo, a, 0)
	cgBlock(fi, node.Block)
	loopPc := fi.emitForLoop(node.LineOfFor, a, 0)

	fi.fixSbx(prepPc, loopPc-prepPc-1)
	fi.fixSbx(loopPc, prepPc-loopPc)

	fi.exitScope(fi.pc())
	fi.fixEndPC(forIndexVar, 1)
	fi.fixEndPC(forLimitVar, 1)
	fi.fixEndPC(forStepVar, 1)
}

func cgForInStat(fi *funcInfo, node *ForInStat) {
	forGeneratorVar := "(for generator)"
	forStateVar := "(for state)"
	forControlVar := "(for control)"

	fi.enterScope(true)

	cgStat(fi, &LocalAssignStat{
		//LastLine: 0,
		NameList: []string{forGeneratorVar, forStateVar, forControlVar},
		ExpList:  node.ExpList,
	})
	for _, name := range node.NameList {
		fi.addLocVar(name, fi.pc()+2)
	}

	jmpToTFC := fi.emitJmp(node.LineOfDo, 0)
	cgBlock(fi, node.Block)
	fi.fixSbx(jmpToTFC, fi.pc()-jmpToTFC)

	line := lineOf(node.ExpList[0])
	rGenerator := fi.slotOfLocVar(forGeneratorVar)
	fi.emitTForCall(line, rGenerator, len(node.NameList))
	fi.emitTForLoop(line, rGenerator+2, jmpToTFC-fi.pc()-1)

	fi.exitScope(fi.pc() - 1)
	fi.fixEndPC(forGeneratorVar, 2)
	fi.fixEndPC(forStateVar, 2)
	fi.fixEndPC(forControlVar, 2)
}

func cgLocalAssignStat(fi *funcInfo, node *LocalAssignStat) {
	exps := removeTailNils(node.ExpList)
	nExps := len(exps)
	nNames := len(node.NameList)

	oldRegs := fi.usedRegs
	if nExps == nNames {
		for _, exp := range exps {
			a := fi.allocReg()
			cgExp(fi, exp, a, 1)
		}
	} else if nExps > nNames {
		for i, exp := range exps {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				cgExp(fi, exp, a, 0)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
	} else { // nNames > nExps
		multRet := false
		for i, exp := range exps {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				multRet = true
				n := nNames - nExps + 1
				cgExp(fi, exp, a, n)
				fi.allocRegs(n - 1)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
		if !multRet {
			n := nNames - nExps
			a := fi.allocRegs(n)
			fi.emitLoadNil(node.LastLine, a, n)
		}
	}

	fi.resetRegs(oldRegs)
	startPC := fi.pc() + 1
	for _, name := range node.NameList {
		fi.addLocVar(name, startPC)
	}
}

func cgAssignStat(fi *funcInfo, node *AssignStat) {
	exps := removeTailNils(node.ExpList)
	nExps := len(exps)
	nVars := len(node.VarList)

	ts := make([]int, nVars)
	ks := make([]int, nVars)
	vs := make([]int, nVars)
	oldRegs := fi.usedRegs

	for i, exp := range node.VarList {
		if taExp, ok := exp.(*TableAccessExp); ok {
			ts[i] = fi.allocReg()
			cgExp(fi, taExp.PrefixExp, ts[i], 1)
			ks[i] = fi.allocReg()
			cgExp(fi, taExp.KeyExp, ks[i], 1)
		} else {
			nameExp := exp.(*NameExp)
			if fi.slotOfLocVar(nameExp.Name) < 0 &&
				fi.indexOfUpval(nameExp.Name) < 0 {
				// global var
				ks[i] = -1
				if fi.indexOfConstant(nameExp.Name) > 0xFF {
					ks[i] = fi.allocReg()
				}
			}
		}
	}
	for i := 0; i < nVars; i++ {
		vs[i] = fi.usedRegs + i
	}

	if nExps >= nVars {
		for i, exp := range exps {
			a := fi.allocReg()
			if i >= nVars && i == nExps-1 && isVarargOrFuncCall(exp) {
				cgExp(fi, exp, a, 0)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
	} else { // nVars > nExps
		multRet := false
		for i, exp := range exps {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				multRet = true
				n := nVars - nExps + 1
				cgExp(fi, exp, a, n)
				fi.allocRegs(n - 1)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
		if !multRet {
			n := nVars - nExps
			a := fi.allocRegs(n)
			fi.emitLoadNil(node.LastLine, a, n)
		}
	}

	lastLine := node.LastLine
	for i, exp := range node.VarList {
		if nameExp, ok := exp.(*NameExp); ok {
			varName := nameExp.Name
			if a := fi.slotOfLocVar(varName); a >= 0 {
				fi.emitMove(lastLine, a, vs[i])
			} else if a := fi.indexOfUpval(varName); a >= 0 {
				fi.emitSetUpval(lastLine, a, vs[i])
			} else { // global var
				a := fi.indexOfUpval("_ENV")
				if ks[i] < 0 {
					b := 0x100 + fi.indexOfConstant(varName)
					fi.emitSetTabUp(lastLine, a, b, vs[i])
				} else {
					fi.emitSetTabUp(lastLine, a, ks[i], vs[i])
				}
			}
		} else {
			fi.emitSetTable(lastLine, ts[i], ks[i], vs[i])
		}
	}

	// todo
	fi.resetRegs(oldRegs)
}
