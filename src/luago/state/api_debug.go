package state

import . "luago/api"

func (self *luaState) GetHook() LuaHook {
	panic("todo: GetHook!")
}

func (self *luaState) SetHook(f LuaHook, mask, count int) {
	panic("todo: SetHook!")
}

func (self *luaState) GetHookCount() int {
	panic("todo: GetHookCount!")
}

func (self *luaState) GetHookMask() int {
	panic("todo: GetHookMask!")
}

func (self *luaState) GetInfo(what string, ar *LuaDebug) {
	panic("todo: GetInfo!")
}

func (self *luaState) GetLocal(ar *LuaDebug, n int) string {
	panic("todo: GetLocal!")
}

func (self *luaState) SetLocal(ar *LuaDebug, n int) string {
	panic("todo: SetLocal!")
}

func (self *luaState) GetStack(level int, ar *LuaDebug) int {
	// todo
	if self.callDepth > 1 {
		return 1
	}
	return 0
}

func (self *luaState) GetUpvalue(funcIdx, n int) string {
	panic("todo: GetUpvalue!")
}

func (self *luaState) SetUpvalue(funcIdx, n int) string {
	panic("todo: SetUpvalue!")
}

func (self *luaState) UpvalueId(funcIdx, n int) interface{} {
	panic("todo: UpvalueId!")
}

func (self *luaState) UpvalueJoin(funcIdx1, n1, funcIdx2, n2 int) {
	panic("todo: UpvalueJoin!")
}
