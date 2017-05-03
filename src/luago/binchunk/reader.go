package binchunk

import "encoding/binary"
import "math"
import "luago/lua"

type reader struct {
	data []byte
}

func (self *reader) readByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

func (self *reader) readBytes(n uint) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return i
}

func (self *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	size := uint(self.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(self.readUint64()) // size_t
	}
	bytes := self.readBytes(size - 1)
	return string(bytes) // todo
}

func (self *reader) readHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk")
	}
	if self.readByte() != LUA_VERSION {
		panic("version mismatch!")
	}
	if self.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if self.readByte() != CINT_SIZE {
		panic("unsupported int size!")
	}
	if self.readByte() != CSZIET_SIZE {
		panic("unsupported size_t size!")
	}
	if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if self.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (self *reader) readProto(parentSource string) *FuncProto {
	source := self.readString()
	if source == "" {
		source = parentSource
	}
	return &FuncProto{
		Source:          source,
		LineDefined:     self.readUint32(),
		LastLineDefined: self.readUint32(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readCode(),
		Constants:       self.readConstants(),
		Upvalues:        self.readUpvalues(),
		Protos:          self.readProtos(source),
		LineInfo:        self.readLineInfo(),
		LocVars:         self.readLocVars(),
		UpvalueNames:    self.readUpvalueNames(),
	}
}

func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}

func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readUint32())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

func (self *reader) readConstant() interface{} {
	switch lua.LuaType(self.readByte()) {
	case lua.LUA_TNIL:
		return nil
	case lua.LUA_TBOOLEAN:
		return self.readByte() != 0
	case lua.LUA_TNUMINT:
		return self.readLuaInteger()
	case lua.LUA_TNUMFLT:
		return self.readLuaNumber()
	case lua.LUA_TSHRSTR, lua.LUA_TLNGSTR:
		return self.readString()
	default:
		panic("corrupted!") // todo
	}
}

func (self *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, self.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}
	return upvalues
}

func (self *reader) readProtos(parentSource string) []*FuncProto {
	protos := make([]*FuncProto, self.readUint32())
	for i := range protos {
		protos[i] = self.readProto(parentSource)
	}
	return protos
}

func (self *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, self.readUint32())
	for i := range lineInfo {
		lineInfo[i] = self.readUint32()
	}
	return lineInfo
}

func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPc: self.readUint32(),
			EndPc:   self.readUint32(),
		}
	}
	return locVars
}

func (self *reader) readUpvalueNames() []string {
	names := make([]string, self.readUint32())
	for i := range names {
		names[i] = self.readString()
	}
	return names
}
