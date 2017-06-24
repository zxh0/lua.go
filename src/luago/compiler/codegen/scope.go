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

type scope struct {
	parent    *scope
	level     int
	nLocals   int
	locVars   []*locVarInfo
	locNames  map[string]*locVarInfo
	upvalues  map[string]upvalInfo
	constants map[interface{}]int
	stackSize int
	stackMax  int
}

func newScope(parent *scope) *scope {
	return &scope{
		parent:    parent,
		locVars:   make([]*locVarInfo, 0, 8),
		locNames:  map[string]*locVarInfo{},
		upvalues:  map[string]upvalInfo{},
		constants: map[interface{}]int{},
	}
}

/* local vars */

func (self *scope) incrLevel() {
	self.level++
}

func (self *scope) decrLevel(endPc int) {
	self.stackSize = 0
	self.level--
	for _, locVar := range self.locNames {
		if locVar.level > self.level { // out of scope
			self.nLocals--
			locVar.endPc = endPc
			self.removeLocVar(locVar)
		}
	}
}

func (self *scope) removeLocVar(locVar *locVarInfo) {
	for locVar.prev != nil {
		if locVar.prev.level == locVar.level {
			locVar.prev.endPc = locVar.endPc
			locVar = locVar.prev
		} else {
			self.locNames[locVar.name] = locVar.prev
			return
		}
	}
	delete(self.locNames, locVar.name)
}

func (self *scope) addLocVar(name string, startPc int) int {
	newVar := &locVarInfo{
		prev:    self.locNames[name],
		name:    name,
		level:   self.level,
		slot:    self.nLocals,
		startPc: startPc,
		endPc:   0,
	}
	self.nLocals++

	self.locVars = append(self.locVars, newVar)
	self.locNames[name] = newVar

	return newVar.slot
}

func (self *scope) slotOf(name string) int {
	if locVar, found := self.locNames[name]; found {
		return locVar.slot
	} else {
		return -1
	}
}

func (self *scope) isLocVarSlot(slot int) bool {
	return slot >= 0 && slot < self.nLocals
}

func (self *scope) isTmpVar(slot int) bool {
	return slot >= self.nLocals && slot < 0x100
}

/* stack */

func (self *scope) allocTmps(n int) int {
	if n > 0 {
		slot := self.allocTmp()
		for i := 1; i < n; i++ {
			self.allocTmp()
		}
		return slot
	} else {
		panic("n <= 0 !")
	}
}

func (self *scope) allocTmp() int {
	if self.stackSize < self.nLocals {
		self.stackSize = self.nLocals
	}
	self.stackSize++
	if self.stackSize > self.stackMax {
		self.stackMax = self.stackSize
	}
	return self.stackSize - 1
}

func (self *scope) freeTmp() {
	if self.stackSize > self.nLocals {
		self.stackSize--
	}
}

func (self *scope) freeTmps(n int) {
	if n >= 0 {
		self.stackSize -= n
	} else {
		panic("n < 0!")
	}
}

/* upvalues */

func (self *scope) setupEnv() {
	self.upvalues["_ENV"] = upvalInfo{
		instack: true,
		idx:     0,
		seq:     0,
	}
}

func (self *scope) lookupUpval(name string) int {
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
		if idx := self.parent.lookupUpval(name); idx >= 0 {
			self.upvalues[name] = upvalInfo{
				instack: false,
				idx:     idx,
				seq:     seq,
			}
			return seq
		}
		self.lookupUpval("_ENV")
		return -1
	}
	return -1
}

/* constants */

func (self *scope) indexOfConstant(k interface{}) int {
	if idx, found := self.constants[k]; found {
		// todo: idx > 0xFF ?
		return idx + 0x100
	}

	idx := len(self.constants)
	self.constants[k] = idx
	return idx + 0x100
}

/* summarize */

func (self *scope) getMaxStack() int {
	maxLocals := 0
	for _, locVar := range self.locVars {
		if locVar.slot+1 > maxLocals {
			maxLocals = locVar.slot + 1
		}
	}
	if self.stackMax > maxLocals {
		return self.stackMax
	} else {
		return maxLocals
	}
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
