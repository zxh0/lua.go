#!/bin/sh
set -ex

go build github.com/zxh0/lua.go/cmd/lua
./lua ./test/lua-5.3.4-tests/attrib.lua     | grep -q OK
./lua ./test/lua-5.3.4-tests/bitwise.lua    | grep -q OK
./lua ./test/lua-5.3.4-tests/calls.lua      | grep -q OK
./lua ./test/lua-5.3.4-tests/closure.lua    | grep -q OK
./lua ./test/lua-5.3.4-tests/constructs.lua | grep -q OK
./lua ./test/lua-5.3.4-tests/events.lua     | grep -q OK
./lua ./test/lua-5.3.4-tests/goto.lua       | grep -q OK
./lua ./test/lua-5.3.4-tests/locals.lua     | grep -q OK
./lua ./test/lua-5.3.4-tests/nextvar.lua    | grep -q OK
./lua ./test/lua-5.3.4-tests/math.lua       | grep -q OK
./lua ./test/lua-5.3.4-tests/sort.lua       | grep -q OK
./lua ./test/lua-5.3.4-tests/strings.lua    | grep -q OK
./lua ./test/lua-5.3.4-tests/utf8.lua       | grep -q ok
./lua ./test/lua-5.3.4-tests/vararg.lua     | grep -q OK
./lua ./test/lua-5.3.4-tests/verybig.lua    | grep -q OK
echo "OK"