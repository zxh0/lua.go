package stdlib

import . "luago/api"

/* key, in the registry, for table of loaded modules */
const LUA_LOADED_TABLE = "_LOADED"

/* key, in the registry, for table of preloaded loaders */
const LUA_PRELOAD_TABLE = "_PRELOAD"

const LUA_PATH_VAR = "LUA_PATH"
const LUA_PATHSUFFIX = "_" + LUA_VERSION_MAJOR + "_" + LUA_VERSION_MINOR
const LUA_PATHVARVERSION = LUA_PATH_VAR + LUA_PATHSUFFIX

var _GOLIBS = "golibs"

var pkgFuncs = map[string]LuaGoFunction{
	"loadlib":    pkgLoadLib,
	"searchpath": pkgSearchPath,
	/* placeholders */
	"preload":   nil,
	"cpath":     nil,
	"path":      nil,
	"searchers": nil,
	"loaded":    nil,
}

var llFuncs = map[string]LuaGoFunction{
	"require": llRequire,
}

func OpenPackageLib(ls LuaState) int {
	//createGoLibsTable(ls)
	ls.NewLib(pkgFuncs) /* create 'package' table */
	createSearchersTable(ls)
	/* set field 'path' */
	// setPath(ls, "path", LUA_PATHVARVERSION, LUA_PATH_VAR, LUA_PATH_DEFAULT)
	// /* set field 'cpath' */
	// setpath(L, "cpath", LUA_CPATHVARVERSION, LUA_CPATH_VAR, LUA_CPATH_DEFAULT);
	// /* store config information */
	// lua_pushliteral(L, LUA_DIRSEP "\n" LUA_PATH_SEP "\n" LUA_PATH_MARK "\n"
	//                    LUA_EXEC_DIR "\n" LUA_IGMARK "\n");
	// ls.Setfield(L, -2, "config");
	// /* set field 'loaded' */
	ls.GetSubTable(LUA_REGISTRYINDEX, LUA_LOADED_TABLE)
	ls.SetField(-2, "loaded")
	/* set field 'preload' */
	ls.GetSubTable(LUA_REGISTRYINDEX, LUA_PRELOAD_TABLE)
	ls.SetField(-2, "preload")
	ls.PushGlobalTable()
	ls.PushValue(-2)        /* set 'package' as upvalue for next lib */
	ls.SetFuncs(llFuncs, 1) /* open lib into global table */
	ls.Pop(1)               /* pop global table */
	return 1                /* return 'package' table */
}

/*
 * create table GOLIBS to keep track of loaded Go libraries,
 * setting a finalizer to close all libraries when closing state.
 */
func createGoLibsTable(ls LuaState) {
	ls.NewTable() /* create GOLIBS table */
	// lua_createtable(L, 0, 1);  /* create metatable for CLIBS */
	// lua_pushcfunction(L, gctm);
	// ls.Setfield(L, -2, "__gc");  /* set finalizer for GOLIBS table */
	// ls.Setmetatable(L, -2);
	ls.RawSetP(LUA_REGISTRYINDEX, &_GOLIBS) /* set GOLIBS table in registry */
}

func createSearchersTable(ls LuaState) {
	searchers := []LuaGoFunction{
		searcherPreload,
		searcherLua,
		searcherGo,
		searcherGoRoot,
	}
	/* create 'searchers' table */
	ls.CreateTable(len(searchers), 0)
	/* fill it with predefined searchers */
	for idx, searcher := range searchers {
		ls.PushValue(-2) /* set 'package' as upvalue for all searchers */
		ls.PushGoClosure(searcher, 1)
		ls.RawSetI(-2, int64(idx+1))
	}
	ls.SetField(-2, "searchers") /* put it in field 'searchers' */
}

func searcherPreload(ls LuaState) int {
	name := ls.CheckString(1)
	ls.GetField(LUA_REGISTRYINDEX, "_PRELOAD")
	if ls.GetField(-1, name) == LUA_TNIL { /* not found? */
		// lua_pushfstring(L, "\n\tno field package.preload['%s']", name);
		ls.PushString("\n\tno field package.preload['" + name + "']")
	}
	return 1
}

func searcherLua(ls LuaState) int {
	panic("todo: searcherLua")
	// const char *filename;
	// const char *name = ls.checkstring(L, 1);
	// filename = findfile(L, name, "path", LUA_LSUBSEP);
	// if (filename == NULL) return 1;  /* module not found in this path */
	// return checkload(L, (ls.loadfile(L, filename) == LUA_OK), filename);
}

func searcherGo(ls LuaState) int {
	panic("todo: searcherGo")
	// const char *name = ls.checkstring(L, 1);
	// const char *filename = findfile(L, name, "cpath", LUA_CSUBSEP);
	// if (filename == NULL) return 1;  /* module not found in this path */
	// return checkload(L, (loadfunc(L, filename, name) == 0), filename);
}

func searcherGoRoot(ls LuaState) int {
	panic("todo: searcherGoRoot")
	// const char *filename;
	// const char *name = ls.checkstring(L, 1);
	// const char *p = strchr(name, '.');
	// int stat;
	// if (p == NULL) return 0;  /* is root */
	// lua_pushlstring(L, name, p - name);
	// filename = findfile(L, lua_tostring(L, -1), "cpath", LUA_CSUBSEP);
	// if (filename == NULL) return 1;  /* root not found */
	// if ((stat = loadfunc(L, filename, name)) != 0) {
	//   if (stat != ERRFUNC)
	//     return checkload(L, 0, filename);  /* real error */
	//   else {  /* open function not found */
	//     lua_pushfstring(L, "\n\tno module '%s' in file '%s'", name, filename);
	//     return 1;
	//   }
	// }
	// lua_pushstring(L, filename);  /* will be 2nd argument to module */
	// return 2;
}

func setPath(ls LuaState, fieldName, envName1, envName2, def string) {
	println(envName1)
	panic("todo setpath")
	// const char *path = getenv(envname1);
	// if (path == NULL)  /* no environment variable? */
	//   path = getenv(envname2);  /* try alternative name */
	// if (path == NULL || noenv(L))  /* no environment variable? */
	//   lua_pushstring(L, def);  /* use default */
	// else {
	//    replace ";;" by ";AUXMARK;" and then AUXMARK by default path
	//   path = ls.gsub(L, path, LUA_PATH_SEP LUA_PATH_SEP,
	//                             LUA_PATH_SEP AUXMARK LUA_PATH_SEP);
	//   ls.gsub(L, path, AUXMARK, def);
	//   lua_remove(L, -2);
	// }
	// setprogdir(L);
	// ls.Setfield(L, -2, fieldname);
}

func pkgLoadLib(ls LuaState) int {
	panic("todo pkgLoadLib")
}

func pkgSearchPath(ls LuaState) int {
	panic("todo pkgSearchPath")
}

func llRequire(ls LuaState) int {
	name := ls.CheckString(1)
	ls.SetTop(1) /* LOADED table will be at index 2 */
	ls.GetField(LUA_REGISTRYINDEX, LUA_LOADED_TABLE)
	ls.GetField(2, name)  /* LOADED[name] */
	if ls.ToBoolean(-1) { /* is it there? */
		return 1 /* package is already loaded */
	}
	/* else must load package */
	ls.Pop(1) /* remove 'getfield' result */
	findLoader(ls, name)
	ls.PushString(name) /* pass name as argument to module loader */
	ls.Insert(-2)       /* name is 1st argument (before search data) */
	ls.Call(2, 1)       /* run loader to load module */
	if !ls.IsNil(-1) {  /* non-nil return? */
		ls.SetField(2, name) /* LOADED[name] = returned value */
	}
	if ls.GetField(2, name) == LUA_TNIL { /* module set no value? */
		ls.PushBoolean(true) /* use true as result */
		ls.PushValue(-1)     /* extra copy to be returned */
		ls.SetField(2, name) /* LOADED[name] = true */
	}
	return 1
}

func findLoader(ls LuaState, name string) {
	// int i;
	// luaL_Buffer msg;  /* to build error message */
	// luaL_buffinit(L, &msg);
	/* push 'package.searchers' to index 3 in the stack */
	if ls.GetField(LuaUpvalueIndex(1), "searchers") != LUA_TTABLE {
		//luaL_error(L, "'package.searchers' must be a table");
		panic("'package.searchers' must be a table") // todo
	}
	/*  iterate over available searchers to find a loader */
	for i := 1; ; i++ {
		if ls.RawGetI(3, int64(i)) == LUA_TNIL { /* no more searchers? */
			ls.Pop(1) /* remove nil */
			//luaL_pushresult(&msg);   /* create error message */
			//luaL_error(L, "module '%s' not found:%s", name, lua_tostring(L, -1));
			panic("todo!")
		}
		ls.PushString(name)
		ls.Call(1, 2)          /* call it */
		if ls.IsFunction(-2) { /* did it find a loader? */
			return /* module loader found */
		} else if ls.IsString(-2) { /* searcher returned error message? */
			ls.Pop(1) /* remove extra return */
			// luaL_addvalue(&msg);  /* concatenate error message */
		} else {
			ls.Pop(2) /* remove both returns */
		}
	}
}
