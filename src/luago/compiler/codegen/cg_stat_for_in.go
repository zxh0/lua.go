package codegen

import . "luago/compiler/ast"

const forGeneratorVar = "(for generator)"
const forStateVar = "(for state)"
const forControlVar = "(for control)"

func (self *cg) forInStat(node *ForInStat) {
	self.enterScope()

	self.stat(&LocalAssignStat{
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

	jmpToTFC := self.jmp(node.LineOfDo, 0)
	self.block(node.Block)
	self.fixSbx(jmpToTFC, self.pc()-jmpToTFC)

	line := lineOfExp(node.ExpList[0])
	slotOfGeneratorVar := self.slotOf(forGeneratorVar)
	self.tForCall(line, slotOfGeneratorVar, len(node.NameList))
	self.tForLoop(line, slotOfGeneratorVar+2, jmpToTFC-self.pc()-1)

	self.exitScope(self.pc() - 1)
	self.fixEndPc(forGeneratorVar, 2)
	self.fixEndPc(forStateVar, 2)
	self.fixEndPc(forControlVar, 2)
}
