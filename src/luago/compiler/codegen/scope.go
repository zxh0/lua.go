package codegen

import "luago/binchunk"

type upvalInfo struct {
	instack bool
	idx     int
	seq     int
}

type locVarInfo struct {
	prev    *locVarInfo
	name    string
	level   int
	slot    int
	startPc int
	endPc   int
}

type breakInfo struct {
	breakable bool
	breakJmps []int
}

type scope struct {
	parent    *scope
	level     int // blockLevel
	locVars   []*locVarInfo
	locNames  map[string]*locVarInfo
	upvalues  map[string]upvalInfo
	constants map[interface{}]int
	stackSize int // usedRegs
	stackMax  int // maxRegs
	breaks    []*breakInfo
}

func newScope(parent *scope) *scope {
	return &scope{
		parent:    parent,
		locVars:   make([]*locVarInfo, 0, 8),
		locNames:  map[string]*locVarInfo{},
		upvalues:  map[string]upvalInfo{},
		constants: map[interface{}]int{},
		breaks:    make([]*breakInfo, 1),
	}
}

/* local vars */

func (self *scope) incrLevel() {
	self.level++
	self.breaks = append(self.breaks, nil)
}

func (self *scope) decrLevel(endPc int) []int {
	self.level--
	for _, locVar := range self.locNames {
		if locVar.level > self.level { // out of scope
			locVar.endPc = endPc
			self.removeLocVar(locVar)
		}
	}

	breakInfo := self.breaks[len(self.breaks)-1]
	self.breaks = self.breaks[:len(self.breaks)-1]
	if breakInfo != nil {
		return breakInfo.breakJmps
	} else {
		return nil
	}
}

func (self *scope) removeLocVar(locVar *locVarInfo) {
	self.freeReg()
	if locVar.prev != nil {
		if locVar.prev.level == locVar.level {
			self.removeLocVar(locVar.prev)
		} else {
			self.locNames[locVar.name] = locVar.prev
		}
	} else {
		delete(self.locNames, locVar.name)
	}
	// for locVar.prev != nil {
	// 	if locVar.prev.level == locVar.level {
	// 		locVar.prev.endPc = locVar.endPc
	// 		locVar = locVar.prev
	// 		self.freeReg()
	// 	} else {
	// 		self.locNames[locVar.name] = locVar.prev
	// 		return
	// 	}
	// }
	// delete(self.locNames, locVar.name)
}

func (self *scope) addLocVar(name string, startPc int) int {
	newVar := &locVarInfo{
		prev:    self.locNames[name],
		name:    name,
		level:   self.level,
		slot:    self.allocReg(),
		startPc: startPc,
		endPc:   0,
	}

	self.locVars = append(self.locVars, newVar)
	self.locNames[name] = newVar

	return newVar.slot
}

func (self *scope) indexOfLocVar(name string) int {
	if locVar, found := self.locNames[name]; found {
		return locVar.slot
	} else {
		return -1
	}
}

/* stack */

func (self *scope) allocRegs(n int) int {
	if n > 0 {
		slot := self.allocReg()
		for i := 1; i < n; i++ {
			self.allocReg()
		}
		return slot
	} else {
		panic("n <= 0 !")
	}
}
func (self *scope) freeRegs(n int) {
	if n >= 0 {
		for i := 0; i < n; i++ {
			self.freeReg()
		}
	} else {
		panic("n < 0!")
	}
}

func (self *scope) allocReg() int {
	self.stackSize++
	if self.stackSize > self.stackMax {
		self.stackMax = self.stackSize
	}
	return self.stackSize - 1
}

func (self *scope) freeReg() {
	self.stackSize--
}

/* upvalues */

func (self *scope) setupEnv() {
	self.upvalues["_ENV"] = upvalInfo{
		instack: true,
		idx:     0,
		seq:     0,
	}
}

func (self *scope) indexOfUpval(name string) int {
	if uvInfo, ok := self.upvalues[name]; ok {
		return uvInfo.seq
	}
	if self.parent != nil {
		seq := len(self.upvalues)
		if locVar, found := self.parent.locNames[name]; found {
			self.upvalues[name] = upvalInfo{
				instack: true,
				idx:     locVar.slot,
				seq:     seq,
			}
			return seq
		}
		if idx := self.parent.indexOfUpval(name); idx >= 0 {
			self.upvalues[name] = upvalInfo{
				instack: false,
				idx:     idx,
				seq:     seq,
			}
			return seq
		}
		self.indexOfUpval("_ENV")
		return -1
	}
	return -1
}

/* constants */

func (self *scope) indexOfConstant(k interface{}) int {
	if idx, found := self.constants[k]; found {
		return idx
	}

	idx := len(self.constants)
	self.constants[k] = idx
	return idx
}

/* break support */

func (self *scope) markBreakable() {
	self.breaks[self.level] = &breakInfo{
		breakable: true,
		breakJmps: []int{},
	}
}

func (self *scope) addBreakJmp(pc int) {
	var breakInfo *breakInfo
	for i := self.level; i >= 0; i-- {
		breakInfo = self.breaks[i]
		if breakInfo != nil {
			break
		}
	}

	if breakInfo == nil {
		panic("<break> at line ? not inside a loop!")
	} else {
		breakInfo.breakJmps = append(breakInfo.breakJmps, pc)
	}
}

/* summarize */

func (self *scope) getMaxStack() int {
	return self.stackMax
}

func (self *scope) getLocVars() []binchunk.LocVar {
	locVars := make([]binchunk.LocVar, len(self.locVars))
	for i, locVar := range self.locVars {
		locVars[i] = binchunk.LocVar{
			VarName: locVar.name,
			StartPc: uint32(locVar.startPc),
			EndPc:   uint32(locVar.endPc),
		}
	}
	return locVars
}

func (self scope) getUpvalues() []binchunk.Upvalue {
	upvals := make([]binchunk.Upvalue, len(self.upvalues))
	for _, uv := range self.upvalues {
		if uv.instack {
			upvals[uv.seq] = binchunk.Upvalue{1, byte(uv.idx)}
		} else {
			upvals[uv.seq] = binchunk.Upvalue{0, byte(uv.idx)}
		}
	}
	return upvals
}

func (self scope) getUpvalueNames() []string {
	names := make([]string, len(self.upvalues))
	for name, uv := range self.upvalues {
		names[uv.seq] = name
	}
	return names
}

func (self scope) getConstants() []interface{} {
	consts := make([]interface{}, len(self.constants))
	for k, idx := range self.constants {
		consts[idx] = k
	}
	return consts
}
