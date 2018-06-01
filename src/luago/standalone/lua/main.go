package main

import "bufio"
import "fmt"
import "os"
import "strings"
import . "luago/api"

const LUA_INIT_VAR = "LUA_INIT"
const LUA_INITVARVERSION = LUA_INIT_VAR + "_5_3"

const LUA_PROMPT = "> "
const LUA_PROMPT2 = ">> "

/* mark in error messages for incomplete statements */
// const EOFMARK = "<eof>"
const EOFMARK = "'EOF'"

/* bits of various argument indicators in 'args' */
const has_error = 1 /* bad option */
const has_i = 2     /* -i */
const has_v = 4     /* -v */
const has_e = 8     /* -e */
const has_E = 16    /* -E */

var progname = "lua"

// lua-5.3.4/src/lua.c#main()
func main() {
	argc := int64(len(os.Args))
	argv := os.Args

	L := luaL_newstate()            /* create state */
	lua_pushcfunction(L, pmain)     /* to call 'pmain' in protected mode */
	lua_pushinteger(L, argc)        /* 1st argument */
	lua_pushlightuserdata(L, argv)  /* 2nd argument */
	status := lua_pcall(L, 2, 1, 0) /* do the call */
	result := lua_toboolean(L, -1)  /* get result */
	report(L, status)
	lua_close(L)
	// return (result && status == LUA_OK) ? EXIT_SUCCESS : EXIT_FAILURE;
	if result && status != LUA_OK {
		// os.Exit(EXIT_FAILURE)
	}
}

/*
** Main body of stand-alone interpreter (to be called in protected mode).
** Reads the options and handles them all.
 */
func pmain(L LuaState) int {
	argc := int(lua_tointeger(L, 1))
	argv := lua_touserdata(L, 2).([]string)
	var script int
	args := collectargs(argv, &script)
	luaL_checkversion(L) /* check that interpreter has correct version */
	if argv[0] != "" {
		progname = argv[0]
	}
	if args == has_error { /* bad arg? */
		print_usage(argv[script]) /* 'script' has index of bad arg. */
		return 0
	}
	if args&has_v != 0 { /* option '-v'? */
		print_version()
	}
	if args&has_E != 0 { /* option '-E'? */
		lua_pushboolean(L, true) /* signal for libraries to ignore env. vars. */
		lua_setfield(L, LUA_REGISTRYINDEX, "LUA_NOENV")
	}
	luaL_openlibs(L)                      /* open standard libraries */
	createargtable(L, argv, argc, script) /* create table 'arg' */
	if args&has_E == 0 {                  /* no option '-E'? */
		if handle_luainit(L) != LUA_OK { /* run LUA_INIT */
			return 0 /* error running LUA_INIT */
		}
	}
	if !runargs(L, argv, script) { /* execute arguments -e and -l */
		return 0 /* something failed */
	}
	if script < argc && /* execute main script (if there is one) */
		handle_script(L, argv[script:]) != LUA_OK {
		return 0
	}
	if args&has_i != 0 { /* -i option? */
		doREPL(L) /* do read-eval-print loop */
	} else if script == argc && args&(has_e|has_v) == 0 { /* no arguments? */
		if lua_stdin_is_tty() { /* running in interactive mode? */
			print_version()
			doREPL(L) /* do read-eval-print loop */
		} else {
			dofile(L, "")
		} /* executes stdin as a file */
	}
	lua_pushboolean(L, true) /* signal no errors */
	return 1
}

/*
** Traverses all arguments from 'argv', returning a mask with those
** needed before running any Lua code (or an error code if it finds
** any invalid argument). 'first' returns the first not-handled argument
** (either the script name or a bad argument in case of error).
 */
func collectargs(argv []string, first *int) int {
	args := 0
	i := 0
	for i = 1; i < len(argv); i++ {
		*first = i
		if argv[i][0] != '-' { /* not an option? */
			return args /* stop handling options */
		}
		if len(argv[i]) == 1 { /* '-' */
			return args /* script "name" is '-' */
		}
		switch argv[i][1] { /* else check option */
		case '-': /* '--' */
			if len(argv[i]) > 2 { /* extra characters after '--'? */
				return has_error /* invalid option */
			}
			*first = i + 1
			return args
		case 'E':
			if len(argv[i]) > 2 { /* extra characters after 1st? */
				return has_error /* invalid option */
			}
			args |= has_E
			break
		case 'i':
			args |= has_i /* (-i implies -v) */ /* FALLTHROUGH */
			fallthrough
		case 'v':
			if len(argv[i]) > 2 { /* extra characters after 1st? */
				return has_error /* invalid option */
			}
			args |= has_v
			break
		case 'e':
			args |= has_e /* FALLTHROUGH */
			fallthrough
		case 'l': /* both options need an argument */
			if len(argv[i]) == 2 { /* no concatenated argument? */
				i++ /* try next 'argv' */
				if i >= len(argv) || argv[i][0] == '-' {
					return has_error /* no next argument or it is another option */
				}
			}
			break
		default: /* invalid option */
			return has_error
		}
	}
	*first = i /* no script name */
	return args
}

func print_version() {
	fmt.Println(LUA_RELEASE)
}

func print_usage(badoption string) {
	lua_writestringerror("%s: ", progname)
	if badoption[1] == 'e' || badoption[1] == 'l' {
		lua_writestringerror("'%s' needs argument\n", badoption)
	} else {
		lua_writestringerror("unrecognized option '%s'\n", badoption)
		lua_writestringerror("usage: %s [options] [script [args]]\n"+
			"Available options are:\n"+
			"  -e stat  execute string 'stat'\n"+
			"  -i       enter interactive mode after executing 'script'\n"+
			"  -l name  require library 'name'\n"+
			"  -v       show version information\n"+
			"  -E       ignore environment variables\n"+
			"  --       stop handling options\n"+
			"  -        stop handling options and execute stdin\n",
			progname)
	}
}

/*
** Create the 'arg' table, which stores all arguments from the
** command line ('argv'). It should be aligned so that, at index 0,
** it has 'argv[script]', which is the script name. The arguments
** to the script (everything after 'script') go to positive indices;
** other arguments (before the script name) go to negative indices.
** If there is no script name, assume interpreter's name as base.
 */
func createargtable(L LuaState, argv []string, argc, script int) {
	var i, narg int
	if script == argc {
		script = 0
	} /* no script name? */
	narg = argc - (script + 1) /* number of positive indices */
	lua_createtable(L, narg, script+1)
	for i = 0; i < argc; i++ {
		lua_pushstring(L, argv[i])
		lua_rawseti(L, -2, int64(i-script))
	}
	lua_setglobal(L, "arg")
}

func handle_luainit(L LuaState) int {
	name := "=" + LUA_INITVARVERSION
	init := getenv(name[1:])
	if init == "" {
		name = "=" + LUA_INIT_VAR
		init = getenv(name[1:]) /* try alternative name */
	}
	if init == "" {
		return LUA_OK
	} else if init[0] == '@' {
		return dofile(L, init[1:])
	} else {
		return dostring(L, init, name)
	}
	return 0
}

/*
** Processes options 'e' and 'l', which involve running Lua code.
** Returns 0 if some code raises an error.
 */
func runargs(L LuaState, argv []string, n int) bool {
	for i := 1; i < n; i++ {
		option := argv[i][1]
		lua_assert(argv[i][0] == '-') /* already checked */
		if option == 'e' || option == 'l' {
			extra := argv[i][2:] /* both options need an argument */
			if extra == "" {
				i += 1
				extra = argv[i]
			}
			lua_assert(extra != "")
			status := 0
			if option == 'e' {
				status = dostring(L, extra, "=(command line)")
			} else {
				status = dolibrary(L, extra)
			}
			if status != LUA_OK {
				return false
			}
		}
	}
	return true
}

func handle_script(L LuaState, argv []string) int {
	fname := argv[0]
	// if (strcmp(fname, "-") == 0 && strcmp(argv[-1], "--") != 0)
	//   fname = NULL;  /* stdin */
	status := luaL_loadfile(L, fname)
	if status == LUA_OK {
		n := pushargs(L) /* push arguments to script */
		status = docall(L, n, LUA_MULTRET)
	}
	return report(L, status)
}

/*
** Push on the stack the contents of table 'arg' from 1 to #arg
 */
func pushargs(L LuaState) int {
	var i, n int
	if lua_getglobal(L, "arg") != LUA_TTABLE {

		luaL_error(L, "'arg' is not a table")
	}
	n = int(luaL_len(L, -1))
	luaL_checkstack(L, n+3, "too many arguments to script")
	for i = 1; i <= n; i++ {
		lua_rawgeti(L, -i, int64(i))
	}
	lua_remove(L, -i) /* remove table from the stack */
	return n
}

func dofile(L LuaState, name string) int {
	return dochunk(L, luaL_loadfile(L, name))
}

func dostring(L LuaState, s, name string) int {
	return dochunk(L, lua_load(L, []byte(s), name, "bt"))
}

func dochunk(L LuaState, status int) int {
	if status == LUA_OK {
		status = docall(L, 0, 0)
	}
	return report(L, status)
}

/*
** Interface to 'lua_pcall', which sets appropriate message function
** and C-signal handler. Used to run all chunks.
 */
func docall(L LuaState, narg, nres int) int {
	base := lua_gettop(L) - narg     /* function index */
	lua_pushcfunction(L, msghandler) /* push message handler */
	lua_insert(L, base)              /* put it under function and args */
	// globalL = L;  /* to be available to 'laction' */
	// signal(SIGINT, laction);  /* set C-signal handler */
	status := lua_pcall(L, narg, nres, base)
	// signal(SIGINT, SIG_DFL); /* reset C-signal handler */
	lua_remove(L, base) /* remove message handler from the stack */
	return status
}

/*
** Calls 'require(name)' and stores the result in a global variable
** with the given name.
 */
func dolibrary(L LuaState, name string) int {
	lua_getglobal(L, "require")
	lua_pushstring(L, name)
	status := docall(L, 1, 1) /* call 'require(name)' */
	if status == LUA_OK {
		lua_setglobal(L, name) /* global[name] = require return */
	}
	return report(L, status)
}

/*
** Check whether 'status' is not OK and, if so, prints the error
** message on the top of the stack. It assumes that the error object
** is a string, as it was either generated by Lua or by 'msghandler'.
 */
func report(L LuaState, status int) int {
	if status != LUA_OK {
		msg, _ := lua_tostring(L, -1)
		l_message(progname, msg)
		lua_pop(L, 1) /* remove message */
	}
	return status
}

/*
** Do the REPL: repeatedly read (load) a line, evaluate (call) it, and
** print any results.
 */
func doREPL(L LuaState) {
	var status int
	oldprogname := progname
	progname = "" /* no 'progname' on errors in interactive mode */
	for {
		status = loadline(L)
		if status == -1 {
			break
		}
		if status == LUA_OK {
			status = docall(L, 0, LUA_MULTRET)
		}
		if status == LUA_OK {
			l_print(L)
		} else {
			report(L, status)
		}
	}
	lua_settop(L, 0) /* clear stack */
	lua_writeline()
	progname = oldprogname
}

/*
** Read a line and try to load (compile) it first as an expression (by
** adding "return " in front of it) and second as a statement. Return
** the final status of load/call with the resulting function (if any)
** in the top of the stack.
 */
func loadline(L LuaState) int {
	lua_settop(L, 0)
	if !pushline(L, true) {
		return -1 /* no input */
	}
	var status int
	if status = addreturn(L); status != LUA_OK { /* 'return ...' did not work? */
		status = multiline(L) /* try as command, maybe with continuation lines */
	}
	lua_remove(L, 1) /* remove line from the stack */
	lua_assert(lua_gettop(L) == 1)
	return status
}

/*
** Prompt the user, read a line, and push it into the Lua stack.
 */
func pushline(L LuaState, firstline bool) bool {
	prmt := get_prompt(L, firstline)
	b, readstatus := lua_readline(L, prmt)
	if !readstatus {
		return false /* no input (prompt will be popped by caller) */
	}
	lua_pop(L, 1)                             /* remove prompt */
	if l := len(b); l > 0 && b[l-1] == '\n' { /* line ends with newline? */
		b = b[:l-1] /* remove it */
	}
	if firstline && len(b) > 0 && b[0] == '=' { /* for compatibility with 5.2, ... */
		lua_pushfstring(L, "return %s", b[1:]) /* change '=' to 'return' */
	} else {
		lua_pushstring(L, b)
	}
	lua_freeline(L, b)
	return true
}

/*
** Returns the string to be used as a prompt by the interpreter.
 */
func get_prompt(L LuaState, firstline bool) string {
	if firstline {
		lua_getglobal(L, "_PROMPT")
	} else {
		lua_getglobal(L, "_PROMPT2")
	}

	p, _ := lua_tostring(L, -1)
	if p == "" {
		if firstline {
			p = LUA_PROMPT
		} else {
			p = LUA_PROMPT2
		}
	}
	return p
}

/*
** Try to compile line on the stack as 'return <line>;'; on return, stack
** has either compiled chunk or original line (if compilation failed).
 */
func addreturn(L LuaState) int {
	line, _ := lua_tostring(L, -1) /* original line */
	retline := lua_pushfstring(L, "return %s;", line)
	status := lua_load(L, []byte(retline), "=stdin", "t")
	if status == LUA_OK {
		lua_remove(L, -2) /* remove modified line */
		if line != "" {   /* non empty? */
			lua_saveline(L, line) /* keep history */
		}
	} else {
		lua_pop(L, 2) /* pop result from 'luaL_loadbuffer' and modified line */
	}
	return status
}

/*
** Read multiple lines until a complete Lua statement
 */
func multiline(L LuaState) int {
	for { /* repeat until gets a complete statement */
		line, _ := lua_tostring(L, 1)                      /* get what it has */
		status := lua_load(L, []byte(line), "=stdin", "t") /* try it */
		if !incomplete(L, status) || !pushline(L, false) {
			lua_saveline(L, line) /* keep history */
			return status         /* cannot or should not try to add continuation line */
		}
		lua_pushliteral(L, "\n") /* add newline... */
		lua_insert(L, -2)        /* ...between the two lines */
		lua_concat(L, 3)         /* join them */
	}
}

/*
** Check whether 'status' signals a syntax error and the error
** message at the top of the stack ends with the above mark for
** incomplete statements.
 */
func incomplete(L LuaState, status int) bool {
	if status == LUA_ERRSYNTAX {
		msg, _ := lua_tostring(L, -1)
		if strings.HasSuffix(msg, EOFMARK) {
			lua_pop(L, 1)
			return true
		}
	}
	return false /* else... */
}

/*
** Prints (calling the Lua 'print' function) any values on the stack
 */
func l_print(L LuaState) {
	n := lua_gettop(L)
	if n > 0 { /* any result to be printed? */
		luaL_checkstack(L, LUA_MINSTACK, "too many results to print")
		lua_getglobal(L, "print")
		lua_insert(L, 1)
		if lua_pcall(L, n, 0, 0) != LUA_OK {
			l_message(progname, lua_pushfstring(L, "error calling 'print' (%s)",
				lua_tostring2(L, -1)))
		}
	}
}

func lua_writeline() {
	fmt.Println()
}

func lua_readline(L LuaState, prmt string) (string, bool) {
	fmt.Print(prmt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	return text, err == nil
}

func lua_saveline(L LuaState, line string) {
	// todo
}

func lua_freeline(L LuaState, line string) {
	// todo
}

/*
** Prints an error message, adding the program name in front of it
** (if present)
 */
func l_message(pname, msg string) {
	if pname != "" {
		lua_writestringerror("%s: ", pname)
	}
	lua_writestringerror("%s\n", msg)
}

func lua_writestringerror(s, p string) {
	fmt.Fprintf(os.Stderr, s, p)
}

func getenv(name string) string {
	return "" // todo
}

func lua_stdin_is_tty() bool {
	return true // todo
}

func lua_assert(c bool) {
	// todo
}

/*
** Message handler used to run all chunks
 */
func msghandler(L LuaState) int {
	msg, isStr := lua_tostring(L, 1)
	if !isStr { /* is error object not a string? */
		if luaL_callmeta(L, 1, "__tostring") && /* does it have a metamethod */
			lua_type(L, -1) == LUA_TSTRING { /* that produces a string? */
			return 1 /* that is the message */
		} else {
			msg = lua_pushfstring(L, "(error object is a %s value)",
				luaL_typename(L, 1))
		}
	}
	luaL_traceback(L, L, msg, 1) /* append a standard traceback */
	return 1                     /* return the traceback */
}
