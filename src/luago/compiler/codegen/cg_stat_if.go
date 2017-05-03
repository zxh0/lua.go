package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

/*
         if
        (exp) ---.
        then     |jmp0
    .--[block]   |
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
	jmps := make([]int, len(node.Exps)*2)
	for i, exp := range node.Exps {
		block := node.Blocks[i]
		line := node.Lines[i]
		isLastExp := (i == len(node.Exps)-1)
		jmps[i*2], jmps[i*2+1] =
			_cgIf(self, exp, block, line, isLastExp)
	}

	// fix jmps
	pc := self.pc() - 1 // todo
	for i := 0; i < len(jmps)/2; i++ {
		if jmps[i*2] > 0 {
			if i < len(jmps)/2-1 {
				self.fixSbx(jmps[i*2], jmps[i*2+1]-jmps[i*2])
			} else {
				self.fixSbx(jmps[i*2], pc-jmps[i*2]-1)
			}
		}
		if jmps[i*2+1] > 0 {
			self.fixSbx(jmps[i*2+1], pc-jmps[i*2+1]-1)
		}
	}
}

// todo
func _cgIf(cg *cg, exp Exp, block *Block,
	lineOfThen int, isLastExp bool) (jmp1, jmp2 int) {

	if isExpTrue(exp) {
		switch x := exp.(type) {
		case *StringExp:
			cg.indexOf(x.Str)
		}
	} else {
		tmp := cg.allocTmp()
		if isRelationalBinopExp(exp) {
			cg.exp(exp, tmp, 0)
		} else {
			cg.exp(exp, tmp, 1)
			cg.test(lineOfThen, tmp, 0)
		}
		cg.freeTmp()

		jmp1 = cg.jmp(lineOfThen, 0)
	}

	cg.block(block)
	if !isLastExp {
		jmp2 = cg.jmp(block.LastLine, 0)
	}

	return
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
