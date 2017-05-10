package codegen

type tmpAllocator struct {
	scope   *scope
	tmpVar  int
	counter int
}

// func (self tmpAllocator) allocTmps(n int) int {
// 	self.counter += n
// 	return self.scope.allocTmps(n)
// }
func (self *tmpAllocator) allocTmp() int {
	if self.tmpVar >= 0 {
		tmp := self.tmpVar
		self.tmpVar = -1
		return tmp
	} else {
		self.counter += 1
		return self.scope.allocTmp()
	}
}
func (self *tmpAllocator) freeTmp() {
	self.counter -= 1
	self.scope.freeTmp()
}
func (self *tmpAllocator) freeTmps(n int) {
	self.counter -= n
	self.scope.freeTmps(n)
}

func (self *tmpAllocator) freeAll() {
	if self.counter > 0 {
		self.scope.freeTmps(self.counter)
		self.counter = 0
	}
}
