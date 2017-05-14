package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// todo: rename
type expListNode struct {
	exp     Exp
	parent  Exp // todo: rename
	op      int
	line    int
	lv      int
	startPc int
	jmpPc   int
	next    *expListNode
	jmpTo   *expListNode
}

func logicalBinopExpToList(exp *BinopExp) *expListNode {
	head, _ := _logicalBinopExpToList(exp, 0, exp)
	for node := head; node != nil; node = node.next {
		if node.next != nil && node.next.parent == node.parent {
			for next := node.next; next != nil; next = next.next {
				if next.lv <= node.lv && next.parent != node.parent {
					node.jmpTo = next
					break
				}
			}
		}
	}
	return head
}

func _logicalBinopExpToList(exp *BinopExp, lv int, parent Exp) (head, tail *expListNode) {
	var head1, tail1, head2, tail2 *expListNode

	if exp1, ok := castToLogicalBinopExp(exp.Exp1); ok {
		if exp1.Op != exp.Op {
			head1, tail1 = _logicalBinopExpToList(exp1, lv+1, exp1)
		} else {
			head1, tail1 = _logicalBinopExpToList(exp1, lv, parent)
		}
	} else {
		head1 = &expListNode{exp: exp.Exp1, lv: lv, parent: parent}
		tail1 = head1
	}

	if exp2, ok := castToLogicalBinopExp(exp.Exp2); ok {
		if exp2.Op != exp.Op {
			head2, tail2 = _logicalBinopExpToList(exp2, lv+1, exp2)
		} else {
			head2, tail2 = _logicalBinopExpToList(exp2, lv, parent)
		}
	} else {
		head2 = &expListNode{exp: exp.Exp2, lv: lv, parent: parent}
		tail2 = head2
	}

	tail1.op = exp.Op
	tail1.line = exp.Line
	tail1.next = head2

	return head1, tail2
}

func castToLogicalBinopExp(exp Exp) (*BinopExp, bool) {
	exp = stripParans(exp)
	if bexp, ok := exp.(*BinopExp); ok {
		if bexp.Op == TOKEN_OP_AND || bexp.Op == TOKEN_OP_OR {
			return bexp, true
		}
	}
	return nil, false
}

// (exp) => exp
func stripParans(exp Exp) Exp {
	if paransExp, ok := exp.(*ParensExp); ok {
		return paransExp.Exp
	} else {
		return exp
	}
}
