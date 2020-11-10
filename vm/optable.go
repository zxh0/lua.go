package vm

import (
	"github.com/zxh0/lua.go/api"
)

type opFn = func(i Instruction, vm api.LuaVM)

type opInfo struct {
	flagA  byte // instruction set register A
	flagT  byte // operator is a test (next instruction must be a jump)
	flagIT byte // instruction uses 'L->top' set by previous instruction (when B == 0)
	flagOT byte // instruction sets 'L->top' for next instruction (when C == 0)
	flagMM byte // instruction is an MM instruction (call a metamethod)
	opMode byte // instruction format
	name   string
	action opFn
}

var opTable = []opInfo{
	/*      MM OT IT T  A  mode   name      action   */
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_MOVE", move),
	_opInfo(0, 0, 0, 0, 1, iAsBx, "OP_LOADI", loadI),
	_opInfo(0, 0, 0, 0, 1, iAsBx, "OP_LOADF", loadF),
	_opInfo(0, 0, 0, 0, 1, iABx, "OP_LOADK", loadK),
	_opInfo(0, 0, 0, 0, 1, iABx, "OP_LOADKX", loadKx),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_LOADFALSE", loadFalse),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_LFALSESKIP", lFalseSkip),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_LOADTRUE", loadTrue),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_LOADNIL", loadNil),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_GETUPVAL", getUpval),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_SETUPVAL", setUpval),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_GETTABUP", getTabUp),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_GETTABLE", getTable),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_GETI", getI),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_GETFIELD", getField),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_SETTABUP", setTabUp),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_SETTABLE", setTable),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_SETI", setI),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_SETFIELD", setField),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_NEWTABLE", newTable),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_SELF", self),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_ADDI", addI),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_ADDK", addK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_SUBK", subK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_MULK", mulK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_MODK", modK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_POWK", powK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_DIVK", divK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_IDIVK", idivK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_BANDK", bandK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_BORK", borK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_BXORK", bxorK),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_SHRI", shrI),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_SHLI", shlI),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_ADD", add),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_SUB", sub),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_MUL", mul),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_MOD", mod),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_POW", pow),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_DIV", div),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_IDIV", idiv),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_BAND", band),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_BOR", bor),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_BXOR", bxor),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_SHL", shl),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_SHR", shr),
	_opInfo(1, 0, 0, 0, 0, iABC, "OP_MMBIN", _todo),
	_opInfo(1, 0, 0, 0, 0, iABC, "OP_MMBINI", _todo),
	_opInfo(1, 0, 0, 0, 0, iABC, "OP_MMBINK", _todo),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_UNM", unm),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_BNOT", bnot),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_NOT", not),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_LEN", length),
	_opInfo(0, 0, 0, 0, 1, iABC, "OP_CONCAT", concat),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_CLOSE", closeUV),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_TBC", _todo),
	_opInfo(0, 0, 0, 0, 0, isJ, "OP_JMP", jmp),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_EQ", eq),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_LT", lt),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_LE", le),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_EQK", eqK),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_EQI", eqI),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_LTI", ltI),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_LEI", leI),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_GTI", gtI),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_GEI", geI),
	_opInfo(0, 0, 0, 1, 0, iABC, "OP_TEST", test),
	_opInfo(0, 0, 0, 1, 1, iABC, "OP_TESTSET", testSet),
	_opInfo(0, 1, 1, 0, 1, iABC, "OP_CALL", call),
	_opInfo(0, 1, 1, 0, 1, iABC, "OP_TAILCALL", tailCall),
	_opInfo(0, 0, 1, 0, 0, iABC, "OP_RETURN", _return),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_RETURN0", _todo),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_RETURN1", _todo),
	_opInfo(0, 0, 0, 0, 1, iABx, "OP_FORLOOP", forLoop),
	_opInfo(0, 0, 0, 0, 1, iABx, "OP_FORPREP", forPrep),
	_opInfo(0, 0, 0, 0, 0, iABx, "OP_TFORPREP", tForPrep),
	_opInfo(0, 0, 0, 0, 0, iABC, "OP_TFORCALL", tForCall),
	_opInfo(0, 0, 0, 0, 1, iABx, "OP_TFORLOOP", tForLoop),
	_opInfo(0, 0, 1, 0, 0, iABC, "OP_SETLIST", setList),
	_opInfo(0, 0, 0, 0, 1, iABx, "OP_CLOSURE", closure),
	_opInfo(0, 1, 0, 0, 1, iABC, "OP_VARARG", vararg),
	_opInfo(0, 0, 1, 0, 1, iABC, "OP_VARARGPREP", varargPrep),
	_opInfo(0, 0, 0, 0, 0, iAx, "OP_EXTRAARG", _todo),
}

func _opInfo(mm, ot, it, t, a, mode byte, name string, action opFn) opInfo {
	return opInfo{a, t, it, ot, mm, mode, name, action}
}
func _todo(i Instruction, vm api.LuaVM) {
	panic("TODO~")
}
