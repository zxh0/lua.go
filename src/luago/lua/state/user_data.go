package state

// todo
type userData struct {
	metaTable *luaTable
	userValue luaValue
	data      interface{} // anything
}
