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

// coroutine.create (f)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.create
func coCreate(ls LuaState) int {
	t := ls.NewThread()
	ls.PushThread(t)
	return 1
}

// coroutine.resume (co [, val1, ···])
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.resume
func coResume(ls LuaState) int {
	panic("todo: coResume!")
}

// coroutine.yield (···)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.yield
func coYield(ls LuaState) int {
	panic("todo: coYield!")
}

// coroutine.status (co)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.status
func coStatus(ls LuaState) int {
	panic("todo: coStatus!")
}

// coroutine.isyieldable ()
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.isyieldable
func coYieldable(ls LuaState) int {
	panic("todo: coYieldable!")
}

// coroutine.running ()
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.running
func coRunning(ls LuaState) int {
	panic("todo: coRunning!")
}

// coroutine.wrap (f)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.wrap
func coWrap(ls LuaState) int {
	panic("todo: coWrap!")
}
