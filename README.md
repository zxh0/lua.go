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

luac -o hello_world.luac test/Pil3/ch01/hello_world.lua
bin/lua hello_world.luac
```
