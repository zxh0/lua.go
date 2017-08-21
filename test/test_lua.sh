#!/bin/sh

echo "compile lua.go ..."
export GOPATH=`cd ..; pwd`
go install luago/standalone/lua

lua53=lua
luago=`pwd`/../bin/lua
filename=$1

if [ ! -z $2 ]; then
  keep_output=true
else
  keep_output=false
fi

test_file() { # $1:dir $2:file
  printf "[test] $1/$2"

  lua53_output=$2".output"
  luago_output=$2"go.output"

  $lua53 $2 > $lua53_output 2>&1
  $luago $2 > $luago_output 2>&1

  x=`diff -B  $lua53_output $luago_output`
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
test_dir "PiL4/ch08"
test_dir "PiL4/ch09"
test_dir "PiL4/ch10"
test_dir "PiL4/ch24"
test_dir "PiL4/ch25"
