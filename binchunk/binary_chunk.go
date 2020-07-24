package binchunk

type binaryChunk struct {
	header
	sizeUpvalues byte
	mainFunc     *Prototype
}

const (
	LuaSignature    = "\x1bLua"
	LuacVersion     = 0x53
	LuacFormat      = 0
	LuacData        = "\x19\x93\r\n\x1a\n"
	CIntSize        = 4
	CSizeTSize      = 8
	InstructionSize = 4
	LuaIntegerSize  = 8
	LuaNumberSize   = 8
	LuacInt         = 0x5678
	LuacNum         = 370.5
)

type header struct {
	signature       [4]byte // ESC、L、U、A 的 ASCII码 "\x1bLua"
	version         byte    // Major Version * 16 + Minor Version
	format          byte    // Format code
	luacData        [6]byte // \x19\x93\r\n\x1a\n
	cintSize        byte    // c_int size
	sizetSize       byte    // size_t size
	instructionSize byte    // instruction size
	luaIntegerSize  byte    // Lua Integer size
	luaNumberSize   byte    // Lua Number(float) size
	luacInt         int64   // Endianness
	luacNum         float64 // Float format
}

const (
	TagNil      = 0x00
	TagBoolean  = 0x01
	TagInteger  = 0x13
	TagNumber   = 0x03
	TagShortStr = 0x04
	TagLongStr  = 0x14
)

type Prototype struct {
	Source          string        // "@": file; "=": stream; "": literal
	LineDefined     uint32        // Start line number
	LastLineDefined uint32        // End line number
	NumParams       byte          // Count of parameters
	IsVararg        byte          // Enable vararg
	MaxStackSize    byte          // Count of registers
	Code            []uint32      // Instruction
	Constants       []interface{} // Constant with tag
	Upvalues        []Upvalue     // Upvalues
	Protos          []*Prototype  // Child function prototype
	LineInfo        []uint32      // Instruction line number
	LocVars         []LocVar      // Local vars
	UpvalueNames    []string      // Upvalue names
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        // header
	reader.readByte()           // sizeUpvalues
	return reader.readProto("") // mainFunc
}
