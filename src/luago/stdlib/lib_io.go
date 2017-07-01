package stdlib

import . "luago/api"

/*
** functions for 'io' library
 */
// static const luaL_Reg iolib[] = {
//   {"close", io_close},
//   {"flush", io_flush},
//   {"input", io_input},
//   {"lines", io_lines},
//   {"open", io_open},
//   {"output", io_output},
//   {"popen", io_popen},
//   {"read", io_read},
//   {"tmpfile", io_tmpfile},
//   {"type", io_type},
//   {"write", io_write},
//   {NULL, NULL}
// };

/*
** methods for file handles
 */
// static const luaL_Reg flib[] = {
//   {"close", io_close},
//   {"flush", f_flush},
//   {"lines", f_lines},
//   {"read", f_read},
//   {"seek", f_seek},
//   {"setvbuf", f_setvbuf},
//   {"write", f_write},
//   {"__gc", f_gc},
//   {"__tostring", f_tostring},
//   {NULL, NULL}
// };

func OpenIOLib(ls LuaState) int {
	panic("todo!")
	//  luaL_newlib(L, iolib);  /* new module */
	// createmeta(L);
	// /* create (and set) default files */
	// createstdfile(L, stdin, IO_INPUT, "stdin");
	// createstdfile(L, stdout, IO_OUTPUT, "stdout");
	// createstdfile(L, stderr, NULL, "stderr");
	// return 1;
}

// static void createmeta (lua_State *L) {
//   luaL_newmetatable(L, LUA_FILEHANDLE);  /* create metatable for file handles */
//   lua_pushvalue(L, -1);  /* push metatable */
//   lua_setfield(L, -2, "__index");  /* metatable.__index = metatable */
//   luaL_setfuncs(L, flib, 0);  /* add file methods to new metatable */
//   lua_pop(L, 1);  /* pop new metatable */
// }

// static void createstdfile (lua_State *L, FILE *f, const char *k,
//                            const char *fname) {
//   LStream *p = newprefile(L);
//   p->f = f;
//   p->closef = &io_noclose;
//   if (k != NULL) {
//     lua_pushvalue(L, -1);
//     lua_setfield(L, LUA_REGISTRYINDEX, k);  /* add file to registry */
//   }
//   lua_setfield(L, -2, fname);  /* add file to module */
// }

// /*
// ** function to (not) close the standard files stdin, stdout, and stderr
// */
// static int io_noclose (lua_State *L) {
//   LStream *p = tolstream(L);
//   p->closef = &io_noclose;  /* keep file opened */
//   lua_pushnil(L);
//   lua_pushliteral(L, "cannot close standard file");
//   return 2;
// }
