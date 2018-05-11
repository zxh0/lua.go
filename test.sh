#!/bin/sh
set -ex

go install luago/standalone/lua
./bin/lua ./test/Lua534TestSuites/bitwise.lua    | grep -q OK
./bin/lua ./test/Lua534TestSuites/calls.lua      | grep -q OK
./bin/lua ./test/Lua534TestSuites/closure.lua    | grep -q OK
./bin/lua ./test/Lua534TestSuites/constructs.lua | grep -q OK
./bin/lua ./test/Lua534TestSuites/events.lua     | grep -q OK
./bin/lua ./test/Lua534TestSuites/locals.lua     | grep -q OK
./bin/lua ./test/Lua534TestSuites/nextvar.lua    | grep -q OK
./bin/lua ./test/Lua534TestSuites/math.lua       | grep -q OK
./bin/lua ./test/Lua534TestSuites/sort.lua       | grep -q OK
./bin/lua ./test/Lua534TestSuites/strings.lua    | grep -q OK
./bin/lua ./test/Lua534TestSuites/utf8.lua       | grep -q ok
./bin/lua ./test/Lua534TestSuites/vararg.lua     | grep -q OK
./bin/lua ./test/Lua534TestSuites/verybig.lua    | grep -q OK