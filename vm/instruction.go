package vm

const MaxArgBx = 1<<18 - 1      // 2^18 - 1 = 262143
const MaxArgSbx = MaxArgBx >> 1 // 262143 / 2 = 131071

type Instruction uint32

func (i Instruction) Opcode() int {
	return int(i & 0x3F)
}

func (i Instruction) ABC() (a, b, c int) {
	a = int(i >> 6 & 0xFF)
	b = int(i >> 14 & 0x1FF)
	c = int(i >> 23 & 0x1FF)
	return
}

func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 6 & 0xFF)
	bx = int(i >> 14)
	return
}

func (i Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - MaxArgSbx
}

func (i Instruction) Ax() int {
	return int(i >> 6)
}

func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

func (i Instruction) BMode() byte {
	return opcodes[i.Opcode()].argBMode
}

func (i Instruction) CMode() byte {
	return opcodes[i.Opcode()].argCMode
}
