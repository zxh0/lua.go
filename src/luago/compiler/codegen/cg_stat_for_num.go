package codegen

import . "luago/compiler/ast"
import . "luago/lua/vm"

const forIdxVar = "(for index)"
const forLmtVar = "(for limit)"
const forStpVar = "(for step)"

func (self *cg) forNumStat(node *ForNumStat) {
	self.enterScope()

	self.stat(&LocalAssignStat{
		LastLine: node.LineOfFor,
		NameList: []string{forIdxVar, forLmtVar, forStpVar},
		ExpList:  []Exp{node.InitExp, node.LimitExp, node.StepExp},
	})
	self.addLocVar(node.VarName, self.pc())

	a := self.scope.stackSize - 3
	prepPc := self.inst(node.LineOfDo, OP_FORPREP, a, 0, 0)
	self.block(node.Block)
	loopPc := self.inst(node.LineOfFor, OP_FORLOOP, a, 0, 0)

	self.fixSbx(prepPc, loopPc-prepPc-1)
	self.fixSbx(loopPc, prepPc-loopPc)

	self.exitScope(self.pc()-1)
	self.fixEndPcOfIdxVar(node.VarName)
}

func (self *cg) fixEndPcOfIdxVar(name string) {
	for i := len(self.scope.locVars)-1; i > 0; i-- {
		locVar := self.scope.locVars[i]
		if locVar.name == name {
			locVar.endPc -= 1
			return
		}
	}
}
