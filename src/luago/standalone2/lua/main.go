package main

import "fmt"
import "os"
import . "luago/api"

/* bits of various argument indicators in 'args' */
const has_error = 1 /* bad option */
const has_i = 2     /* -i */
const has_v = 4     /* -v */
const has_e = 8     /* -e */
const has_E = 16    /* -E */

const LUA_INIT_VAR = "LUA_INIT"
const LUA_INITVARVERSION = LUA_INIT_VAR + "_5_3"

var progname = "lua"

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
		//   doREPL(L);  /* do read-eval-print loop */
	} else if script == argc && args&(has_e|has_v) == 0 { /* no arguments? */
		if lua_stdin_is_tty() { /* running in interactive mode? */
			print_version()
			//     doREPL(L);  /* do read-eval-print loop */
			println("doREPL...")
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
		//lua_assert(argv[i][0] == '-');  /* already checked */
		if option == 'e' || option == 'l' {
			extra := argv[i][2:] /* both options need an argument */
			if extra == "" {
				i += 1
				extra = argv[i]
			}
			//lua_assert(extra != NULL);
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

/*
** Message handler used to run all chunks
 */
func msghandler(L LuaState) int {
	// const char *msg = lua_tostring(L, 1);
	// if (msg == NULL) {  /* is error object not a string? */
	//   if (luaL_callmeta(L, 1, "__tostring") &&  /* does it have a metamethod */
	//       lua_type(L, -1) == LUA_TSTRING)   that produces a string?
	//     return 1;  /* that is the message */
	//   else
	//     msg = lua_pushfstring(L, "(error object is a %s value)",
	//                              luaL_typename(L, 1));
	// }
	// luaL_traceback(L, L, msg, 1);  /* append a standard traceback */
	println("msghandler...")
	return 1 /* return the traceback */
}
