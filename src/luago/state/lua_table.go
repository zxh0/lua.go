package state

import "luago/luanum"

type pair struct {
	key luaValue
	val luaValue
}

type luaTable struct {
	metatable *luaTable
	_map      map[luaValue]luaValue
	arr       []luaValue
	_len      int
	pairs     []pair // used by next()
	nextPairs []pair // used by next()
}

func newLuaTable(nArr, nRec int) *luaTable {
	return &luaTable{
		_map: make(map[luaValue]luaValue),
		arr:  make([]luaValue, 0, nArr),
	}
}

func (self *luaTable) hasMetafield(fieldName string) bool {
	return self.metatable != nil &&
		self.metatable.get(fieldName) != nil
}

func (self *luaTable) len() int {
	if self._len < 0 {
		// calc & cache len of sequence
		self._len = 0
		for _, val := range self.arr {
			if val != nil {
				self._len++
			} else {
				break
			}
		}
	}

	return self._len
}

func (self *luaTable) get(key luaValue) luaValue {
	key = _floatToIntger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(self.arr))
		if idx <= arrLen {
			return self.arr[idx-1]
		}
	}
	if val, found := self._map[key]; found {
		return val
	} else {
		return nil
	}
}

func (self *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}

	self.pairs = nil // invalidate pairs
	self._len = -1

	key = _floatToIntger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(self.arr))
		if idx <= arrLen {
			self.arr[idx-1] = val
			return
		}
		if idx == arrLen+1 {
			delete(self._map, key)
			self.arr = append(self.arr, val)
			self.expandArray()
			return
		}
	}
	if val != nil {
		self._map[key] = val
	} else {
		delete(self._map, key)
	}
}

func (self *luaTable) expandArray() {
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

func _floatToIntger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := luanum.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}
