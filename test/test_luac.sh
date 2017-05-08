#! bash

export GOPATH=`cd ..; pwd`
go install luago/standalone/luac

luac=luac5.3
luagoc=`pwd`/../bin/luac

test_dir() {
  old_dir=`pwd`
  cd $1

  for f in *.lua; do
    printf "$1/$f"
    luacll=$f"c.ll.txt"
    luagocll=$f"goc.ll.txt"

    # luac -l -l x.lua > x.luac.ll.txt
    $luac -l -l $f > $luacll
    sed -i.bak -E 's/ (at|for) 0x[0-9a-f]+//g' $luacll
    sed -i.bak -E 's/( |	); 0x[0-9a-f]+//g' $luacll
    rm luac.out
    rm *.bak

    # luagoc -ll x.lua > x.luagoc.ll.txt
    $luagoc -ll $f > $luagocll

    x=`diff -B  $luacll $luagocll`
    if [ ! -z "$x" ]; then
      echo " !!!"
    else
      echo ""
    fi

    rm *.ll.txt
  done

  cd $old_dir
}

test_dir "ravi"
test_dir "PiL3/ch01"
test_dir "PiL3/ch02"
test_dir "PiL3/ch03"
# for ch in PiL3/ch*; do
#   test_dir $ch
# done
