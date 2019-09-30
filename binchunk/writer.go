package binchunk

import (
	"encoding/binary"
	"math"

	. "github.com/zxh0/lua.go/api"
)

type writer struct {
	buf []byte
	idx int
	cap int
}

func (w *writer) data() []byte {
	return w.buf[:w.idx]
}

func (w *writer) grow(n int) {
	if w.cap-w.idx < n {
		w.cap += w.cap/2 + n
		newBuf := make([]byte, w.cap)
		copy(newBuf, w.buf)
		w.buf = newBuf
	}
}

func (w *writer) writeByte(b byte) {
	w.grow(1)
	w.buf[w.idx] = b
	w.idx += 1
}

func (w *writer) writeBytes(s []byte) {
	w.grow(len(s))
	copy(w.buf[w.idx:], s)
	w.idx += len(s)
}

func (w *writer) writeUint32(i uint32) {
	w.grow(4)
	binary.LittleEndian.PutUint32(w.buf[w.idx:], i)
	w.idx += 4
}

func (w *writer) writeUint64(i uint64) {
	w.grow(8)
	binary.LittleEndian.PutUint64(w.buf[w.idx:], i)
	w.idx += 8
}

func (w *writer) writeLuaInteger(i int64) {
	w.writeUint64(uint64(i))
}

func (w *writer) writeLuaNumber(f float64) {
	w.writeUint64(math.Float64bits(f))
}

func (w *writer) writeString(s string) {
	size := len(s)
	if size == 0 {
		w.writeByte(0)
		return
	}

	size += 1
	if size < 0xFF {
		w.writeByte(byte(size))
	} else {
		w.writeByte(0xFF)
		w.writeUint64(uint64(size)) // size_t
	}

	w.writeBytes([]byte(s))
}

func (w *writer) writeHeader() {
	w.writeBytes([]byte(LUA_SIGNATURE))
	w.writeByte(LUAC_VERSION)
	w.writeByte(LUAC_FORMAT)
	w.writeBytes([]byte(LUAC_DATA))
	w.writeByte(CINT_SIZE)
	w.writeByte(CSIZET_SIZE)
	w.writeByte(INSTRUCTION_SIZE)
	w.writeByte(LUA_INTEGER_SIZE)
	w.writeByte(LUA_NUMBER_SIZE)
	w.writeLuaInteger(LUAC_INT)
	w.writeLuaNumber(LUAC_NUM)
}

func (w *writer) writeProto(proto *Prototype, parentSource string) {
	if proto.Source == parentSource {
		w.writeString("")
	} else {
		w.writeString(proto.Source)
	}
	w.writeUint32(proto.LineDefined)
	w.writeUint32(proto.LastLineDefined)
	w.writeByte(proto.NumParams)
	w.writeByte(proto.IsVararg)
	w.writeByte(proto.MaxStackSize)
	w.writeCode(proto.Code)
	w.writeConstants(proto.Constants)
	w.writeUpvalues(proto.Upvalues)
	w.writeProtos(proto.Protos, proto.Source)
	w.writeLineInfo(proto.LineInfo)
	w.writeLocVars(proto.LocVars)
	w.writeUpvalueNames(proto.UpvalueNames)
}

func (w *writer) writeCode(code []uint32) {
	w.writeUint32(uint32(len(code)))
	for _, inst := range code {
		w.writeUint32(inst)
	}
}

func (w *writer) writeConstants(constants []interface{}) {
	w.writeUint32(uint32(len(constants)))
	for _, constant := range constants {
		w.writeConstant(constant)
	}
}

func (w *writer) writeConstant(constant interface{}) {
	switch x := constant.(type) {
	case nil:
		w.writeByte(byte(LUA_TNIL))
	case bool:
		w.writeByte(byte(LUA_TBOOLEAN))
		if x {
			w.writeByte(1)
		} else {
			w.writeByte(0)
		}
	case int64:
		w.writeByte(byte(LUA_TNUMINT))
		w.writeLuaInteger(x)
	case float64:
		w.writeByte(byte(LUA_TNUMFLT))
		w.writeLuaNumber(x)
	case string: // todo
		w.writeByte(byte(LUA_TSHRSTR))
		w.writeString(x)
	default:
		panic("unreachable!")
	}
}

func (w *writer) writeUpvalues(upvalues []Upvalue) {
	w.writeUint32(uint32(len(upvalues)))
	for _, upvalue := range upvalues {
		w.writeByte(upvalue.Instack)
		w.writeByte(upvalue.Idx)
	}
}

func (w *writer) writeProtos(protos []*Prototype, parentSource string) {
	w.writeUint32(uint32(len(protos)))
	for _, proto := range protos {
		w.writeProto(proto, parentSource)
	}
}

func (w *writer) writeLineInfo(lineInfo []uint32) {
	w.writeUint32(uint32(len(lineInfo)))
	for _, line := range lineInfo {
		w.writeUint32(line) // todo
	}
}

func (w *writer) writeLocVars(locVars []LocVar) {
	w.writeUint32(uint32(len(locVars)))
	for _, locVar := range locVars {
		w.writeString(locVar.VarName)
		w.writeUint32(locVar.StartPC)
		w.writeUint32(locVar.EndPC)
	}
}

func (w *writer) writeUpvalueNames(names []string) {
	w.writeUint32(uint32(len(names)))
	for _, name := range names {
		w.writeString(name)
	}
}
