package state

import . "luago/api"
import "luago/binchunk"

type goClosure struct {
	goFunc GoFunction
	upvals []luaValue
}

type luaClosure struct {
	proto  *binchunk.Prototype
	upvals []*luaValue
}

func newLuaClosure(proto *binchunk.Prototype) *luaClosure {
	upvals := make([]*luaValue, len(proto.Upvalues))
	return &luaClosure{
		proto:  proto,
		upvals: upvals,
	}
}
