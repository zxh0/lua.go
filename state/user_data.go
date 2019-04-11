package state

// todo
type userdata struct {
	metatable *luaTable
	userValue luaValue
	data      interface{} // anything
}
