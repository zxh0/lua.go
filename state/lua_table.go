package state

import (
	"math"

	"github.com/zxh0/lua.go/number"
)

type luaTable struct {
	metatable *luaTable
	arr       []luaValue
	_map      map[luaValue]luaValue
	keys      map[luaValue]luaValue // used by next()
	lastKey   luaValue              // used by next()
	changed   bool                  // used by next()
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (table *luaTable) hasMetafield(fieldName string) bool {
	return table.metatable != nil &&
		table.metatable.get(fieldName) != nil
}

func (table *luaTable) len() int {
	return len(table.arr)
}

func (table *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(table.arr)) {
			return table.arr[idx-1]
		}
	}
	return table._map[key]
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (table *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	table.changed = true
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(table.arr))
		if idx <= arrLen {
			table.arr[idx-1] = val
			if idx == arrLen && val == nil {
				table._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(table._map, key)
			if val != nil {
				table.arr = append(table.arr, val)
				table._expandArray()
			}
			return
		}
	}
	if val != nil {
		if table._map == nil {
			table._map = make(map[luaValue]luaValue, 8)
		}
		table._map[key] = val
	} else {
		delete(table._map, key)
	}
}

func (table *luaTable) _shrinkArray() {
	for i := len(table.arr) - 1; i >= 0; i-- {
		if table.arr[i] == nil {
			table.arr = table.arr[0:i]
		}
	}
}

func (table *luaTable) _expandArray() {
	for idx := int64(len(table.arr)) + 1; true; idx++ {
		if val, found := table._map[idx]; found {
			delete(table._map, idx)
			table.arr = append(table.arr, val)
		} else {
			break
		}
	}
}

func (table *luaTable) nextKey(key luaValue) luaValue {
	if table.keys == nil || (key == nil && table.changed) {
		table.initKeys()
		table.changed = false
	}

	nextKey := table.keys[key]
	if nextKey == nil && key != nil && key != table.lastKey {
		panic("invalid key to 'next'")
	}

	return nextKey
}

func (table *luaTable) initKeys() {
	table.keys = make(map[luaValue]luaValue)
	var key luaValue = nil
	for i, v := range table.arr {
		if v != nil {
			table.keys[key] = int64(i + 1)
			key = int64(i + 1)
		}
	}
	for k, v := range table._map {
		if v != nil {
			table.keys[key] = k
			key = k
		}
	}
	table.lastKey = key
}
