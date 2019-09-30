package state

import (
	. "github.com/zxh0/lua.go/api"
	"github.com/zxh0/lua.go/binchunk"
)

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

func (c *closure) getUpvalueName(n int) string {
	if c.proto != nil {
		if len(c.proto.UpvalueNames) > n {
			return c.proto.UpvalueNames[n]
		}
	}
	return ""
}

func (c *closure) getUpvalue(n int) luaValue {
	if c.upvals[n] == nil || c.upvals[n].val == nil {
		return nil
	}
	return *(c.upvals[n].val)
}

func (c *closure) setUpvalue(n int, uv luaValue) {
	if c.upvals[n] == nil {
		c.upvals[n] = &upvalue{}
	}
	if c.upvals[n].val == nil {
		c.upvals[n].val = &uv
	} else {
		*(c.upvals[n].val) = uv
	}
}
