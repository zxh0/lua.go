package codegen

import . "luago/compiler/ast"

const forIdxVar = "(for index)"
const forLmtVar = "(for limit)"
const forStpVar = "(for step)"

func (self *codeGen) forNumStat(node *ForNumStat) {
	self.enterScope()

	self.stat(&LocalAssignStat{
		LastLine: node.LineOfFor,
		NameList: []string{forIdxVar, forLmtVar, forStpVar},
		ExpList:  []Exp{node.InitExp, node.LimitExp, node.StepExp},
	})
	self.addLocVar(node.VarName, self.pc()+2)

	a := self.scope.stackSize - 3
	prepPc := self.emitForPrep(node.LineOfDo, a, 0)
	self.block(node.Block)
	loopPc := self.emitForLoop(node.LineOfFor, a, 0)

	self.fixSbx(prepPc, loopPc-prepPc-1)
	self.fixSbx(loopPc, prepPc-loopPc)

	self.exitScope(self.pc())
	self.fixEndPc(forIdxVar, 1)
	self.fixEndPc(forLmtVar, 1)
	self.fixEndPc(forStpVar, 1)
}
