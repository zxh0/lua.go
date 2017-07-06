package stdlib

import "os"
import "time"
import . "luago/api"

var sysLib = map[string]GoFunction{
	"clock":     osClock,
	"difftime":  osDiffTime,
	"time":      osTime,
	"date":      osDate,
	"remove":    osRemove,
	"rename":    osRename,
	"tmpname":   osTmpName,
	"getenv":    osGetEnv,
	"execute":   osExecute,
	"exit":      osExit,
	"setlocale": osSetLocale,
}

func OpenOSLib(ls LuaState) int {
	ls.NewLib(sysLib)
	return 1
}

// os.clock ()
// http://www.lua.org/manual/5.3/manual.html#pdf-os.clock
func osClock(ls LuaState) int {
	panic("todo: osClock!")
}

// os.difftime (t2, t1)
// http://www.lua.org/manual/5.3/manual.html#pdf-os.difftime
func osDiffTime(ls LuaState) int {
	t2 := ls.ToInteger(1)
	t1 := ls.ToInteger(2)
	ls.PushInteger(t2 - t1)
	return 1
}

// os.time ([table])
// http://www.lua.org/manual/5.3/manual.html#pdf-os.time
func osTime(ls LuaState) int {
	if ls.GetTop() == 0 {
		t := time.Now().Unix()
		ls.PushInteger(t)
	} else {
		year := _getTimeField(ls, "year", -1)
		month := _getTimeField(ls, "month", -1)
		day := _getTimeField(ls, "day", -1)
		hour := _getTimeField(ls, "hour", 12)
		min := _getTimeField(ls, "min", 0)
		sec := _getTimeField(ls, "sec", 0)
		// todo: isdst
		t := time.Date(year, time.Month(month), day,
			hour, min, sec, 0, time.Local).Unix()
		ls.PushInteger(t)
	}
	return 1
}

func _getTimeField(ls LuaState, field string, defaultVal int) int {
	ls.GetField(1, field)
	if ls.IsNil(-1) {
		if defaultVal >= 0 {
			return defaultVal
		} else {
			panic("field '" + field + "' missing in date table")
		}
	}
	if i, ok := ls.ToIntegerX(-1); ok {
		return int(i)
	} else {
		panic("field '" + field + "' is not an integer")
	}
}

// os.date ([format [, time]])
// http://www.lua.org/manual/5.3/manual.html#pdf-os.date
func osDate(ls LuaState) int {
	panic("todo: osDate!")
}

// os.remove (filename)
// http://www.lua.org/manual/5.3/manual.html#pdf-os.remove
func osRemove(ls LuaState) int {
	name := ls.CheckString(1)
	if err := os.Remove(name); err != nil {
		ls.PushNil()
		ls.PushString(err.Error())
		return 2
	} else {
		ls.PushBoolean(true)
		return 1
	}
}

// os.rename (oldname, newname)
// http://www.lua.org/manual/5.3/manual.html#pdf-os.rename
func osRename(ls LuaState) int {
	oldName := ls.CheckString(1)
	newName := ls.CheckString(2)
	if err := os.Rename(oldName, newName); err != nil {
		ls.PushNil()
		ls.PushString(err.Error())
		return 2
	} else {
		ls.PushBoolean(true)
		return 1
	}
}

// os.tmpname ()
// http://www.lua.org/manual/5.3/manual.html#pdf-os.tmpname
func osTmpName(ls LuaState) int {
	panic("todo: osTmpName!")
}

// os.getenv (varname)
// http://www.lua.org/manual/5.3/manual.html#pdf-os.getenv
func osGetEnv(ls LuaState) int {
	key := ls.CheckString(1)
	if env := os.Getenv(key); env != "" {
		ls.PushString(env)
	} else {
		ls.PushNil()
	}
	return 1
}

// os.execute ([command])
// http://www.lua.org/manual/5.3/manual.html#pdf-os.execute
func osExecute(ls LuaState) int {
	panic("todo: osExecute!")
}

// os.exit ([code [, close]])
// http://www.lua.org/manual/5.3/manual.html#pdf-os.exit
func osExit(ls LuaState) int {
	if ls.IsBoolean(1) {
		if ls.ToBoolean(1) {
			os.Exit(0)
		} else {
			os.Exit(1) // todo
		}
	} else {
		code := ls.ToInteger(1)
		os.Exit(int(code))
	}
	return 0
}

// os.setlocale (locale [, category])
// http://www.lua.org/manual/5.3/manual.html#pdf-os.setlocale
func osSetLocale(ls LuaState) int {
	panic("todo: osSetLocale!")
}
