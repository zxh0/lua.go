#!/bin/sh

echo "compile lua.go ..."
export GOPATH=`cd ..; pwd`
go install luago/standalone/lua

lua=lua5.3
luago=`pwd`/../bin/lua
filename=$1

if [ ! -z $2 ]; then
  keep_output=true
else
  keep_output=false
fi

test_file() { # $1:dir $2:file
  printf "[test] $1/$2"

  lua_output=$2".output"
  luago_output=$2"go.output"

  $lua $2 > $lua_output 2>&1
  $luago $2 > $luago_output 2>&1

  x=`diff -B  $lua_output $luago_output`
  if [ ! -z "$x" ]; then
    echo " !!!"
  else
    echo ""
  fi

  if ! $keep_output; then
    rm *.output
  fi
}

test_dir() { # $1:dir
  old_dir=`pwd`
  cd $1

  for f in *.lua; do
    if [[ "$f" == *"$filename"* ]]; then
      if [[ "$f" == _* ]]; then
        echo "[skip] $1/$f"
      else
        test_file $1 $f
      fi
    fi
  done

  cd $old_dir
}

test_dir "PiL4/ch01"
test_dir "PiL4/ch03"
test_dir "PiL4/ch04"
test_dir "PiL4/ch05"
test_dir "PiL4/ch06"
