package main

import "flag"
import "fmt"
import "io/ioutil"
import "os"
import "luago/binchunk"
import "luago/compiler"

const version = `Lua 5.3.3 (lua.go)`
const usage = `
usage: %s [options] [filename]
Available options are:
  -l       list (use -ll for full listing)
  -o name  output to file 'name' (default is "luac.out")
  -p       parse only
  -s       strip debug information
  -v       show version information
`

func main() {
	_l := flag.Bool("l", false, "")
	_ll := flag.Bool("ll", false, "")
	_o := flag.String("o", "luac.out", "")
	_p := flag.Bool("p", false, "")
	_s := flag.Bool("s", false, "")
	_v := flag.Bool("v", false, "")
	flag.Usage = printUsage
	flag.Parse()

	if *_v {
		fmt.Println(version)
		return
	}
	if len(flag.Args()) != 1 {
		printUsage()
		return
	}

	filename := flag.Args()[0]
	proto := loadOrCompile(filename)

	if *_p {
		return
	}
	if *_s {
		binchunk.StripDebug(proto)
	}
	if *_l || *_ll {
		output := binchunk.List(proto, *_ll)
		fmt.Println(output)
	} else {
		// write to disk
		data := binchunk.Dump(proto)
		ioutil.WriteFile(*_o, data, 0644) // todo
	}
}

func printUsage() {
	fmt.Printf(usage, os.Args[0])
}

func loadOrCompile(filename string) *binchunk.Prototype {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if binchunk.IsBinaryChunk(data) {
		return binchunk.Undump(data)
	} else {
		return compiler.Compile(string(data), "@"+filename)
	}
}
