package state

type pair struct {
	key luaValue
	val luaValue
}

// todo: move to types?
type luaTable struct {
	metaTable *luaTable
	_map      map[luaValue]luaValue
	list      []luaValue
	pairs     []pair // used by next()
	nextPairs []pair // used by next()
}

func newLuaTable(nArr, nRec int) *luaTable {
	return &luaTable{
		list: make([]luaValue, 0, nArr),
		_map: make(map[luaValue]luaValue),
	}
}

func (self *luaTable) hasMetaField(fieldName string) bool {
	return self.metaTable != nil &&
		self.metaTable.get(fieldName) != nil
}

func (self *luaTable) len() int {
	listLen := len(self.list)
	for listLen > 0 { // remove tail nils
		if self.list[listLen-1] == nil {
			self.list = self.list[0 : listLen-1]
			listLen -= 1
		} else {
			break
		}
	}
	return listLen
}

func (self *luaTable) get(key luaValue) luaValue {
	// todo: try to cast float key to integer
	if idx, ok := key.(int64); ok && idx >= 1 {
		listLen := int64(len(self.list))
		if idx <= listLen {
			return self.list[idx-1]
		}
	}
	if val, found := self._map[key]; found {
		return val
	} else {
		return nil
	}
}

func (self *luaTable) put(key, val luaValue) {
	// todo: try to cast float key to integer
	self.pairs = nil // invalidate pairs
	if idx, ok := key.(int64); ok && idx >= 1 {
		listLen := int64(len(self.list))
		if idx <= listLen {
			self.list[idx-1] = val
			return
		}
		if idx == listLen+1 {
			self.list = append(self.list, val)
			if val != nil {
				self.growList()
			}
			return
		}
	}
	if val != nil {
		self._map[key] = val
	} else {
		delete(self._map, key)
	}
}

// todo: rename?
func (self *luaTable) growList() {
	for {
		key := int64(len(self.list) + 1) // todo
		if val, found := self._map[key]; found {
			delete(self._map, key)
			self.list = append(self.list, val)
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
			return int64(1), self.list[0]
		}
		if idx, ok := key.(int64); ok && idx >= 1 {
			if idx < int64(len(self.list)) {
				return int64(idx + 1), self.list[idx]
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
