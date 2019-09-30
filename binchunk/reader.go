package binchunk

import (
	"encoding/binary"
	"math"

	. "github.com/zxh0/lua.go/api"
)

type reader struct {
	data []byte
}

func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *reader) readBytes(n uint) []byte {
	bytes := r.data[:n]
	r.data = r.data[n:]
	return bytes
}

func (r *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

func (r *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

func (r *reader) readString() string {
	size := uint(r.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(r.readUint64()) // size_t
	}
	bytes := r.readBytes(size - 1)
	return string(bytes) // todo
}

func (r *reader) checkHeader() {
	if string(r.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk")
	}
	if r.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	if r.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(r.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if r.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}
	if r.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	}
	if r.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if r.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if r.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if r.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if r.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (r *reader) readProto(parentSource string) *Prototype {
	source := r.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

func (r *reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case LUA_TNIL:
		return nil
	case LUA_TBOOLEAN:
		return r.readByte() != 0
	case LUA_TNUMINT:
		return r.readLuaInteger()
	case LUA_TNUMFLT:
		return r.readLuaNumber()
	case LUA_TSHRSTR, LUA_TLNGSTR:
		return r.readString()
	default:
		panic("corrupted!") // todo
	}
}

func (r *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upvalues
}

func (r *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(parentSource)
	}
	return protos
}

func (r *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, r.readUint32())
	for i := range lineInfo {
		lineInfo[i] = r.readUint32()
	}
	return lineInfo
}

func (r *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return locVars
}

func (r *reader) readUpvalueNames() []string {
	names := make([]string, r.readUint32())
	for i := range names {
		names[i] = r.readString()
	}
	return names
}
