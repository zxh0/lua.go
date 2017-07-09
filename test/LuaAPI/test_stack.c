#include <stdio.h>
#include <string.h>
#include "lua.h"
#include "lauxlib.h"
#include "lualib.h"

// gcc -I/usr/local/include -L/usr/local/lib -o test_stack test_stack.c -llua
int main(void) {
  lua_State *L = luaL_newstate();
  lua_pushinteger(L, 35);
  fprintf(stderr, "%d\n", lua_absindex(L, 0));     // 2
  fprintf(stderr, "%d\n", lua_absindex(L, 1));     // 1
  fprintf(stderr, "%d\n", lua_absindex(L, 100));   // 100
  fprintf(stderr, "%d\n", lua_absindex(L, -100));  // -98
  fprintf(stderr, "%d\n", lua_isnumber(L, 0));     // 0
  fprintf(stderr, "%d\n", lua_isnumber(L, 100));   // 0
  fprintf(stderr, "%d\n", lua_isnumber(L, -100));  // 0
  fprintf(stderr, "%d\n", lua_isboolean(L, 100));  // 0
  fprintf(stderr, "%d\n", lua_isboolean(L, -100)); // 0
  fprintf(stderr, "%d\n", lua_type(L, 100));       // -1
  fprintf(stderr, "%d\n", lua_type(L, -100));      // 0
  lua_close(L);
  return 0; 
}
