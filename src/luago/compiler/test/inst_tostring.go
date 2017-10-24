package test

import "fmt"
import "strings"
import . "luago/vm"

func instToStr(_i uint32) string {
	i := Instruction(_i)
	opName := strings.ToLower(i.OpName())
	opName = strings.TrimSpace(opName)

	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		return fmt.Sprintf("%s(%d,%s,%s)", opName, a,
			argBCToStr(b, i.BMode()),
			argBCToStr(c, i.CMode()))
	case IABx:
		a, bx := i.ABx()
		return fmt.Sprintf("%s(%d,%s)", opName, a,
			argBxToStr(bx, i.BMode()))
	case IAsBx:
		a, sbx := i.AsBx()
		return fmt.Sprintf("%s(%d,%d)", opName, a, sbx)
	case IAx:
		ax := i.Ax()
		return fmt.Sprintf("%s(%d)", opName, -1-ax)
	default:
		panic("unreachable!")
	}
}

func argBCToStr(arg int, mode byte) string {
	if mode == OpArgN {
		return "_"
	}
	if arg > 0xFF {
		return fmt.Sprintf("%d", -1-arg&0xFF)
	}
	return fmt.Sprintf("%d", arg)
}

func argBxToStr(bx int, mode byte) string {
	if mode == OpArgK {
		return fmt.Sprintf("%d", -1-bx)
	}
	if mode == OpArgU {
		return fmt.Sprintf("%d", bx)
	}
	return "_"
}
