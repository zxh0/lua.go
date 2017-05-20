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

	jmpToTFC := self.jmp(node.LineOfDo, 0)
	self.block(node.Block)
	self.fixSbx(jmpToTFC, self.pc()-jmpToTFC)

	line := lineOfExp(node.ExpList[0])
	slotOfGeneratorVar := self.slotOf(forGeneratorVar)
	self.tForCall(line, slotOfGeneratorVar, len(node.NameList))
	self.tForLoop(line, slotOfGeneratorVar+2, jmpToTFC-self.pc()-1)

	self.exitScope(self.pc() + 1)
	for _, name := range node.NameList {
		self.fixEndPcOfLocVar(name, -2)
	}
}
