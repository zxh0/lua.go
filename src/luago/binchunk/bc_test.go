package binchunk

import "encoding/hex"
import "strings"
import "testing"
import "assert"

var fibonacciSrc = `
function fibonacci(n)
  if n < 2 then
    return n
  else
    return fibonacci(n-1) + fibonacci(n-2)
  end
end

local x = fibonacci(16)
print(x)
`

var fibonacciLuacOut = `
1b4c 7561 5300 1993 0d0a 1a0a 0408 0408
0878 5600 0000 0000 0000 0000 0000 2877
4001 1440 7465 7374 5c66 6962 6f6e 6163
6369 2e6c 7561 0000 0000 0000 0000 0002
0309 0000 002c 0000 0008 0000 8006 4040
0041 8000 0024 8000 0146 c040 0080 0000
0064 4000 0126 0080 0004 0000 0004 0a66
6962 6f6e 6163 6369 040a 6669 626f 6e61
6363 6913 1000 0000 0000 0000 0406 7072
696e 7401 0000 0001 0001 0000 0000 0100
0000 0700 0000 0100 040d 0000 0020 0040
001e 4000 8026 0000 011e c001 8046 4040
008e 8040 0064 8000 0186 4040 00ce 0040
00a4 8000 014d 8080 0066 0000 0126 0080
0003 0000 0013 0200 0000 0000 0000 040a
6669 626f 6e61 6363 6913 0100 0000 0000
0000 0100 0000 0000 0000 0000 0d00 0000
0200 0000 0200 0000 0300 0000 0300 0000
0500 0000 0500 0000 0500 0000 0500 0000
0500 0000 0500 0000 0500 0000 0500 0000
0700 0000 0100 0000 026e 0000 0000 0d00
0000 0100 0000 055f 454e 5609 0000 0007
0000 0001 0000 0009 0000 0009 0000 0009
0000 000a 0000 000a 0000 000a 0000 000a
0000 0001 0000 0002 7805 0000 0009 0000
0001 0000 0005 5f45 4e56
	`

var fibonacciLuacOutWithoutDebug = `
1b4c 7561 5300 1993 0d0a 1a0a 0408 0408
0878 5600 0000 0000 0000 0000 0000 2877
4001 0000 0000 0000 0000 0000 0203 0900
0000 2c00 0000 0800 0080 0640 4000 4180
0000 2480 0001 46c0 4000 8000 0000 6440
0001 2600 8000 0400 0000 040a 6669 626f
6e61 6363 6904 0a66 6962 6f6e 6163 6369
1310 0000 0000 0000 0004 0670 7269 6e74
0100 0000 0100 0100 0000 0001 0000 0007
0000 0001 0004 0d00 0000 2000 4000 1e40
0080 2600 0001 1ec0 0180 4640 4000 8e80
4000 6480 0001 8640 4000 ce00 4000 a480
0001 4d80 8000 6600 0001 2600 8000 0300
0000 1302 0000 0000 0000 0004 0a66 6962
6f6e 6163 6369 1301 0000 0000 0000 0001
0000 0000 0000 0000 0000 0000 0000 0000
0000 0000 0000 0000 0000 0000 0000 0000
00
	`

var fibonacciFullList = `
main <test\fibonacci.lua:0,0> (9 instructions)
0+ params, 3 slots, 1 upvalue, 1 local, 4 constants, 1 function
	1	[7]	CLOSURE 	0 0
	2	[1]	SETTABUP	0 -1 0	; _ENV "fibonacci"
	3	[9]	GETTABUP	0 0 -2	; _ENV "fibonacci"
	4	[9]	LOADK   	1 -3	; 16
	5	[9]	CALL    	0 2 2
	6	[10]	GETTABUP	1 0 -4	; _ENV "print"
	7	[10]	MOVE    	2 0
	8	[10]	CALL    	1 2 1
	9	[10]	RETURN  	0 1
constants (4):
	1	"fibonacci"
	2	"fibonacci"
	3	16
	4	"print"
locals (1):
	0	x	6	10
upvalues (1):
	0	_ENV	1	0

function <test\fibonacci.lua:1,7> (13 instructions)
1 param, 4 slots, 1 upvalue, 1 local, 3 constants, 0 functions
	1	[2]	LT      	0 0 -1	; - 2
	2	[2]	JMP     	0 2	; to 5
	3	[3]	RETURN  	0 2
	4	[3]	JMP     	0 8	; to 13
	5	[5]	GETTABUP	1 0 -2	; _ENV "fibonacci"
	6	[5]	SUB     	2 0 -3	; - 1
	7	[5]	CALL    	1 2 2
	8	[5]	GETTABUP	2 0 -2	; _ENV "fibonacci"
	9	[5]	SUB     	3 0 -1	; - 2
	10	[5]	CALL    	2 2 2
	11	[5]	ADD     	1 1 2
	12	[5]	RETURN  	1 2
	13	[7]	RETURN  	0 1
constants (3):
	1	2
	2	"fibonacci"
	3	1
locals (1):
	0	n	1	14
upvalues (1):
	0	_ENV	0	0
`

func TestHeader(t *testing.T) {
	writer := &writer{}
	writer.writeHeader()
	data := writer.data()

	reader := &reader{data}
	reader.readHeader()
}

func TestString(t *testing.T) {
	emptyStr := ""
	shortStr1 := "hello"
	shortStr2 := strings.Repeat("a", 0xFE)
	longStr1 := strings.Repeat("a", 0xFF)
	longStr2 := strings.Repeat("a", 0x100)

	writer := &writer{}
	writer.writeString(emptyStr)
	writer.writeString(shortStr1)
	writer.writeString(shortStr2)
	writer.writeString(longStr1)
	writer.writeString(longStr2)

	reader := &reader{writer.data()}
	assert.StringEqual(t, reader.readString(), emptyStr)
	assert.StringEqual(t, reader.readString(), shortStr1)
	assert.StringEqual(t, reader.readString(), shortStr2)
	assert.StringEqual(t, reader.readString(), longStr1)
	assert.StringEqual(t, reader.readString(), longStr2)
}

func TestFibonacci(t *testing.T) {
	bytes1 := decodeHexStr(fibonacciLuacOut)
	bytes2 := Dump(Undump(bytes1))
	hex1 := hex.EncodeToString(bytes1)
	hex2 := hex.EncodeToString(bytes2)
	assert.StringEqual(t, hex1, hex2)
}

func TestStripDebug(t *testing.T) {
	bytes0 := decodeHexStr(fibonacciLuacOutWithoutDebug)
	bytes1 := decodeHexStr(fibonacciLuacOut)
	bytes2 := Dump(_stripDebug(Undump(bytes1)))
	hex0 := hex.EncodeToString(bytes0)
	hex2 := hex.EncodeToString(bytes2)
	assert.StringEqual(t, hex0, hex2)
}

func TestList(t *testing.T) {
	mainFunc := Undump(decodeHexStr(fibonacciLuacOut))
	output := strings.TrimSpace(List(mainFunc, true))
	expected := strings.TrimSpace(fibonacciFullList)

	if output != expected {
		t.Errorf(output)
	}
}

func decodeHexStr(str string) []byte {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)

	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	return bytes
}

func _stripDebug(proto *FuncProto) *FuncProto {
	StripDebug(proto)
	return proto
}
