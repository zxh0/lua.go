package state

import "math"
import "luago/number"

type pair struct {
	key luaValue
	val luaValue
}

type luaTable struct {
	metatable *luaTable
	_map      map[luaValue]luaValue
	arr       []luaValue
	pairs     []pair // used by next()
	nextPairs []pair // used by next()
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

func (self *luaTable) hasMetafield(fieldName string) bool {
	return self.metatable != nil &&
		self.metatable.get(fieldName) != nil
}

func (self *luaTable) len() int {
	return len(self.arr)
}

func (self *luaTable) get(key luaValue) luaValue {
	key = _floatToIntger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(self.arr))
		if idx <= arrLen {
			return self.arr[idx-1]
		}
	}
	return self._map[key]
}

func _floatToIntger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (self *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	self.pairs = nil // invalidate pairs

	key = _floatToIntger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(self.arr))
		if idx <= arrLen {
			self.arr[idx-1] = val
			if idx == arrLen && val == nil {
				self._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(self._map, key)
			if val != nil {
				self.arr = append(self.arr, val)
				self._expandArray()
			}
			return
		}
	}
	if val != nil {
		if self._map == nil {
			self._map = make(map[luaValue]luaValue, 8)
		}
		self._map[key] = val
	} else {
		delete(self._map, key)
	}
}

func (self *luaTable) _shrinkArray() {
	for i := len(self.arr) - 1; i >= 0; i-- {
		if self.arr[i] == nil {
			self.arr = self.arr[0:i]
		}
	}
}

func (self *luaTable) _expandArray() {
	for idx := int64(len(self.arr)) + 1; true; idx++ {
		if val, found := self._map[idx]; found {
			delete(self._map, idx)
			self.arr = append(self.arr, val)
		} else {
			break
		}
	}
}

func (self *luaTable) next(key luaValue) (nextKey, nextVal luaValue) {
	if key == nil {
		if self.pairs == nil {
			self.initPairs()
		}
		self.nextPairs = self.pairs
	}

	if self.len() > 0 {
		if key == nil {
			return int64(1), self.arr[0]
		}
		if idx, ok := key.(int64); ok && idx >= 1 {
			if idx < int64(len(self.arr)) {
				return int64(idx + 1), self.arr[idx]
			}
		}
	}

	if len(self.nextPairs) > 0 {
		pair := self.nextPairs[0]
		self.nextPairs = self.nextPairs[1:]
		return pair.key, pair.val
	} else {
		return nil, nil
	}
}

func (self *luaTable) initPairs() {
	n := len(self._map)
	if n > 0 {
		self.pairs = make([]pair, 0, n)
		for key, val := range self._map {
			self.pairs = append(self.pairs, pair{key, val})
		}
	}
}
