#!/bin/sh

export GOPATH=`cd ..; pwd`
go install luago/standalone/luac

luac=luac5.3
luacgo=`pwd`/../bin/luac
filename=$1

if [ ! -z $2 ]; then
  keep_ll=true
else
  keep_ll=false
fi

test_file() { # $1:dir $2:file
  printf "$1/$2"

  luacll=$2"c.ll.txt"
  luacgoll=$2"goc.ll.txt"

  # luac -l -l x.lua > x.luac.ll.txt
  $luac -l -l $2 > $luacll
  sed -i.bak -E 's/ (at|for) 0x[0-9a-f]+//g' $luacll
  sed -i.bak -E 's/( |	); 0x[0-9a-f]+//g' $luacll
  rm luac.out
  rm *.bak

  # luacgo -ll x.lua > x.luacgo.ll.txt
  $luacgo -ll $2 > $luacgoll

  x=`diff -B  $luacll $luacgoll`
  if [ ! -z "$x" ]; then
    echo " !!!"
  else
    echo ""
  fi

  if ! $keep_ll; then
    rm *.ll.txt
  fi
}

test_dir() { # $1:dir
  old_dir=`pwd`
  cd $1

  for f in *.lua; do
    if [[ "$f" == *"$filename"* ]]; then
      if [[ "$f" == _* ]]; then
        echo "$1/$f ?"
      else
        test_file $1 $f
      fi
    fi
  done

  cd $old_dir
}

test_dir "ravi"
test_dir "syntax"
test_dir "PiL3/ch01"
test_dir "PiL3/ch02"
test_dir "PiL3/ch03"
# for ch in PiL3/ch*; do
#   test_dir $ch
# done
