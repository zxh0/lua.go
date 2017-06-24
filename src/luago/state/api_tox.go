package state

import "strconv"
import . "luago/lua"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (self *luaState) ToBoolean(index int) bool {
	val := self.stack.get(index)
	return valToBoolean(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointeger
func (self *luaState) ToInteger(index int) int64 {
	i, _ := self.ToIntegerX(index)
	return i
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointegerx
func (self *luaState) ToIntegerX(index int) (int64, bool) {
	val := self.stack.get(index)
	return valToInteger(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumber
func (self *luaState) ToNumber(index int) float64 {
	n, _ := self.ToNumberX(index)
	return n
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumberx
func (self *luaState) ToNumberX(index int) (float64, bool) {
	val := self.stack.get(index)
	return valToNumber(val)
}

// [-0, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_tostring
// http://www.lua.org/manual/5.3/manual.html#lua_tolstring
func (self *luaState) ToString(index int) (string, bool) {
	val := self.stack.get(index)

	s := ""
	switch x := val.(type) {
	case int64:
		s = strconv.FormatInt(x, 10) // todo
	case float64:
		s = strconv.FormatFloat(x, 'f', -1, 64) // todo
	case string:
		return x, true
	default:
		return "", false
	}

	// val is a number
	self.CheckStack(1)
	self.PushString(s)
	self.Replace(index)
	return s, true
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tocfunction
func (self *luaState) ToGoFunction(index int) LuaGoFunction {
	val := self.stack.get(index)
	switch x := val.(type) {
	case LuaGoFunction:
		return x
	default:
		return nil
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_touserdata
func (self *luaState) ToUserData(index int) LuaUserData {
	val := self.stack.get(index)
	if val != nil {
		if ud, ok := val.(*userData); ok {
			return ud.data
		}
	}
	return nil
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_topointer
func (self *luaState) ToPointer(index int) interface{} {
	return self.stack.get(index)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tothread
func (self *luaState) ToThread(index int) LuaState {
	panic("todo!")
}
