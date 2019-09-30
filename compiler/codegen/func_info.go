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
	codeBuf
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
	line      int
	lastLine  int
	numParams int
	isVararg  bool
}

func newFuncInfo(parent *funcInfo, fd *FuncDefExp) *funcInfo {
	return &funcInfo{
		codeBuf: codeBuf{
			insts:    make([]uint32, 0, 8),
			lineNums: make([]uint32, 0, 8),
		},
		parent:    parent,
		subFuncs:  []*funcInfo{},
		locVars:   make([]*locVarInfo, 0, 8),
		locNames:  map[string]*locVarInfo{},
		upvalues:  map[string]upvalInfo{},
		constants: map[interface{}]int{},
		labels:    map[string]labelInfo{},
		gotos:     nil,
		breaks:    make([][]int, 1),
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
