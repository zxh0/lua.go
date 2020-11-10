package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *reader) readBytes(n uint32) []byte {
	bytes := r.data[:n]
	r.data = r.data[n:]
	return bytes
}

func (r *reader) readVarUint32() uint32 {
	return uint32(r.readVarUint(math.MaxUint32))
}
func (r *reader) readVarUint(limit uint64) uint64 {
	var x uint64 = 0
	limit >>= 7
	for {
		b := r.readByte()
		if x >= limit {
			panic("integer overflow")
		}
		x = (x << 7) | uint64(b&0x7f)

		if (b & 0x80) != 0 {
			break
		}
	}
	return x
}

func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}
func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}
func (r *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

func (r *reader) readInstruction() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

func (r *reader) readString() string {
	size := r.readVarUint32()
	if size == 0 {
		return ""
	}
	bytes := r.readBytes(size - 1)
	return string(bytes) // TODO
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
		LineDefined:     r.readVarUint32(),
		LastLineDefined: r.readVarUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		AbsLineInfo:     r.readAbsLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readVarUint32())
	for i := range code {
		code[i] = r.readInstruction()
	}
	return code
}

func (r *reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readVarUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case LUA_VNIL:
		return nil
	case LUA_VFALSE:
		return false
	case LUA_VTRUE:
		return true
	case LUA_VNUMINT:
		return r.readLuaInteger()
	case LUA_VNUMFLT:
		return r.readLuaNumber()
	case LUA_TSHRSTR, LUA_TLNGSTR:
		return r.readString()
	default:
		panic("corrupted!") // TODO
	}
}

func (r *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readVarUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
			Kind:    r.readByte(),
		}
	}
	return upvalues
}

func (r *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, r.readVarUint32())
	for i := range protos {
		protos[i] = r.readProto(parentSource)
	}
	return protos
}

func (r *reader) readLineInfo() []byte {
	size := r.readVarUint32()
	return r.readBytes(size)
}

func (r *reader) readAbsLineInfo() []AbsLineInfo {
	absLineInfo := make([]AbsLineInfo, r.readVarUint32())
	for i := range absLineInfo {
		absLineInfo[i] = AbsLineInfo{
			PC:   r.readVarUint32(),
			Line: r.readVarUint32(),
		}
	}
	return absLineInfo
}

func (r *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readVarUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readVarUint32(),
			EndPC:   r.readVarUint32(),
		}
	}
	return locVars
}

func (r *reader) readUpvalueNames() []string {
	names := make([]string, r.readVarUint32())
	for i := range names {
		names[i] = r.readString()
	}
	return names
}
