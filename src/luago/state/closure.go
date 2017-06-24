package state

import . "luago/lua"
import "luago/binchunk"

type goClosure struct {
	goFunc LuaGoFunction
	upvals []luaValue
}

type luaClosure struct {
	proto  *binchunk.FuncProto
	upvals []*luaValue
}

func newLuaClosure(proto *binchunk.FuncProto) *luaClosure {
	upvals := make([]*luaValue, len(proto.Upvalues))
	return &luaClosure{
		proto:  proto,
		upvals: upvals,
	}
}
