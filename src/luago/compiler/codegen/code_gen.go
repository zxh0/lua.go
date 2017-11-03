package codegen

import . "luago/binchunk"
import . "luago/compiler/ast"

type codeGen struct {
	scope  *scope
	insts  []uint32
	lines  []uint32
	protos []*Prototype
}

func newCodeGen(parentScope *scope) *codeGen {
	return &codeGen{
		scope:  newScope(parentScope),
		insts:  make([]uint32, 0, 8),
		lines:  make([]uint32, 0, 8),
		protos: []*Prototype{},
	}
}

// constants
func (self *codeGen) indexOfConstant(k interface{}) int {
	return self.scope.indexOfConstant(k)
}

// registers
func (self *codeGen) usedRegs() int {
	return self.scope.stackSize
}
func (self *codeGen) resetRegs(n int) {
	self.scope.stackSize = n
}
func (self *codeGen) allocReg() int {
	return self.scope.allocReg()
}
func (self *codeGen) freeReg() {
	self.scope.freeReg()
}
func (self *codeGen) allocRegs(n int) int {
	return self.scope.allocRegs(n)
}
func (self *codeGen) freeRegs(n int) {
	self.scope.freeRegs(n)
}

// lexical scope
func (self *codeGen) enterScope(breakable bool) {
	self.scope.incrLevel(breakable)
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
func (self *codeGen) indexOfLocVar(name string) int {
	return self.scope.indexOfLocVar(name)
}
func (self *codeGen) addBreakJmp(pc int) {
	self.scope.addBreakJmp(pc)
}

// upvalues
func (self *codeGen) indexOfUpval(name string) int {
	return self.scope.indexOfUpval(name)
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
		Code:            self.insts,
		Constants:       self.scope.getConstants(),
		Upvalues:        self.scope.getUpvalues(),
		Protos:          self.protos,
		LineInfo:        self.lines,
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
