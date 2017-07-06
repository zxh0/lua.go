# lua.go
A TOY Lua 5.3 implementation written in Go (WORK IN PROGRESS).

![lua.go Logo](https://github.com/zxh0/lua.go/blob/master/logo.png?raw=true)

# Build & Test
```shell
cd .
git clone https://github.com/zxh0/lua.go.git

cd lua.go
export GOPATH=`pwd`
go install luago/standalone/lua

luac -o hello_world.luac test/Pil4/ch01/hello_world.lua
bin/lua hello_world.luac
```

# Links
* [Lua 5.3 Reference Manual](http://www.lua.org/manual/5.3/manual.html)
* [Lua 5.3 Source Code](http://www.lua.org/ftp/lua-5.3.4.tar.gz)
* [Lua 5.3 Test suites](http://www.lua.org/tests/lua-5.3.4-tests.tar.gz)
* [Lua 5.3 Bytecode Reference](http://the-ravi-programming-language.readthedocs.io/en/latest/lua_bytecode_reference.html#lua-5-3-bytecode-reference)
* [A No-Frills Introduction to Lua 5.1 VM Instructions](http://luaforge.net/docman/83/98/ANoFrillsIntroToLua51VMInstructions.pdf)
* [The Evolution of Lua](http://www.lua.org/doc/hopl.pdf)
* [The Implementation of Lua 5.0](http://www.lua.org/doc/jucs05.pdf)
* [Syntax-Diagrams for Lua 5.0](http://lua.lickert.net/syntax/Lua_syntax.pdf)
* [Lua-Source-Internal](https://github.com/lichuang/Lua-Source-Internal)
* [Linear Scan Register Allocation](http://www.cs.ucla.edu/~palsberg/course/cs132/linearscan.pdf)
* [Programming in Lua](https://www.lua.org/pil/)
* https://www.0value.com/implementing-lua-coroutines-in-go
