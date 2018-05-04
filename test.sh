#!/bin/sh
set -ex

go install luago/standalone/lua
./bin/lua ./test/Lua534TestSuites/bitwise.lua    | grep OK
./bin/lua ./test/Lua534TestSuites/closure.lua    | grep OK
./bin/lua ./test/Lua534TestSuites/constructs.lua | grep OK
./bin/lua ./test/Lua534TestSuites/events.lua     | grep OK
./bin/lua ./test/Lua534TestSuites/locals.lua     | grep OK
./bin/lua ./test/Lua534TestSuites/nextvar.lua    | grep -q OK
./bin/lua ./test/Lua534TestSuites/math.lua       | grep OK
./bin/lua ./test/Lua534TestSuites/vararg.lua     | grep OK
