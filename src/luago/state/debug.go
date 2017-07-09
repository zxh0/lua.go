package state

import "bytes"
import "fmt"
import "reflect"
import "runtime"
import "strings"
import . "luago/api"

func stackToString(stack *luaStack) string {
	var buf bytes.Buffer

	for i := 0; i < stack.top; i++ {
		buf.WriteString("[")
		buf.WriteString(valToString(stack.slots[i]))
		buf.WriteString("]")
	}

	return buf.String()
}

func valToString(val luaValue) string {
	switch x := val.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%f", val)
	case string:
		return fmt.Sprintf("%q", val)
	case *luaTable:
		return fmt.Sprintf("{@%p}", val)
	case *luaClosure:
		return luaFuncToString(x)
	case *goClosure:
		return goFuncToString(x.goFunc) + "!"
	case GoFunction:
		return goFuncToString(val)
	default:
		fmt.Printf("%T\n", val)
		panic("todo!")
	}
}

func luaFuncToString(luaf *luaClosure) string {
	return fmt.Sprintf("<%s:%d,%d>",
		luaf.proto.Source, // todo
		luaf.proto.LineDefined,
		luaf.proto.LastLineDefined)
}

func goFuncToString(gof luaValue) string {
	pc := reflect.ValueOf(gof).Pointer()
	if f := runtime.FuncForPC(pc); f != nil {
		name := f.Name()[strings.LastIndex(f.Name(), ".")+1:]
		return fmt.Sprintf("%s()", name)
	}
	return fmt.Sprintf("(@%p)", gof)
}
