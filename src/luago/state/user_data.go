package state

// todo
type userData struct {
	metatable *luaTable
	userValue luaValue
	data      interface{} // anything
}
