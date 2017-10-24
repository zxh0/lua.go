package codegen

type tmpAllocator struct {
	scope   *scope
	tmpVar  int
	counter int
}

func (self *tmpAllocator) allocReg() int {
	if self.tmpVar >= 0 {
		tmp := self.tmpVar
		self.tmpVar = -1
		return tmp
	} else {
		self.counter += 1
		return self.scope.allocReg()
	}
}
func (self *tmpAllocator) freeReg() {
	self.counter -= 1
	self.scope.freeReg()
}

func (self *tmpAllocator) allocRegs(n int) int {
	self.counter += n
	return self.scope.allocRegs(n)
}
func (self *tmpAllocator) freeRegs(n int) {
	self.counter -= n
	self.scope.freeRegs(n)
}

func (self *tmpAllocator) freeAll() {
	if self.counter > 0 {
		self.scope.freeRegs(self.counter)
		self.counter = 0
	}
}
