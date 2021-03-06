package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

// byte
func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

// c_int
func (r *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

// size_t
func (r *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

// Lua int
func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

// Lua number
func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

// string
func (r *reader) readString() string {
	size := uint(r.readByte()) // Short string
	if size == 0 {             // NULL
		return ""
	}
	if size == 0xFF { // Long string
		size = uint(r.readUint64())
	}
	bytes := r.readBytes(size - 1)
	return string(bytes)
}

func (r *reader) readBytes(n uint) []byte {
	bytes := r.data[:n]
	r.data = r.data[n:]
	return bytes
}

func (r *reader) checkHeader() {
	if LuaSignature != string(r.readBytes(4)) {
		panic("not a binary chunk!")
	}
	if LuacVersion != r.readByte() {
		panic("version mismatch!")
	}
	if LuacFormat != r.readByte() {
		panic("format mismatch!")
	}
	if LuacData != string(r.readBytes(6)) {
		panic("corrupted chunk!")
	}
	if CIntSize != r.readByte() {
		panic("int size mismatch!")
	}
	if CSizeTSize != r.readByte() {
		panic("size_t size mismatch!")
	}
	if InstructionSize != r.readByte() {
		panic("Instruction size mismatch!")
	}
	if LuaIntegerSize != r.readByte() {
		panic("lua_Integer size mismatch!")
	}
	if LuaNumberSize != r.readByte() {
		panic("lua_Number!")
	}
	if LuacInt != r.readLuaInteger() {
		panic("integer format mismatch!")
	}
	if LuacNum != r.readLuaNumber() {
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
	case TagNil:
		return nil
	case TagBoolean:
		return r.readByte() != 0
	case TagInteger:
		return r.readLuaInteger()
	case TagNumber:
		return r.readLuaNumber()
	case TagShortStr, TagLongStr:
		return r.readString()
	default:
		panic("constant tag corrupted!")
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
	upvalueNames := make([]string, r.readUint32())
	for i := range upvalueNames {
		upvalueNames[i] = r.readString()
	}
	return upvalueNames
}
