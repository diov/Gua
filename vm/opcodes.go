package vm

const (
	IABC  = iota // B:9 C:9 A:8 Opcode: 6
	IABx         // Bx: 18 A: 8 Opcode: 6
	IAsBx        // sBx: 18 A: 8 Opcode: 6
	IAx          // Ax: 26 Opcode: 6
)

const (
	OpMove = iota
	OpLoadK
	OpLoadKX
	OpLoadBool
	OpLoadNil
	OpGetUpval
	OpGetTabup
	OpSetTabup
	OpSetUpval
	OpSetTable
	OpNewTable
	OpSelf
	OpAdd
	OpSub
	OpMul
	OpMod
	OpPow
	OpDiv
	OpIDiv
	OpBand
	OpBOr
	OpBXor
	OpShl
	OpShr
	OpUnm
	OpBNot
	OpNot
	OpLen
	OpConcat
	OpJmp
	OpEq
	OpLt
	OpLe
	OpTest
	OpTestSet
	OpCall
	OpTailCall
	OpReturn
	OpForLoop
	OpForPrep
	OpTForCall
	OpTForLoop
	OpSetList
	OpClosure
	OpVararg
	OpExtraArg
)

const (
	OpArgN = iota // is not used
	OpArgU        // is used
	OpArgR        // a register or a jump offset
	OpArgK        // constant or register/constant
)

type opcode struct {
	testFlag byte // whether operator is a test
	setAFlag byte // instruction set register A
	argBMode byte // B arg mode
	argCMode byte // C arg mode
	opMode   byte // op mode
	name     string
}

var opcodes = []opcode{
	{0, 1, OpArgR, OpArgN, IABC /* */, "MOVE    "}, // R(A) := R(B)
	{0, 1, OpArgK, OpArgN, IABx /* */, "LOADK   "}, // R(A) := Kst(Bx)
	{0, 1, OpArgN, OpArgN, IABx /* */, "LOADKX  "}, // R(A) := Kst(extra arg)
	{0, 1, OpArgU, OpArgU, IABC /* */, "LOADBOOL"}, // R(A) := (bool)B; if (C) pc++
	{0, 1, OpArgU, OpArgN, IABC /* */, "LOADNIL "}, // R(A), R(A+1), ..., R(A+B) := nil
	{0, 1, OpArgU, OpArgN, IABC /* */, "GETUPVAL"}, // R(A) := UpValue[B]
	{0, 1, OpArgU, OpArgK, IABC /* */, "GETTABUP"}, // R(A) := UpValue[B][RK(C)]
	{0, 1, OpArgR, OpArgK, IABC /* */, "GETTABLE"}, // R(A) := R(B)[RK(C)]
	{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABUP"}, // UpValue[A][RK(B)] := RK(C)
	{0, 0, OpArgU, OpArgN, IABC /* */, "SETUPVAL"}, // UpValue[B] := R(A)
	{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABLE"}, // R(A)[RK(B)] := RK(C)
	{0, 1, OpArgU, OpArgU, IABC /* */, "NEWTABLE"}, // R(A) := {} (size = B,C)
	{0, 1, OpArgR, OpArgK, IABC /* */, "SELF    "}, // R(A+1) := R(B); R(A) := R(B)[RK(C)]
	{0, 1, OpArgK, OpArgK, IABC /* */, "ADD     "}, // R(A) := RK(B) + RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "SUB     "}, // R(A) := RK(B) - RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "MUL     "}, // R(A) := RK(B) * RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "MOD     "}, // R(A) := RK(B) % RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "POW     "}, // R(A) := RK(B) ^ RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "DIV     "}, // R(A) := RK(B) / RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "IDIV    "}, // R(A) := RK(B) // RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "BAND    "}, // R(A) := RK(B) & RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "BOR     "}, // R(A) := RK(B) | RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "BXOR    "}, // R(A) := RK(B) ~ RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "SHL     "}, // R(A) := RK(B) << RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "SHR     "}, // R(A) := RK(B) >> RK(C)
	{0, 1, OpArgR, OpArgN, IABC /* */, "UNM     "}, // R(A) := -R(B)
	{0, 1, OpArgR, OpArgN, IABC /* */, "BNOT    "}, // R(A) := ~R(B)
	{0, 1, OpArgR, OpArgN, IABC /* */, "NOT     "}, // R(A) := not R(B)
	{0, 1, OpArgR, OpArgN, IABC /* */, "LEN     "}, // R(A) := length of R(B)
	{0, 1, OpArgR, OpArgR, IABC /* */, "CONCAT  "}, // R(A) := R(B).. ... ..R(C)
	{0, 0, OpArgR, OpArgN, IAsBx /**/, "JMP     "}, // pc+=sBx; if (A) close all upvalues >= R(A - 1)
	{1, 0, OpArgK, OpArgK, IABC /* */, "EQ      "}, // if ((RK(B) == RK(C)) ~= A) then pc++
	{1, 0, OpArgK, OpArgK, IABC /* */, "LT      "}, // if ((RK(B) <  RK(C)) ~= A) then pc++
	{1, 0, OpArgK, OpArgK, IABC /* */, "LE      "}, // if ((RK(B) <= RK(C)) ~= A) then pc++
	{1, 0, OpArgN, OpArgU, IABC /* */, "TEST    "}, // if not (R(A) <=> C) then pc++
	{1, 1, OpArgR, OpArgU, IABC /* */, "TESTSET "}, // if (R(B) <=> C) then R(A) := R(B) else pc++
	{0, 1, OpArgU, OpArgU, IABC /* */, "CALL    "}, // R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
	{0, 1, OpArgU, OpArgU, IABC /* */, "TAILCALL"}, // return R(A)(R(A+1), ... ,R(A+B-1))
	{0, 0, OpArgU, OpArgN, IABC /* */, "RETURN  "}, // return R(A), ... ,R(A+B-2)
	{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORLOOP "}, // R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }
	{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORPREP "}, // R(A)-=R(A+2); pc+=sBx
	{0, 0, OpArgN, OpArgU, IABC /* */, "TFORCALL"}, // R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
	{0, 1, OpArgR, OpArgN, IAsBx /**/, "TFORLOOP"}, // if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }
	{0, 0, OpArgU, OpArgU, IABC /* */, "SETLIST "}, // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
	{0, 1, OpArgU, OpArgN, IABx /* */, "CLOSURE "}, // R(A) := closure(KPROTO[Bx])
	{0, 1, OpArgU, OpArgN, IABC /* */, "VARARG  "}, // R(A), R(A+1), ..., R(A+B-2) = vararg
	{0, 0, OpArgU, OpArgU, IAx /*  */, "EXTRAARG"}, // extra (larger) argument for previous opcode
}
