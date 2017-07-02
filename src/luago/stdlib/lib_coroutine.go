package stdlib

import . "luago/api"

var coFuncs = map[string]LuaGoFunction{
	"create":      coCreate,
	"resume":      coResume,
	"running":     coRunning,
	"status":      coStatus,
	"wrap":        coWrap,
	"yield":       coYield,
	"isyieldable": coYieldable,
}

func OpenCoroutineLib(ls LuaState) int {
	ls.NewLib(coFuncs)
	return 1
}

func coCreate(ls LuaState) int {
	panic("todo: coCreate!")
}

func coResume(ls LuaState) int {
	panic("todo: coResume!")
}

func coRunning(ls LuaState) int {
	panic("todo: coRunning!")
}

func coStatus(ls LuaState) int {
	panic("todo: coStatus!")
}

func coWrap(ls LuaState) int {
	panic("todo: coWrap!")
}

func coYield(ls LuaState) int {
	panic("todo: coYield!")
}

func coYieldable(ls LuaState) int {
	panic("todo: coYieldable!")
}
