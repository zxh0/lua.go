package state

import "fmt"
import "reflect"
import "runtime"
import "strings"
import . "luago/lua"
import . "luago/number"

var _mtOfNil *luaTable = nil //?
var _mtOfBool *luaTable = nil
var _mtOfNumber *luaTable = nil
var _mtOfString *luaTable = nil
var _mtOfFunc *luaTable = nil
var _mtOfThread *luaTable = nil

type luaValue interface{}

func getMetaTable(val luaValue) *luaTable {
	switch x := val.(type) {
	case nil:
		return _mtOfNil
	case bool:
		return _mtOfBool
	case int64, float64:
		return _mtOfNumber
	case string:
		return _mtOfString
	case *luaClosure, *goClosure, LuaGoFunction:
		return _mtOfFunc
	case *luaTable:
		return x.metaTable
	case *userData:
		return x.metaTable
	default:
		panic("todo!")
	}
}

func setMetaTable(val luaValue, mt *luaTable) {
	switch x := val.(type) {
	case nil:
		_mtOfNil = mt
	case bool:
		_mtOfBool = mt
	case int64, float64:
		_mtOfNumber = mt
	case string:
		_mtOfString = mt
	case *luaClosure, *goClosure, LuaGoFunction:
		_mtOfFunc = mt
	case *luaTable:
		x.metaTable = mt
	case *userData:
		x.metaTable = mt
	default:
		panic("todo!")
	}
}

func getMetaField(val luaValue, fieldName string) luaValue {
	if mt := getMetaTable(val); mt != nil {
		return mt.get(fieldName)
	} else {
		return nil
	}
}

func typeOf(val luaValue) LuaType {
	return fullTypeOf(val) & 0x0F
}

func fullTypeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMINT
	case float64:
		return LUA_TNUMFLT
	case string:
		return LUA_TSTRING
	case *luaTable:
		return LUA_TTABLE
	case *userData:
		return LUA_TUSERDATA
	case *luaClosure:
		return LUA_TLCL
	case *goClosure:
		return LUA_TGCL
	case LuaGoFunction:
		return LUA_TLGF
	default:
		panic("todo")
	}
}

func valToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

// todo
func valToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return CastToInteger(x)
	case string:
		x = strings.TrimSpace(x)
		if i, ok := ParseInteger(x); ok {
			return i, true
		}
		if f, ok := ParseFloat(x); ok {
			return CastToInteger(f)
		}
	}
	return 0, false
}

func valToNumber(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case int64:
		return float64(x), true
	case float64:
		return x, true
	case string:
		x = strings.TrimSpace(x)
		return ParseFloat(x)
	}
	return 0, false
}

// debug
func valToString(val luaValue) string {
	switch x := val.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%f", val)
	case string:
		return fmt.Sprintf("%q", val)
	case *luaTable:
		return fmt.Sprintf("{@%p}", val)
	case *luaClosure:
		return luaFuncToString(x)
	case *goClosure:
		return goFuncToString(x.goFunc) + "!"
	case LuaGoFunction:
		return goFuncToString(val)
	default:
		fmt.Printf("%T\n", val)
		panic("todo")
	}
}

func luaFuncToString(luaf *luaClosure) string {
	return fmt.Sprintf("<%s:%d,%d>",
		luaf.proto.Source, // todo
		luaf.proto.LineDefined,
		luaf.proto.LastLineDefined)
}

func goFuncToString(gof luaValue) string {
	pc := reflect.ValueOf(gof).Pointer()
	if f := runtime.FuncForPC(pc); f != nil {
		name := f.Name()[strings.LastIndex(f.Name(), ".")+1:]
		return fmt.Sprintf("%s()", name)
	}
	return fmt.Sprintf("(@%p)", gof)
}
