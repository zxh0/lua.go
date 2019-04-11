package stdlib

import . "luago/api"

/* functions for 'io' library */
var ioLib = map[string]GoFunction{
	"close":   ioClose,
	"flush":   ioFlush,
	"input":   ioInput,
	"lines":   ioLines,
	"open":    ioOpen,
	"output":  ioOutput,
	"popen":   ioPopen,
	"read":    ioRead,
	"tmpfile": ioTmpFile,
	"type":    ioType,
	"write":   ioWrite,
}

func OpenIOLib(ls LuaState) int {
	ls.NewLib(ioLib) /* new module */
	//createmeta(L);
	/* create (and set) default files */
	//createstdfile(L, stdin, IO_INPUT, "stdin");
	//createstdfile(L, stdout, IO_OUTPUT, "stdout");
	//createstdfile(L, stderr, NULL, "stderr");
	return 1
}

func ioClose(ls LuaState) int {
	panic("todo: ioClose!")
}

func ioFlush(ls LuaState) int {
	panic("todo: ioFlush!")
}

func ioInput(ls LuaState) int {
	panic("todo: ioInput!")
}

func ioLines(ls LuaState) int {
	panic("todo: ioLines!")
}

func ioOpen(ls LuaState) int {
	panic("todo: ioOpen!")
}

func ioOutput(ls LuaState) int {
	panic("todo: ioOutput!")
}

func ioPopen(ls LuaState) int {
	panic("todo: ioPopen!")
}

func ioRead(ls LuaState) int {
	panic("todo: ioRead!")
}

func ioTmpFile(ls LuaState) int {
	panic("todo: ioTmpFile!")
}

func ioType(ls LuaState) int {
	panic("todo: ioType!")
}

func ioWrite(ls LuaState) int {
	panic("todo: ioWrite!")
}
