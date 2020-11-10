package binchunk

import (
	"fmt"
	"strings"

	. "github.com/zxh0/lua.go/vm"
)

type printer struct {
	buf []string
}

func (p *printer) printf(format string, a ...interface{}) {
	p.buf = append(p.buf, fmt.Sprintf(format, a...))
}

func (p *printer) printFunc(f *Prototype, full bool) string {
	p.printHeader(f)
	//p.printCode(f)
	//if full {
	//	p.printDebug(f)
	//}
	for _, proto := range f.Protos {
		p.printFunc(proto, full)
	}
	return strings.Join(p.buf, "")
}

func (p *printer) printHeader(f *Prototype) {
	p.printf("\n%s <%s:%d,%d> (%d instruction%s)\n",
		_t(f.LineDefined == 0, "main", "function"),
		_t(f.Source == "", "=?", f.Source)[1:], // TODO
		f.LineDefined, f.LastLineDefined,
		len(f.Code), _s(len(f.Code)),
	)

	p.printf("%d%s param%s, %d slot%s, %d upvalue%s, ",
		f.NumParams, _t(f.IsVararg > 0, "+", ""), _s(int(f.NumParams)),
		f.MaxStackSize, _s(int(f.MaxStackSize)),
		len(f.Upvalues), _s(len(f.Upvalues)),
	)

	p.printf("%d local%s, %d constant%s, %d function%s\n",
		len(f.LocVars), _s(len(f.LocVars)), // TODO
		len(f.Constants), _s(len(f.Constants)),
		len(f.Protos), _s(len(f.Protos)),
	)
}

func (p *printer) printCode(f *Prototype) {
	for pc := 0; pc < len(f.Code); pc++ {
		i := Instruction(f.Code[pc])
		a, b, c := i.ABC()
		_, bx := i.ABx()
		_, sbx := i.AsBx()
		ax := i.Ax()

		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		p.printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName()) // todo

		//switch i.OpMode() {
		//case IABC:
		//	p.printf("%d", a)
		//	if i.BMode() != OpArgN {
		//		if isK(b) {
		//			p.printf(" %d", myk(indexK(b)))
		//		} else {
		//			p.printf(" %d", b)
		//		}
		//	}
		//	if i.CMode() != OpArgN {
		//		if isK(c) {
		//			p.printf(" %d", myk(indexK(c)))
		//		} else {
		//			p.printf(" %d", c)
		//		}
		//	}
		//case IABx:
		//	p.printf("%d", a)
		//	if i.BMode() == OpArgK {
		//		p.printf(" %d", myk(bx))
		//	}
		//	if i.BMode() == OpArgU {
		//		p.printf(" %d", bx)
		//	}
		//case IAsBx:
		//	p.printf("%d %d", a, sbx)
		//case IAx:
		//	p.printf("%d", myk(ax))
		//}

		switch i.Opcode() {
		case OP_LOADK:
			p.printf("\t; ")
			p.printConstant(f, bx)
		case OP_GETUPVAL, OP_SETUPVAL:
			p.printf("\t; %s", upvalName(f, b))
		case OP_GETTABUP:
			p.printf("\t; %s", upvalName(f, b))
			if isK(c) {
				p.printf(" ")
				p.printConstant(f, indexK(c))
			}
		case OP_SETTABUP:
			p.printf("\t; %s", upvalName(f, a))
			if isK(b) {
				p.printf(" ")
				p.printConstant(f, indexK(b))
			}
			if isK(c) {
				p.printf(" ")
				p.printConstant(f, indexK(c))
			}
		case OP_GETTABLE, OP_SELF:
			if isK(c) {
				p.printf("\t; ")
				p.printConstant(f, indexK(c))
			}
		case OP_SETTABLE, OP_ADD, OP_SUB, OP_MUL, OP_POW, OP_DIV, OP_IDIV,
			OP_BAND, OP_BOR, OP_BXOR, OP_SHL, OP_SHR, OP_EQ, OP_LT, OP_LE:
			if isK(b) || isK(c) {
				p.printf("\t; ")
				if isK(b) {
					p.printConstant(f, indexK(b))
				} else {
					p.printf("-")
				}
				p.printf(" ")
				if isK(c) {
					p.printConstant(f, indexK(c))
				} else {
					p.printf("-")
				}
			}
		case OP_JMP, OP_FORLOOP, OP_FORPREP, OP_TFORLOOP:
			p.printf("\t; to %d", sbx+pc+2)
		case OP_CLOSURE:
			// p.printf("\t; %p",VOID(f->p[bx]));
		case OP_SETLIST:
			if c == 0 {
				pc += 1
				p.printf("\t; %d", f.Code[pc])
			} else {
				p.printf("\t; %d", c)
			}
		case OP_EXTRAARG:
			p.printf("\t; ")
			p.printConstant(f, ax)
		}

		p.printf("\n")
	}
}

func (p *printer) printConstant(f *Prototype, i int) {
	k := f.Constants[i]
	switch x := k.(type) {
	case nil:
		p.printf("nil")
	case bool:
		p.printf("%t", x)
	case float64:
		p.printf("%.14g", x) // todo
	case int64:
		p.printf("%d", x) // todo
	case string:
		p.printf("%q", x) // todo
	default: /* cannot happen */
		p.printf("?")
	}
}

func (p *printer) printDebug(f *Prototype) {
	p.printf("constants (%d):\n", len(f.Constants))
	for i, _ := range f.Constants {
		p.printf("\t%d\t", i+1)
		p.printConstant(f, i)
		p.printf("\n")
	}
	p.printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		p.printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}
	p.printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		p.printf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

/* test whether value is a constant */
// #define ISK(x)		((x) & BITRK)
// #define BITRK		(1 << (SIZE_B - 1))
// #define SIZE_B		9
func isK(x int) bool {
	return x&(1<<8) != 0
}

/* gets the index of the constant */
// #define INDEXK(r)	((int)(r) & ~BITRK)
func indexK(r int) int {
	return r & ^(1 << 8)
}

// #define MYK(x)		(-1-(x))
func myk(x int) int {
	return -1 - x
}

// #define UPVALNAME(x) ((f->upvalues[x].name) ? getstr(f->upvalues[x].name) : "-")
func upvalName(f *Prototype, x int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[x]
	} else {
		return "-"
	}
}

// n == 1 ? "" : "s"
func _s(n int) string {
	return _t(n == 1, "", "s")
}

// arg1 ? arg2 : arg3
func _t(arg1 bool, arg2, arg3 string) string {
	if arg1 {
		return arg2
	} else {
		return arg3
	}
}
