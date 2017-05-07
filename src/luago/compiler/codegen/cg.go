package codegen

import . "luago/binchunk"
import . "luago/compiler/ast"

type cg struct {
	scope  *scope
	insts  insts
	protos []*FuncProto
}

func newCG(parentScope *scope) *cg {
	return &cg{
		scope:  newScope(parentScope),
		insts:  make([]instruction, 0, 8),
		protos: []*FuncProto{},
	}
}

// todo: use pc0()
func (self *cg) pc() int {
	return len(self.insts) + 1 // ???
}

func (self *cg) pc0() int {
	return len(self.insts) - 1
}

func (self *cg) inst(line, opcode, a, b, c int) int {
	i := instruction{line, opcode, a, b, c}
	self.insts = append(self.insts, i)
	return len(self.insts) - 1
}

func (self *cg) fixA(pc, a int) {
	i := self.insts[pc]
	i.a = a
	self.insts[pc] = i
}

func (self *cg) fixSbx(pc, sbx int) {
	i := self.insts[pc]
	i.b = sbx
	self.insts[pc] = i
}

func (self *cg) enterScope() {
	self.scope.incrLevel()
}
func (self *cg) exitScope(endPc int) {
	self.scope.decrLevel(endPc)
}
func (self *cg) addLocVar(name string, startPc int) int {
	return self.scope.addLocVar(name, startPc)
}
func (self *cg) slotOf(name string) int {
	return self.scope.slotOf(name)
}
func (self *cg) lookupUpval(name string) int {
	return self.scope.lookupUpval(name)
}
func (self *cg) allocTmps(n int) int {
	return self.scope.allocTmps(n)
}
func (self *cg) allocTmp() int {
	return self.scope.allocTmp()
}
func (self *cg) freeTmp() {
	self.scope.freeTmp()
}
func (self *cg) freeTmps(n int) {
	self.scope.freeTmps(n)
}
func (self *cg) indexOf(k interface{}) int {
	return self.scope.indexOf(k)
}

func (self *cg) genSubProto(fd *FuncDefExp) int {
	proto := newCG(self.scope).genProto(fd)
	self.protos = append(self.protos, proto)
	return len(self.protos) - 1
}

func (self *cg) genProto(fd *FuncDefExp) *FuncProto {
	if fd.Line == 0 { // main
		self.scope.setupEnv()
	}

	for _, param := range fd.ParList {
		self.addLocVar(param, 0)
	}

	self.block(fd.Block)

	endPc := self.pc()
	self.exitScope(endPc)

	return self.toProto(fd)
}

func (self *cg) toProto(fd *FuncDefExp) *FuncProto {
	proto := &FuncProto{
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

func GenProto(fd *FuncDefExp) *FuncProto {
	return newCG(nil).genProto(fd)
}
