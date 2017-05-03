#! bash
# luac5.3 -o $f"c" $f

ll=$1"c.ll.txt"
echo $ll
luac5.3 -l -l $1 > $ll
sed -i.bak -E 's/ (at|for) 0x[0-9a-f]+//g' $ll
sed -i.bak -E 's/( |	); 0x[0-9a-f]+//g' $ll

rm luac.out
rm *.bak
