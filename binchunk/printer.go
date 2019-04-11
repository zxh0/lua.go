package binchunk

import "fmt"
import "strings"
import . "github.com/zxh0/lua.go/vm"

type printer struct {
	buf []string
}

func (self *printer) printf(format string, a ...interface{}) {
	self.buf = append(self.buf, fmt.Sprintf(format, a...))
}

func (self *printer) printFunc(f *Prototype, full bool) string {
	self.printHeader(f)
	self.printCode(f)
	if full {
		self.printDebug(f)
	}
	for _, p := range f.Protos {
		self.printFunc(p, full)
	}
	return strings.Join(self.buf, "")
}

func (self *printer) printHeader(f *Prototype) {
	self.printf("\n%s <%s:%d,%d> (%d instruction%s)\n",
		_t(f.LineDefined == 0, "main", "function"),
		_t(f.Source == "", "=?", f.Source)[1:], // todo
		f.LineDefined, f.LastLineDefined,
		len(f.Code), _s(len(f.Code)),
	)

	self.printf("%d%s param%s, %d slot%s, %d upvalue%s, ",
		f.NumParams, _t(f.IsVararg > 0, "+", ""), _s(int(f.NumParams)),
		f.MaxStackSize, _s(int(f.MaxStackSize)),
		len(f.Upvalues), _s(len(f.Upvalues)),
	)

	self.printf("%d local%s, %d constant%s, %d function%s\n",
		len(f.LocVars), _s(len(f.LocVars)), // todo
		len(f.Constants), _s(len(f.Constants)),
		len(f.Protos), _s(len(f.Protos)),
	)
}

func (self *printer) printCode(f *Prototype) {
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
		self.printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName()) // todo

		switch i.OpMode() {
		case IABC:
			self.printf("%d", a)
			if i.BMode() != OpArgN {
				if isK(b) {
					self.printf(" %d", myk(indexK(b)))
				} else {
					self.printf(" %d", b)
				}
			}
			if i.CMode() != OpArgN {
				if isK(c) {
					self.printf(" %d", myk(indexK(c)))
				} else {
					self.printf(" %d", c)
				}
			}
		case IABx:
			self.printf("%d", a)
			if i.BMode() == OpArgK {
				self.printf(" %d", myk(bx))
			}
			if i.BMode() == OpArgU {
				self.printf(" %d", bx)
			}
		case IAsBx:
			self.printf("%d %d", a, sbx)
		case IAx:
			self.printf("%d", myk(ax))
		}

		switch i.Opcode() {
		case OP_LOADK:
			self.printf("\t; ")
			self.printConstant(f, bx)
		case OP_GETUPVAL, OP_SETUPVAL:
			self.printf("\t; %s", upvalName(f, b))
		case OP_GETTABUP:
			self.printf("\t; %s", upvalName(f, b))
			if isK(c) {
				self.printf(" ")
				self.printConstant(f, indexK(c))
			}
		case OP_SETTABUP:
			self.printf("\t; %s", upvalName(f, a))
			if isK(b) {
				self.printf(" ")
				self.printConstant(f, indexK(b))
			}
			if isK(c) {
				self.printf(" ")
				self.printConstant(f, indexK(c))
			}
		case OP_GETTABLE, OP_SELF:
			if isK(c) {
				self.printf("\t; ")
				self.printConstant(f, indexK(c))
			}
		case OP_SETTABLE, OP_ADD, OP_SUB, OP_MUL, OP_POW, OP_DIV, OP_IDIV,
			OP_BAND, OP_BOR, OP_BXOR, OP_SHL, OP_SHR, OP_EQ, OP_LT, OP_LE:
			if isK(b) || isK(c) {
				self.printf("\t; ")
				if isK(b) {
					self.printConstant(f, indexK(b))
				} else {
					self.printf("-")
				}
				self.printf(" ")
				if isK(c) {
					self.printConstant(f, indexK(c))
				} else {
					self.printf("-")
				}
			}
		case OP_JMP, OP_FORLOOP, OP_FORPREP, OP_TFORLOOP:
			self.printf("\t; to %d", sbx+pc+2)
		case OP_CLOSURE:
			// self.printf("\t; %p",VOID(f->p[bx]));
		case OP_SETLIST:
			if c == 0 {
				pc += 1
				self.printf("\t; %d", f.Code[pc])
			} else {
				self.printf("\t; %d", c)
			}
		case OP_EXTRAARG:
			self.printf("\t; ")
			self.printConstant(f, ax)
		}

		self.printf("\n")
	}
}

func (self *printer) printConstant(f *Prototype, i int) {
	k := f.Constants[i]
	switch x := k.(type) {
	case nil:
		self.printf("nil")
	case bool:
		self.printf("%t", x)
	case float64:
		self.printf("%.14g", x) // todo
	case int64:
		self.printf("%d", x) // todo
	case string:
		self.printf("%q", x) // todo
	default: /* cannot happen */
		self.printf("?")
	}
}

func (self *printer) printDebug(f *Prototype) {
	self.printf("constants (%d):\n", len(f.Constants))
	for i, _ := range f.Constants {
		self.printf("\t%d\t", i+1)
		self.printConstant(f, i)
		self.printf("\n")
	}
	self.printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		self.printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}
	self.printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		self.printf("\t%d\t%s\t%d\t%d\n",
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
