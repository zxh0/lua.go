package state

import . "github.com/zxh0/lua.go/api"
import "github.com/zxh0/lua.go/binchunk"

type upvalue struct {
	val *luaValue
}

type closure struct {
	proto  *binchunk.Prototype // lua closure
	goFunc GoFunction          // go closure
	upvals []*upvalue
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	c := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(f GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func (self *closure) getUpvalueName(n int) string {
	if self.proto != nil {
		if len(self.proto.UpvalueNames) > n {
			return self.proto.UpvalueNames[n]
		}
	}
	return ""
}

func (self *closure) getUpvalue(n int) luaValue {
	if self.upvals[n] == nil || self.upvals[n].val == nil {
		return nil
	}
	return *(self.upvals[n].val)
}

func (self *closure) setUpvalue(n int, uv luaValue) {
	if self.upvals[n] == nil {
		self.upvals[n] = &upvalue{}
	}
	if self.upvals[n].val == nil {
		self.upvals[n].val = &uv
	} else {
		*(self.upvals[n].val) = uv
	}
}
