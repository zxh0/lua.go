package binchunk

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x54
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	LUA_VNIL    = 0x00 // nil
	LUA_VFALSE  = 0x01 // false
	LUA_VTRUE   = 0x11 // true
	LUA_VNUMINT = 0x03 // integer numbers
	LUA_VNUMFLT = 0x13 // float numbers
	LUA_TSHRSTR = 0x04 // short strings
	LUA_TLNGSTR = 0x14 // long strings
)

type BinaryChunk struct {
	Header
	sizeUpvalues byte // ?
	mainFunc     *Prototype
}

type Header struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64
	luacNum         float64
}

// function prototype
type Prototype struct {
	Source          string // debug
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []byte        // debug
	AbsLineInfo     []AbsLineInfo // debug
	LocVars         []LocVar      // debug
	UpvalueNames    []string      // debug
}

type Upvalue struct {
	Instack byte
	Idx     byte
	Kind    byte
}

type AbsLineInfo struct {
	PC   uint32
	Line uint32
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

func IsBinaryChunk(data []byte) bool {
	return len(data) > 4 &&
		string(data[:4]) == LUA_SIGNATURE
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte() // size_upvalues
	return reader.readProto("")
}

func Dump(proto *Prototype) []byte {
	writer := &writer{}
	writer.writeHeader()
	writer.writeByte(byte(len(proto.Upvalues)))
	writer.writeProto(proto, "")
	return writer.data()
}

func List(proto *Prototype, full bool) string {
	printer := &printer{make([]string, 0, 64)}
	return printer.printFunc(proto, full)
}

func StripDebug(proto *Prototype) {
	proto.Source = ""
	proto.LineInfo = nil
	proto.LocVars = nil
	proto.UpvalueNames = nil
	for _, p := range proto.Protos {
		StripDebug(p)
	}
}
