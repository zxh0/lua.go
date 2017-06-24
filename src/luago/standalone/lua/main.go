package main

import "os"
import . "luago/lua"
import "luago/state"

func main() {
	if len(os.Args) > 1 {
		ls := state.NewLuaState()
		ls.OpenLibs()
		if ls.LoadFile(os.Args[1]) == LUA_OK {
			ls.PCall(0, LUA_MULTRET, -1)
		}
	}
}

// func main2() {
// 	argv := os.Args
// 	argc := LuaInteger(len(argv))

// 	ls := NewLuaState()         /* create state */
// 	ls.PushGoFunction(pmain)    /* to call 'pmain' in protected mode */
// 	ls.PushInteger(argc)        /* 1st argument */
// 	ls.PushUserData(argv)       /* 2nd argument */
// 	status := ls.PCall(2, 1, 0) /* do the call */
// 	result := ls.ToBoolean(-1)  /* get result */
// 	println(status)
// 	println(result)
// 	// println(ls.String())
// 	//   report(L, status);
// 	ls.Close()
// 	//   return (result && status == LUA_OK) ? EXIT_SUCCESS : EXIT_FAILURE;
// }

// func pmain(ls LuaState) int {
// 	argc := ls.ToInteger(1)
// 	argv := ls.ToUserData(2).([]string)
// 	println("pnnnnn")
// 	println(argc)
// 	println(argv[0])
// 	// int script;
// 	// int args = collectargs(argv, &script);
// 	ls.CheckVersion() /* check that interpreter has correct version */
// 	// if (argv[0] && argv[0][0]) progname = argv[0];
// 	// if (args == has_error) {  /* bad arg? */
// 	//   print_usage(argv[script]);  /* 'script' has index of bad arg. */
// 	//   return 0;
// 	// }
// 	// if (args & has_v)  /* option '-v'? */
// 	//   print_version();
// 	// if (args & has_E) {  /* option '-E'? */
// 	//   lua_pushboolean(L, 1);  /* signal for libraries to ignore env. vars. */
// 	//   lua_setfield(L, LUA_REGISTRYINDEX, "LUA_NOENV");
// 	// }
// 	// luaL_openlibs(L);  /* open standard libraries */
// 	// createargtable(L, argv, argc, script);  /* create table 'arg' */
// 	// if (!(args & has_E)) {   no option '-E'?
// 	//   if (handle_luainit(L) != LUA_OK)  /* run LUA_INIT */
// 	//     return 0;  /* error running LUA_INIT */
// 	// }
// 	// if (!runargs(L, argv, script))  /* execute arguments -e and -l */
// 	//   return 0;  /* something failed */
// 	// if (script < argc &&  /* execute main script (if there is one) */
// 	//     handle_script(L, argv + script) != LUA_OK)
// 	//   return 0;
// 	// if (args & has_i)  /* -i option? */
// 	//   doREPL(L);  /* do read-eval-print loop */
// 	// else if (script == argc && !(args & (has_e | has_v))) {  /* no arguments? */
// 	//   if (lua_stdin_is_tty()) {  /* running in interactive mode? */
// 	//     print_version();
// 	//     doREPL(L);  /* do read-eval-print loop */
// 	//   }
// 	//   else dofile(L, NULL);  /* executes stdin as a file */
// 	// }
// 	// lua_pushboolean(L, 1);  /* signal no errors */
// 	//ls.PushBoolean(true)
// 	return 1
// }
