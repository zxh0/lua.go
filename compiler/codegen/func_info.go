package codegen

import (
	"fmt"

	. "github.com/zxh0/lua.go/compiler/ast"
	. "github.com/zxh0/lua.go/compiler/lexer"
	. "github.com/zxh0/lua.go/vm"
)

var arithAndBitwiseBinops = map[int]int{
	TOKEN_OP_ADD:  OP_ADD,
	TOKEN_OP_SUB:  OP_SUB,
	TOKEN_OP_MUL:  OP_MUL,
	TOKEN_OP_MOD:  OP_MOD,
	TOKEN_OP_POW:  OP_POW,
	TOKEN_OP_DIV:  OP_DIV,
	TOKEN_OP_IDIV: OP_IDIV,
	TOKEN_OP_BAND: OP_BAND,
	TOKEN_OP_BOR:  OP_BOR,
	TOKEN_OP_BXOR: OP_BXOR,
	TOKEN_OP_SHL:  OP_SHL,
	TOKEN_OP_SHR:  OP_SHR,
}

type upvalInfo struct {
	locVarSlot int
	upvalIndex int
	index      int
}

type locVarInfo struct {
	prev     *locVarInfo
	name     string
	scopeLv  int
	slot     int
	startPC  int
	endPC    int
	captured bool
}

type labelInfo struct {
	line    int
	pc      int
	scopeLv int
}

type gotoInfo struct {
	jmpPC   int
	scopeLv int
	label   string
	pending bool
}

type funcInfo struct {
	parent    *funcInfo
	subFuncs  []*funcInfo
	usedRegs  int
	maxRegs   int
	scopeLv   int
	locVars   []*locVarInfo
	locNames  map[string]*locVarInfo
	upvalues  map[string]upvalInfo
	constants map[interface{}]int
	labels    map[string]labelInfo
	gotos     []*gotoInfo
	breaks    [][]int
	insts     []uint32
	lineNums  []uint32
	line      int
	lastLine  int
	numParams int
	isVararg  bool
}

func newFuncInfo(parent *funcInfo, fd *FuncDefExp) *funcInfo {
	return &funcInfo{
		parent:    parent,
		subFuncs:  []*funcInfo{},
		locVars:   make([]*locVarInfo, 0, 8),
		locNames:  map[string]*locVarInfo{},
		upvalues:  map[string]upvalInfo{},
		constants: map[interface{}]int{},
		labels:    map[string]labelInfo{},
		gotos:     nil,
		breaks:    make([][]int, 1),
		insts:     make([]uint32, 0, 8),
		lineNums:  make([]uint32, 0, 8),
		line:      fd.Line,
		lastLine:  fd.LastLine,
		numParams: len(fd.ParList),
		isVararg:  fd.IsVararg,
	}
}

/* constants */

func (fi *funcInfo) indexOfConstant(k interface{}) int {
	if idx, found := fi.constants[k]; found {
		return idx
	}

	idx := len(fi.constants)
	fi.constants[k] = idx
	return idx
}

/* registers */

func (fi *funcInfo) allocReg() int {
	fi.usedRegs++
	if fi.usedRegs >= 255 {
		panic("function or expression needs too many registers")
	}
	if fi.usedRegs > fi.maxRegs {
		fi.maxRegs = fi.usedRegs
	}
	return fi.usedRegs - 1
}

func (fi *funcInfo) freeReg() {
	if fi.usedRegs <= 0 {
		panic("usedRegs <= 0 !")
	}
	fi.usedRegs--
}

func (fi *funcInfo) allocRegs(n int) int {
	if n <= 0 {
		panic("n <= 0 !")
	}
	for i := 0; i < n; i++ {
		fi.allocReg()
	}
	return fi.usedRegs - n
}

func (fi *funcInfo) freeRegs(n int) {
	if n < 0 {
		panic("n < 0 !")
	}
	for i := 0; i < n; i++ {
		fi.freeReg()
	}
}

/* lexical scope */

func (fi *funcInfo) enterScope(breakable bool) {
	fi.scopeLv++
	if breakable {
		fi.breaks = append(fi.breaks, []int{})
	} else {
		fi.breaks = append(fi.breaks, nil)
	}
}

func (fi *funcInfo) exitScope(endPC int) {
	pendingBreakJmps := fi.breaks[len(fi.breaks)-1]
	fi.breaks = fi.breaks[:len(fi.breaks)-1]

	a := fi.getJmpArgA()
	for _, pc := range pendingBreakJmps {
		sBx := fi.pc() - pc
		i := (sBx+MAXARG_sBx)<<14 | a<<6 | OP_JMP
		fi.insts[pc] = uint32(i)
	}

	fi.fixGotoJmps()

	fi.scopeLv--
	for _, locVar := range fi.locNames {
		if locVar.scopeLv > fi.scopeLv { // out of scope
			locVar.endPC = endPC
			fi.removeLocVar(locVar)
		}
	}
}

func (fi *funcInfo) removeLocVar(locVar *locVarInfo) {
	fi.freeReg()
	if locVar.prev == nil {
		delete(fi.locNames, locVar.name)
	} else if locVar.prev.scopeLv == locVar.scopeLv {
		fi.removeLocVar(locVar.prev)
	} else {
		fi.locNames[locVar.name] = locVar.prev
	}
}

func (fi *funcInfo) addLocVar(name string, startPC int) int {
	newVar := &locVarInfo{
		name:    name,
		prev:    fi.locNames[name],
		scopeLv: fi.scopeLv,
		slot:    fi.allocReg(),
		startPC: startPC,
		endPC:   0,
	}

	fi.locVars = append(fi.locVars, newVar)
	fi.locNames[name] = newVar

	return newVar.slot
}

func (fi *funcInfo) slotOfLocVar(name string) int {
	if locVar, found := fi.locNames[name]; found {
		return locVar.slot
	}
	return -1
}

func (fi *funcInfo) addBreakJmp(pc int) {
	for i := fi.scopeLv; i >= 0; i-- {
		if fi.breaks[i] != nil { // breakable
			fi.breaks[i] = append(fi.breaks[i], pc)
			return
		}
	}

	panic("<break> at line ? not inside a loop!")
}

/* upvalues */

func (fi *funcInfo) indexOfUpval(name string) int {
	if upval, ok := fi.upvalues[name]; ok {
		return upval.index
	}
	if fi.parent != nil {
		if locVar, found := fi.parent.locNames[name]; found {
			idx := len(fi.upvalues)
			fi.upvalues[name] = upvalInfo{locVar.slot, -1, idx}
			locVar.captured = true
			return idx
		}
		if uvIdx := fi.parent.indexOfUpval(name); uvIdx >= 0 {
			idx := len(fi.upvalues)
			fi.upvalues[name] = upvalInfo{-1, uvIdx, idx}
			return idx
		}
	}
	return -1
}

func (fi *funcInfo) closeOpenUpvals(line int) {
	a := fi.getJmpArgA()
	if a > 0 {
		fi.emitJmp(line, a, 0)
	}
}

func (fi *funcInfo) getJmpArgA() int {
	hasCapturedLocVars := false
	minSlotOfLocVars := fi.maxRegs
	for _, locVar := range fi.locNames {
		if locVar.scopeLv == fi.scopeLv {
			for v := locVar; v != nil && v.scopeLv == fi.scopeLv; v = v.prev {
				if v.captured {
					hasCapturedLocVars = true
				}
				if v.slot < minSlotOfLocVars && v.name[0] != '(' {
					minSlotOfLocVars = v.slot
				}
			}
		}
	}
	if hasCapturedLocVars {
		return minSlotOfLocVars + 1
	} else {
		return 0
	}
}

/* labels */

func (fi *funcInfo) addLabel(label string, line int) {
	key := fmt.Sprintf("%s@%d", label, fi.scopeLv)
	if labelInfo, ok := fi.labels[key]; ok {
		panic(fmt.Sprintf("label '%s' already defined on line %d",
			label, labelInfo.line))
	}
	fi.labels[key] = labelInfo{line, fi.pc() + 1, fi.scopeLv}
}

func (fi *funcInfo) addGoto(jmpPC, scopeLv int, label string) {
	fi.gotos = append(fi.gotos, &gotoInfo{jmpPC, scopeLv, label, false})
}

func (fi *funcInfo) fixGotoJmps() {
	for i, gotoInfo := range fi.gotos {
		if gotoInfo == nil || gotoInfo.scopeLv < fi.scopeLv {
			continue
		}
		if gotoInfo.scopeLv == fi.scopeLv && gotoInfo.pending {
			continue
		}

		dstPC := fi.getGotoDst(gotoInfo.label)
		if dstPC >= 0 {
			if dstPC > gotoInfo.jmpPC && dstPC < fi.pc() {
				for _, locVar := range fi.locNames {
					if locVar.startPC > gotoInfo.jmpPC && locVar.startPC <= dstPC {
						panic(fmt.Sprintf("<goto %s> at line %d jumps into the scope of local '%s'",
							gotoInfo.label, fi.lineNums[gotoInfo.jmpPC], locVar.name))
					}
				}
			}

			a := 0
			for _, locVar := range fi.locVars {
				if locVar.startPC > dstPC {
					a = locVar.slot + 1
					break
				}
			}

			sBx := dstPC - gotoInfo.jmpPC - 1
			inst := (sBx+MAXARG_sBx)<<14 | a<<6 | OP_JMP
			fi.insts[gotoInfo.jmpPC] = uint32(inst)
			fi.gotos[i] = nil
		} else if fi.scopeLv == 0 {
			panic(fmt.Sprintf("no visible label '%s' for <goto> at line %d",
				gotoInfo.label, fi.lineNums[gotoInfo.jmpPC]))
		} else {
			gotoInfo.pending = true
		}
	}
	for key, labelInfo := range fi.labels {
		if labelInfo.scopeLv == fi.scopeLv {
			delete(fi.labels, key)
		}
	}
}

func (fi *funcInfo) getGotoDst(label string) int {
	for i := fi.scopeLv; i >= 0; i-- {
		key := fmt.Sprintf("%s@%d", label, i)
		if labelInfo, ok := fi.labels[key]; ok {
			return labelInfo.pc
		}
	}
	return -1
}

/* code */

func (fi *funcInfo) pc() int {
	return len(fi.insts) - 1
}

func (fi *funcInfo) fixSbx(pc, sBx int) {
	if sBx > 0 && sBx > MAXARG_sBx || sBx < 0 && -sBx > MAXARG_sBx {
		panic("control structure too long")
	}

	i := fi.insts[pc]
	i = i << 18 >> 18                  // clear sBx
	i = i | uint32(sBx+MAXARG_sBx)<<14 // reset sBx
	fi.insts[pc] = i
}

// todo: rename?
func (fi *funcInfo) fixEndPC(name string, delta int) {
	for i := len(fi.locVars) - 1; i >= 0; i-- {
		locVar := fi.locVars[i]
		if locVar.name == name {
			locVar.endPC += delta
			return
		}
	}
}

func (fi *funcInfo) emitABC(line, opcode, a, b, c int) {
	i := b<<23 | c<<14 | a<<6 | opcode
	fi.insts = append(fi.insts, uint32(i))
	fi.lineNums = append(fi.lineNums, uint32(line))
}

func (fi *funcInfo) emitABx(line, opcode, a, bx int) {
	i := bx<<14 | a<<6 | opcode
	fi.insts = append(fi.insts, uint32(i))
	fi.lineNums = append(fi.lineNums, uint32(line))
}

func (fi *funcInfo) emitAsBx(line, opcode, a, sBx int) {
	i := (sBx+MAXARG_sBx)<<14 | a<<6 | opcode
	fi.insts = append(fi.insts, uint32(i))
	fi.lineNums = append(fi.lineNums, uint32(line))
}

func (fi *funcInfo) emitAx(line, opcode, ax int) {
	i := ax<<6 | opcode
	fi.insts = append(fi.insts, uint32(i))
	fi.lineNums = append(fi.lineNums, uint32(line))
}

// r[a] = r[b]
func (fi *funcInfo) emitMove(line, a, b int) {
	fi.emitABC(line, OP_MOVE, a, b, 0)
}

// r[a], r[a+1], ..., r[a+b] = nil
func (fi *funcInfo) emitLoadNil(line, a, n int) {
	fi.emitABC(line, OP_LOADNIL, a, n-1, 0)
}

// r[a] = (bool)b; if (c) pc++
func (fi *funcInfo) emitLoadBool(line, a, b, c int) {
	fi.emitABC(line, OP_LOADBOOL, a, b, c)
}

// r[a] = kst[bx]
func (fi *funcInfo) emitLoadK(line, a int, k interface{}) {
	idx := fi.indexOfConstant(k)
	if idx < (1 << 18) {
		fi.emitABx(line, OP_LOADK, a, idx)
	} else {
		fi.emitABx(line, OP_LOADKX, a, 0)
		fi.emitAx(line, OP_EXTRAARG, idx)
	}
}

// r[a], r[a+1], ..., r[a+b-2] = vararg
func (fi *funcInfo) emitVararg(line, a, n int) {
	fi.emitABC(line, OP_VARARG, a, n+1, 0)
}

// r[a] = emitClosure(proto[bx])
func (fi *funcInfo) emitClosure(line, a, bx int) {
	fi.emitABx(line, OP_CLOSURE, a, bx)
}

// r[a] = {}
func (fi *funcInfo) emitNewTable(line, a, nArr, nRec int) {
	fi.emitABC(line, OP_NEWTABLE,
		a, Int2fb(nArr), Int2fb(nRec))
}

// r[a][(c-1)*FPF+i] := r[a+i], 1 <= i <= b
func (fi *funcInfo) emitSetList(line, a, b, c int) {
	fi.emitABC(line, OP_SETLIST, a, b, c)
}

// r[a] := r[b][rk(c)]
func (fi *funcInfo) emitGetTable(line, a, b, c int) {
	fi.emitABC(line, OP_GETTABLE, a, b, c)
}

// r[a][rk(b)] = rk(c)
func (fi *funcInfo) emitSetTable(line, a, b, c int) {
	fi.emitABC(line, OP_SETTABLE, a, b, c)
}

// r[a] = upval[b]
func (fi *funcInfo) emitGetUpval(line, a, b int) {
	fi.emitABC(line, OP_GETUPVAL, a, b, 0)
}

// upval[b] = r[a]
func (fi *funcInfo) emitSetUpval(line, a, b int) {
	fi.emitABC(line, OP_SETUPVAL, a, b, 0)
}

// r[a] = upval[b][rk(c)]
func (fi *funcInfo) emitGetTabUp(line, a, b, c int) {
	fi.emitABC(line, OP_GETTABUP, a, b, c)
}

// upval[a][rk(b)] = rk(c)
func (fi *funcInfo) emitSetTabUp(line, a, b, c int) {
	fi.emitABC(line, OP_SETTABUP, a, b, c)
}

// r[a], ..., r[a+c-2] = r[a](r[a+1], ..., r[a+b-1])
func (fi *funcInfo) emitCall(line, a, nArgs, nRet int) {
	fi.emitABC(line, OP_CALL, a, nArgs+1, nRet+1)
}

// return r[a](r[a+1], ... ,r[a+b-1])
func (fi *funcInfo) emitTailCall(line, a, nArgs int) {
	fi.emitABC(line, OP_TAILCALL, a, nArgs+1, 0)
}

// return r[a], ... ,r[a+b-2]
func (fi *funcInfo) emitReturn(line, a, n int) {
	fi.emitABC(line, OP_RETURN, a, n+1, 0)
}

// r[a+1] := r[b]; r[a] := r[b][rk(c)]
func (fi *funcInfo) emitSelf(line, a, b, c int) {
	fi.emitABC(line, OP_SELF, a, b, c)
}

// pc+=sBx; if (a) close all upvalues >= r[a - 1]
func (fi *funcInfo) emitJmp(line, a, sBx int) int {
	fi.emitAsBx(line, OP_JMP, a, sBx)
	return len(fi.insts) - 1
}

// if not (r[a] <=> c) then pc++
func (fi *funcInfo) emitTest(line, a, c int) {
	fi.emitABC(line, OP_TEST, a, 0, c)
}

// if (r[b] <=> c) then r[a] := r[b] else pc++
func (fi *funcInfo) emitTestSet(line, a, b, c int) {
	fi.emitABC(line, OP_TESTSET, a, b, c)
}

func (fi *funcInfo) emitForPrep(line, a, sBx int) int {
	fi.emitAsBx(line, OP_FORPREP, a, sBx)
	return len(fi.insts) - 1
}

func (fi *funcInfo) emitForLoop(line, a, sBx int) int {
	fi.emitAsBx(line, OP_FORLOOP, a, sBx)
	return len(fi.insts) - 1
}

func (fi *funcInfo) emitTForCall(line, a, c int) {
	fi.emitABC(line, OP_TFORCALL, a, 0, c)
}

func (fi *funcInfo) emitTForLoop(line, a, sBx int) {
	fi.emitAsBx(line, OP_TFORLOOP, a, sBx)
}

// r[a] = op r[b]
func (fi *funcInfo) emitUnaryOp(line, op, a, b int) {
	switch op {
	case TOKEN_OP_NOT:
		fi.emitABC(line, OP_NOT, a, b, 0)
	case TOKEN_OP_BNOT:
		fi.emitABC(line, OP_BNOT, a, b, 0)
	case TOKEN_OP_LEN:
		fi.emitABC(line, OP_LEN, a, b, 0)
	case TOKEN_OP_UNM:
		fi.emitABC(line, OP_UNM, a, b, 0)
	}
}

// r[a] = rk[b] op rk[c]
// arith & bitwise & relational
func (fi *funcInfo) emitBinaryOp(line, op, a, b, c int) {
	if opcode, found := arithAndBitwiseBinops[op]; found {
		fi.emitABC(line, opcode, a, b, c)
	} else {
		switch op {
		case TOKEN_OP_EQ:
			fi.emitABC(line, OP_EQ, 1, b, c)
		case TOKEN_OP_NE:
			fi.emitABC(line, OP_EQ, 0, b, c)
		case TOKEN_OP_LT:
			fi.emitABC(line, OP_LT, 1, b, c)
		case TOKEN_OP_GT:
			fi.emitABC(line, OP_LT, 1, c, b)
		case TOKEN_OP_LE:
			fi.emitABC(line, OP_LE, 1, b, c)
		case TOKEN_OP_GE:
			fi.emitABC(line, OP_LE, 1, c, b)
		}
		fi.emitJmp(line, 0, 1)
		fi.emitLoadBool(line, a, 0, 1)
		fi.emitLoadBool(line, a, 1, 0)
	}
}
