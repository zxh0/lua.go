package main

import "os"
import . "luago/api"
import "luago/state"

func main() {
	if len(os.Args) > 1 {
		ls := state.New()
		ls.OpenLibs()
		if ls.LoadFile(os.Args[1]) == LUA_OK {
			ls.PCall(0, LUA_MULTRET, -1)
		} else {
			panic(ls.CheckString(-1)) // todo
		}
	}
}
