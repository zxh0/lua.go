package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

/*
         if
        (exp) ---.
        then     |jmp0
     ,-[block]   |
jmp1|   elseif <-'
    |   (exp) ---.
    |   then     |jmp2
    |,-[block]   |
jmp3|   elseif <-'
    |   (exp) ---.
    |   then     |jmp4
    |,-[block]   |
jmp5|   elseif <-'
    |   true
    |   then
    |,-[block]
    '-> end
*/
func (self *cg) ifStat(node *IfStat) {
	jmp2elseIfs := map[int]bool{}
	jmp2ends := map[int]bool{}

	for i := 0; i < len(node.Exps); i++ {
		if i > 0 {
			for pc, _ := range jmp2elseIfs {
				self.fixSbx(pc, self.pc0()-pc)
			}
			jmp2elseIfs = map[int]bool{} // clear map
		}

		self._cgIf(node, i, jmp2elseIfs, jmp2ends)
	}

	for pc, _ := range jmp2elseIfs {
		self.fixSbx(pc, self.pc0()-pc)
	}
	for pc, _ := range jmp2ends {
		self.fixSbx(pc, self.pc0()-pc)
	}
}

// todo: rename
func (self *cg) _cgIf(node *IfStat, i int,
	jmp2elseIfs, jmp2ends map[int]bool) {

	exp := node.Exps[i]
	block := node.Blocks[i]
	lineOfThen := node.Lines[i]

	if isExpTrue(exp) {
		switch x := exp.(type) {
		case *StringExp:
			self.indexOf(x.Str)
		}
	} else {
		if slot, ok := self.isLocVar(exp); ok {
			self.test(lineOfThen, slot, 0)
			pc := self.jmp(lineOfThen, 0)
			jmp2elseIfs[pc] = true
		} else if bexp, ok := exp.(*BinopExp); ok && bexp.Op == TOKEN_OP_AND {
			jmps := self.testLogicalAndExp(bexp, lineOfThen)
			for _, pc := range jmps {
				jmp2elseIfs[pc] = true
			}
		} else if bexp, ok := exp.(*BinopExp); ok && bexp.Op == TOKEN_OP_OR {
			jmp := self.testLogicalOrExp(bexp, lineOfThen)
			jmp2elseIfs[jmp] = true
		} else {
			tmp := self.allocTmp()
			self.testExp(exp, tmp) // todo
			if !isRelationalBinopExp(exp) {
				self.test(lineOfThen, tmp, 0)
			}
			self.freeTmp()
			pc := self.jmp(lineOfThen, 0)
			jmp2elseIfs[pc] = true
		}

	}

	self.block(block)
	if i < len(node.Exps)-1 {
		pc := self.jmp(block.LastLine, 0)
		jmp2ends[pc] = true
	}
}

func isRelationalBinopExp(exp Exp) bool {
	if binopExp, ok := exp.(*BinopExp); ok {
		switch binopExp.Op {
		case TOKEN_OP_EQ, TOKEN_OP_NE,
			TOKEN_OP_LT, TOKEN_OP_LE,
			TOKEN_OP_GT, TOKEN_OP_GE:
			return true
		}
	}
	return false
}
