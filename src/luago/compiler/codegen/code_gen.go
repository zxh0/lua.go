package codegen

import . "luago/binchunk"
import . "luago/compiler/ast"

type codeGen struct {
	scope  *scope
	insts  insts
	protos []*Prototype
}

func newCodeGen(parentScope *scope) *codeGen {
	return &codeGen{
		scope:  newScope(parentScope),
		insts:  make([]instruction, 0, 8),
		protos: []*Prototype{},
	}
}

func (self *codeGen) pc() int {
	return len(self.insts) - 1
}

func (self *codeGen) emit(line, opcode, a, b, c int) int {
	i := instruction{line, opcode, a, b, c}
	self.insts = append(self.insts, i)
	return len(self.insts) - 1
}

func (self *codeGen) fixA(pc, a int) {
	i := self.insts[pc]
	i.a = a
	self.insts[pc] = i
}

func (self *codeGen) fixSbx(pc, sbx int) {
	i := self.insts[pc]
	i.b = sbx
	self.insts[pc] = i
}

func (self *codeGen) enterScope(breakable bool) {
	self.scope.incrLevel()
	if breakable {
		self.scope.markBreakable()
	}
}
func (self *codeGen) exitScope(endPc int) {
	pendingBreakJmps := self.scope.decrLevel(endPc)
	for _, pc := range pendingBreakJmps {
		self.fixSbx(pc, self.pc()-pc)
	}
}
func (self *codeGen) addLocVar(name string, startPc int) int {
	return self.scope.addLocVar(name, startPc)
}
func (self *codeGen) slotOf(name string) int {
	return self.scope.slotOf(name)
}
func (self *codeGen) lookupUpval(name string) int {
	return self.scope.lookupUpval(name)
}

func (self *codeGen) addBreakJmp(pc int) {
	self.scope.addBreakJmp(pc)
}

// todo: rename?
func (self *codeGen) fixEndPc(name string, delta int) {
	for i := len(self.scope.locVars) - 1; i >= 0; i-- {
		locVar := self.scope.locVars[i]
		if locVar.name == name {
			locVar.endPc += delta
			return
		}
	}
}

func (self *codeGen) newTmpAllocator(a int) *tmpAllocator {
	if self.isTmpVar(a) {
		return &tmpAllocator{self.scope, a, 0}
	} else {
		return &tmpAllocator{self.scope, -1, 0}
	}
}
func (self *codeGen) allocTmps(n int) int {
	return self.scope.allocTmps(n)
}
func (self *codeGen) allocTmp() int {
	return self.scope.allocTmp()
}
func (self *codeGen) freeTmps(n int) {
	self.scope.freeTmps(n)
}
func (self *codeGen) freeTmp() {
	self.scope.freeTmp()
}
func (self *codeGen) isLocVarSlot(slot int) bool {
	return self.scope.isLocVarSlot(slot)
}
func (self *codeGen) isTmpVar(slot int) bool {
	return self.scope.isTmpVar(slot)
}
func (self *codeGen) isGlobalVar(name string) (int, int, bool) {
	if self.slotOf(name) < 0 && self.lookupUpval(name) < 0 {
		envIdx := self.lookupUpval("_ENV")
		nameIdx := self.indexOfConstant(name)
		return envIdx, nameIdx, true
	} else {
		return -1, -1, false
	}
}
func (self *codeGen) indexOfConstant(k interface{}) int {
	return self.scope.indexOfConstant(k)
}

func (self *codeGen) genSubProto(fd *FuncDefExp) int {
	proto := newCodeGen(self.scope).genProto(fd)
	self.protos = append(self.protos, proto)
	return len(self.protos) - 1
}

func (self *codeGen) genProto(fd *FuncDefExp) *Prototype {
	if fd.Line == 0 { // main
		self.scope.setupEnv()
	}

	for _, param := range fd.ParList {
		self.addLocVar(param, 0)
	}

	self.cgBlock(fd.Block)

	endPc := self.pc() + 2
	self.exitScope(endPc)

	return self.toProto(fd)
}

func (self *codeGen) toProto(fd *FuncDefExp) *Prototype {
	proto := &Prototype{
		LineDefined:     uint32(fd.Line),
		LastLineDefined: uint32(fd.LastLine),
		NumParams:       byte(len(fd.ParList)),
		MaxStackSize:    byte(self.scope.getMaxStack()),
		Code:            self.insts.encode(),
		Constants:       self.scope.getConstants(),
		Upvalues:        self.scope.getUpvalues(),
		Protos:          self.protos,
		LineInfo:        self.insts.getLineNumTable(),
		LocVars:         self.scope.getLocVars(),
		UpvalueNames:    self.scope.getUpvalueNames(),
	}

	if fd.Line == 0 {
		proto.LastLineDefined = 0
	}
	if proto.MaxStackSize < 2 {
		proto.MaxStackSize = 2 // todo
	}
	if fd.IsVararg {
		proto.IsVararg = 1 // todo
	}

	proto.Code = append(proto.Code, 0x00800026) // todo
	proto.LineInfo = append(proto.LineInfo, uint32(fd.LastLine))

	return proto
}

func GenProto(fd *FuncDefExp) *Prototype {
	return newCodeGen(nil).genProto(fd)
}
