package binchunk

const CINT_SIZE = 4
const CSZIET_SIZE = 8
const INSTRUCTION_SIZE = 4
const LUA_INTEGER_SIZE = 8
const LUA_NUMBER_SIZE = 8

const LUA_SIGNATURE = "\x1bLua"
const LUA_VERSION byte = 0x53
const LUAC_FORMAT byte = 0
const LUAC_DATA = "\x19\x93\r\n\x1a\n"
const LUAC_INT int64 = 0x5678
const LUAC_NUM float64 = 370.5

type binaryChunk struct {
	binaryChunkHeader
	sizeUpvalues byte // ?
	mainFunc     *Prototype
}

type binaryChunkHeader struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
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
	LineInfo        []uint32 // debug
	LocVars         []LocVar // debug
	UpvalueNames    []string // debug
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPc uint32
	EndPc   uint32
}

func IsBinaryChunk(data []byte) bool {
	return len(data) > 4 &&
		string(data[:4]) == LUA_SIGNATURE
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.readHeader()
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
