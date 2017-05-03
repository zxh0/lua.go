package compiler

// import "fmt"
import "io/ioutil"
import "strings"
import "testing"
import "luago/binchunk"

const luaGoDir = "../../../../" // lua.go
const pil3Dir = luaGoDir + "test/PiL3/" // lua.go/PiL3
const syntaxDir = luaGoDir + "test/syntax/"
const llPostfix = "c.ll.txt"

func TestPiL3(t *testing.T) {
	testLuaFiles(t, pil3Dir + "ch01/",
		// "hello_world.lua",
		// "factorial.lua",
		// "lib1.lua",
	)
	testLuaFiles(t, pil3Dir + "ch03/",
		// "is_turnback1.lua",
		// "factorial.lua",
		// "lib1.lua",
	)
	testLuaFiles(t, syntaxDir,
		"stats.lua",
	)
}

func _TestCompile(t *testing.T) {
	testLuaFiles(t, luaGoDir + "test/",
		// "blank.lua",
		// "hello_world.lua",
		// "for_num.lua",
		// "func_call.lua",
		// "if.lua",
		"binop.lua",
	)
}

func testLuaFiles(t *testing.T, testDir string, luaFiles ...string) {
	for _, luaFile := range luaFiles {
		testLuaFile(t, testDir, luaFile)
	}
}

func testLuaFile(t *testing.T, testDir, luaFile string) {
	luaSrc := readFile(testDir + luaFile)
	fullList := readList(testDir + luaFile + llPostfix)

	proto := Compile(luaSrc)
	listOutput := binchunk.List(proto, true)
	listOutput = strings.TrimSpace(listOutput)

	if listOutput != fullList {
		//fmt.Printf("%q\n", listOutput)
		//fmt.Printf("%q\n", fullList)
		//t.Errorf(luaFile)
		println(listOutput)
		panic(luaFile)
	}
}

func readList(filename string) string {
	list := readFile(filename)
	list = strings.TrimSpace(list)
	list = strings.Replace(list, "\r\n", "\n", -1)
	return list
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err == nil {
		return string(bytes)
	} else {
		panic(err)
	}
}
