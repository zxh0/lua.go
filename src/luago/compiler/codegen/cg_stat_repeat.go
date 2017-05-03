package codegen

import . "luago/compiler/ast"
import . "luago/lua/vm"

/*
 repeat
[block]<-.
 until   |jmp
 (exp) --'
*/
func (self *cg) repeatStat(node *RepeatStat) {
	if nilExp, ok := node.Exp.(*NilExp); ok {
		node.Exp = &FalseExp{nilExp.Line}
	}

	pc1 := self.pc()
	self.block(node.Block)
	if !isExpTrue(node.Exp) {
		//self.exp(node.Exp, STAT_REPEAT, 0)

		line := LineOfExp(node.Exp)
		pc2 := self.inst(line, OP_TEST, 0, 0, 0) // todo
		self.inst(line, OP_JMP, 0, pc1-pc2-2, 0) // todo
		self.freeTmp()
	}
}

// todo
func LineOfExp(exp Exp) int {
	switch x := exp.(type) {
	case *NilExp:
		return x.Line
	case *TrueExp:
		return x.Line
	case *FalseExp:
		return x.Line
	case *VarargExp:
		return x.Line
	case *IntegerExp:
		return x.Line
	case *FuncDefExp:
		return x.Line
	case *FloatExp:
		return x.Line
	case *StringExp:
		return x.Line
	case *BinopExp: // todo
		return x.Line
	case *UnopExp: // todo
		return x.Line
	case *TableConstructorExp:
		return 0 // todo
	case *NameExp:
		return x.Line
	case *BracketsExp:
		return x.Line
	case *FuncCallExp:
		return x.Line
	}

	panic("todo!")
}