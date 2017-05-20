package codegen

import . "luago/compiler/ast"

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
				self.fixSbx(pc, self.pc()-pc)
			}
			jmp2elseIfs = map[int]bool{} // clear map
		}

		self._cgIf(node, i, jmp2elseIfs, jmp2ends)
	}

	for pc, _ := range jmp2elseIfs {
		self.fixSbx(pc, self.pc()-pc)
	}
	for pc, _ := range jmp2ends {
		self.fixSbx(pc, self.pc()-pc)
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
		pendingJmps := self.testExp(exp, lineOfThen)
		for _, pc := range pendingJmps {
			jmp2elseIfs[pc] = true
		}
	}

	self.blockWithNewScope(block)
	if i < len(node.Exps)-1 {
		pc := self.jmp(block.LastLine, 0)
		jmp2ends[pc] = true
	}
}
