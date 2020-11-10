#!/usr/bin/env bash
set -e

alias luac54=/Users/zxh/me/lua-5.4.1/luac
alias luacgo='go run github.com/zxh0/lua.go/cmd/luac'
alias luago='go run github.com/zxh0/lua.go/cmd/lua'
luac54 -v | grep 'Lua 5.4.1'

# test decoder
# for f in test/lua-5.4.1-tests/*.lua ; do
#   echo $f
#   luac54 $f
#   luacgo -l luac.out
# done

# test hw&ops
# luac54 test/hello_world.lua && luago luac.out | grep "Hello, World!"
# luac54 test/ops/0x00_move.lua && luago luac.out | grep "OK"
# luac54 test/ops/0x01_load.lua && luago luac.out | grep "123"
# luac54 test/ops/0x06_lfalseskip.lua && luago luac.out | grep "false"
# luac54 test/ops/0x0b_table.lua && luago luac.out | tr '\n' ';' | grep "123;456"
# luac54 test/ops/0x15_arith_i.lua && luago luac.out | grep "32"
# luac54 test/ops/0x16_arith_k.lua && luago luac.out | grep "1234567891"
# luac54 test/ops/0x3d_cmp.lua && luago luac.out | tr '\n' ';' | grep "false;false;false;true;true;false"
# luac54 test/ops/0x35_concat.lua && luago luac.out | grep "123"
# luac54 test/ops/0x42_test.lua && luago luac.out | tr '\t' ';' | grep "false;true;false;true"
# luac54 test/ops/0x49_for.lua && luago luac.out | tr '\n' ';' | grep "1;2;3;4;5;6;7;8;9;10;"
# luac54 test/ops/0x4d_tfor.lua && luago luac.out | sort | tr '\t' ':' | tr '\n' ';' | grep 'a:1;b:2;c:3;'

echo ""
echo "OK!"
