package codegen

import . "luago/compiler/ast"

func (self *cg) stat(stat Stat) {
	switch node := stat.(type) {
	case *IfStat:
		self.ifStat(node)
	case *WhileStat:
		self.whileStat(node)
	case *RepeatStat:
		self.repeatStat(node)
	case *ForNumStat:
		self.forNumStat(node)
	case FuncCallStat:
		self.funcCallStat(node)
	case *LocalAssignStat:
		self.localAssignStat(node)
	case *AssignStat:
		self.assignStat(node)
	case *BreakStat:
		// todo
	case DoStat:
		self.block(node)
	default:
		panic("todo: stat!")
	}
}

func (self *cg) funcCallStat(stat FuncCallStat) {
	fcExp := (*FuncCallExp)(stat)
	tmp := self.allocTmp()
	self.exp(fcExp, tmp, 0)
	self.freeTmp()
}
