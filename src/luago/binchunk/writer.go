package binchunk

import "encoding/binary"
import "math"
import . "luago/api"

type writer struct {
	buf []byte
	idx int
	cap int
}

func (self *writer) data() []byte {
	return self.buf[:self.idx]
}

func (self *writer) grow(n int) {
	if self.cap-self.idx < n {
		self.cap += self.cap/2 + n
		newBuf := make([]byte, self.cap)
		copy(newBuf, self.buf)
		self.buf = newBuf
	}
}

func (self *writer) writeByte(b byte) {
	self.grow(1)
	self.buf[self.idx] = b
	self.idx += 1
}

func (self *writer) writeBytes(s []byte) {
	self.grow(len(s))
	copy(self.buf[self.idx:], s)
	self.idx += len(s)
}

func (self *writer) writeUint32(i uint32) {
	self.grow(4)
	binary.LittleEndian.PutUint32(self.buf[self.idx:], i)
	self.idx += 4
}

func (self *writer) writeUint64(i uint64) {
	self.grow(8)
	binary.LittleEndian.PutUint64(self.buf[self.idx:], i)
	self.idx += 8
}

func (self *writer) writeLuaInteger(i int64) {
	self.writeUint64(uint64(i))
}

func (self *writer) writeLuaNumber(f float64) {
	self.writeUint64(math.Float64bits(f))
}

func (self *writer) writeString(s string) {
	size := len(s)
	if size == 0 {
		self.writeByte(0)
		return
	}

	size += 1
	if size < 0xFF {
		self.writeByte(byte(size))
	} else {
		self.writeByte(0xFF)
		self.writeUint64(uint64(size)) // size_t
	}

	self.writeBytes([]byte(s))
}

func (self *writer) writeHeader() {
	self.writeBytes([]byte(LUA_SIGNATURE))
	self.writeByte(LUAC_VERSION)
	self.writeByte(LUAC_FORMAT)
	self.writeBytes([]byte(LUAC_DATA))
	self.writeByte(CINT_SIZE)
	self.writeByte(CSZIET_SIZE)
	self.writeByte(INSTRUCTION_SIZE)
	self.writeByte(LUA_INTEGER_SIZE)
	self.writeByte(LUA_NUMBER_SIZE)
	self.writeLuaInteger(LUAC_INT)
	self.writeLuaNumber(LUAC_NUM)
}

func (self *writer) writeProto(proto *Prototype, parentSource string) {
	if proto.Source == parentSource {
		self.writeString("")
	} else {
		self.writeString(proto.Source)
	}
	self.writeUint32(proto.LineDefined)
	self.writeUint32(proto.LastLineDefined)
	self.writeByte(proto.NumParams)
	self.writeByte(proto.IsVararg)
	self.writeByte(proto.MaxStackSize)
	self.writeCode(proto.Code)
	self.writeConstants(proto.Constants)
	self.writeUpvalues(proto.Upvalues)
	self.writeProtos(proto.Protos, proto.Source)
	self.writeLineInfo(proto.LineInfo)
	self.writeLocVars(proto.LocVars)
	self.writeUpvalueNames(proto.UpvalueNames)
}

func (self *writer) writeCode(code []uint32) {
	self.writeUint32(uint32(len(code)))
	for _, inst := range code {
		self.writeUint32(inst)
	}
}

func (self *writer) writeConstants(constants []interface{}) {
	self.writeUint32(uint32(len(constants)))
	for _, constant := range constants {
		self.writeConstant(constant)
	}
}

func (self *writer) writeConstant(constant interface{}) {
	switch x := constant.(type) {
	case nil:
		self.writeByte(byte(LUA_TNIL))
	case bool:
		self.writeByte(byte(LUA_TBOOLEAN))
		if x {
			self.writeByte(1)
		} else {
			self.writeByte(0)
		}
	case int64:
		self.writeByte(byte(LUA_TNUMINT))
		self.writeLuaInteger(x)
	case float64:
		self.writeByte(byte(LUA_TNUMFLT))
		self.writeLuaNumber(x)
	case string: // todo
		self.writeByte(byte(LUA_TSHRSTR))
		self.writeString(x)
	default:
		panic("unreachable!")
	}
}

func (self *writer) writeUpvalues(upvalues []Upvalue) {
	self.writeUint32(uint32(len(upvalues)))
	for _, upvalue := range upvalues {
		self.writeByte(upvalue.Instack)
		self.writeByte(upvalue.Idx)
	}
}

func (self *writer) writeProtos(protos []*Prototype, parentSource string) {
	self.writeUint32(uint32(len(protos)))
	for _, proto := range protos {
		self.writeProto(proto, parentSource)
	}
}

func (self *writer) writeLineInfo(lineInfo []uint32) {
	self.writeUint32(uint32(len(lineInfo)))
	for _, line := range lineInfo {
		self.writeUint32(line) // todo
	}
}

func (self *writer) writeLocVars(locVars []LocVar) {
	self.writeUint32(uint32(len(locVars)))
	for _, locVar := range locVars {
		self.writeString(locVar.VarName)
		self.writeUint32(locVar.StartPc)
		self.writeUint32(locVar.EndPc)
	}
}

func (self *writer) writeUpvalueNames(names []string) {
	self.writeUint32(uint32(len(names)))
	for _, name := range names {
		self.writeString(name)
	}
}
