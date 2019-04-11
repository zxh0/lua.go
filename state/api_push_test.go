package state

import "testing"
import . "github.com/zxh0/lua.go/api"

func TestPushGoClosure(t *testing.T) {
	ls := New()
	ls.PushInteger(1)
	ls.PushInteger(2)
	ls.PushInteger(3)
	ls.PushGoClosure(func(ls LuaState) int {
		if ls.CheckInteger(LUA_REGISTRYINDEX-1) != 1 {
			t.Fail()
		}
		if ls.CheckInteger(LUA_REGISTRYINDEX-2) != 2 {
			t.Fail()
		}
		if ls.CheckInteger(LUA_REGISTRYINDEX-3) != 3 {
			t.Fail()
		}
		return 0
	}, 3)
	ls.Call(0, 0)
}
