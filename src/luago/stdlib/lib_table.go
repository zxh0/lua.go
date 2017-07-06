package stdlib

import "sort"
import "strings"
import . "luago/api"

var tabFuncs = map[string]GoFunction{
	"insert": tabInsert,
	"remove": tabRemove,
	"move":   tabMove,
	"sort":   tabSort,
	"concat": tabConcat,
	"pack":   tabPack,
	"unpack": tabUnpack,
}

func OpenTableLib(ls LuaState) int {
	ls.NewLib(tabFuncs)
	return 1
}

// table.insert (list, [pos,] value)
// http://www.lua.org/manual/5.3/manual.html#pdf-table.insert
func tabInsert(ls LuaState) int {
	tabLen := ls.LenL(1)

	switch ls.GetTop() {
	case 2:
		ls.SetI(1, tabLen+1) // append
	case 3:
		if pos, ok := ls.ToIntegerX(2); ok {
			if pos == tabLen+1 {
				ls.SetI(1, pos) // append
			} else if pos >= 1 && pos <= tabLen { // insert
				for i := tabLen; i >= pos; i-- {
					ls.GetI(1, i)
					ls.SetI(1, i+1)
				}
				ls.SetI(1, pos)
			} else {
				panic("todo!")
			}
		} else {
			panic("todo!")
		}
	default:
		panic("todo!")
	}

	return 0
}

// table.remove (list [, pos])
// http://www.lua.org/manual/5.3/manual.html#pdf-table.remove
func tabRemove(ls LuaState) int {
	tabLen := ls.LenL(1)
	pos := ls.OptInteger(2, tabLen)

	if pos < 0 || pos == 0 && tabLen > 0 || pos > tabLen+1 {
		panic("todo!")
	}

	if pos == 0 && tabLen == 0 || pos == tabLen+1 {
		ls.PushNil()
		return 1
	}

	ls.GetI(1, pos)
	for i := pos; i < tabLen; i++ {
		ls.GetI(1, i+1)
		ls.SetI(1, i)
	}
	ls.PushNil()
	ls.SetI(1, tabLen)
	return 1
}

// table.move (a1, f, e, t [,a2])
// http://www.lua.org/manual/5.3/manual.html#pdf-table.move
func tabMove(ls LuaState) int {
	panic("todo: tabMove!")
}

// table.sort (list [, comp])
// http://www.lua.org/manual/5.3/manual.html#pdf-table.sort
func tabSort(ls LuaState) int {
	sort.Sort(wrapper{ls})
	return 0
}

// table.concat (list [, sep [, i [, j]]])
// http://www.lua.org/manual/5.3/manual.html#pdf-table.concat
func tabConcat(ls LuaState) int {
	tabLen := ls.LenL(1)

	sep := ls.OptString(2, "")
	i := ls.OptInteger(3, 1)
	j := ls.OptInteger(4, tabLen)

	if i < 1 || j > tabLen {
		panic("todo!")
	}

	if i > j {
		ls.PushString("")
		return 1
	}

	buf := make([]string, j-i+1)
	for k := i; k <= j; k++ {
		ls.GetI(1, int64(k))
		buf[k-1], _ = ls.ToString(-1)
		ls.Pop(1)
	}
	ls.PushString(strings.Join(buf, sep))

	return 1
}

// table.pack (···)
// http://www.lua.org/manual/5.3/manual.html#pdf-table.pack
func tabPack(ls LuaState) int {
	n := int64(ls.GetTop())   /* number of elements to pack */
	ls.CreateTable(int(n), 1) /* create result table */
	ls.Insert(1)              /* put it at index 1 */
	for i := n; i >= 1; i-- { /* assign elements */
		ls.SetI(1, i)
	}
	ls.PushInteger(n)
	ls.SetField(1, "n") /* t.n = number of elements */
	return 1            /* return table */
}

// table.unpack (list [, i [, j]])
// http://www.lua.org/manual/5.3/manual.html#pdf-table.unpack
func tabUnpack(ls LuaState) int {
	i := ls.OptInteger(2, 1)
	e := ls.OptInteger(3, ls.LenL(1))
	if i > e { /* empty range */
		return 0
	}
	// n = (lua_Unsigned)e - i;  /* number of elements minus 1 (avoid overflows) */
	// if (n >= (unsigned int)INT_MAX  || !lua_checkstack(L, (int)(++n)))
	//   return luaL_error(L, "too many results to unpack");
	n := int(e - i + 1)
	ls.CheckStack(n)
	for ; i < e; i++ { /* push arg[i..e - 1] (to avoid overflows) */
		ls.GetI(1, i)
	}
	ls.GetI(1, e) /* push last element */
	return n
}

/* table.sort */

type wrapper struct {
	ls LuaState
}

func (self wrapper) Len() int {
	return int(self.ls.LenL(1))
}

func (self wrapper) Less(i, j int) bool {
	ls := self.ls
	if ls.GetTop() == 2 { // cmp is given
		ls.PushValue(2)
		ls.GetI(1, int64(i+1))
		ls.GetI(1, int64(j+1))
		ls.Call(2, 1)
		b := ls.ToBoolean(-1)
		ls.Pop(1)
		return b
	} else { // cmp is missing
		ls.GetI(1, int64(i+1))
		ls.GetI(1, int64(j+1))
		b := ls.Compare(-2, -1, LUA_OPLT)
		ls.Pop(2)
		return b
	}
}

func (self wrapper) Swap(i, j int) {
	ls := self.ls
	ls.GetI(1, int64(i+1))
	ls.GetI(1, int64(j+1))
	ls.SetI(1, int64(i+1))
	ls.SetI(1, int64(j+1))
}
