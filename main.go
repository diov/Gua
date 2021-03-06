package main

import (
	"Gua/api"
	"Gua/binchunk"
	"Gua/state"
	"Gua/vm"
	"fmt"
)

func main() {
	//if len(os.Args) > 1 {
	//	file, err := ioutil.ReadFile(os.Args[1])
	//	if err != nil {
	//		panic(err)
	//	}
	//	prototype := binchunk.Undump(file)
	//	listProto(prototype)
	//}

	s := state.New()
	s.PushBoolean(true)
	printStack(s)
	s.PushInteger(10)
	printStack(s)
	s.PushNil()
	printStack(s)
	s.PushString("Hello")
	printStack(s)
	s.PushValue(-4)
	printStack(s)
	s.Replace(3)
	printStack(s)
	s.SetTop(6)
	printStack(s)
	s.Remove(-3)
	printStack(s)
	s.SetTop(-5)
	printStack(s)
}

func listProto(p *binchunk.Prototype) {
	printHeader(p)
	printCode(p)
	printDetail(p)
	for _, proto := range p.Protos {
		listProto(proto)
	}
}

func printHeader(p *binchunk.Prototype) {
	funcType := "main"
	if p.LineDefined > 0 {
		funcType = "function"
	}
	varargFlag := ""
	if p.IsVararg > 0 {
		varargFlag = "+"
	}
	fmt.Printf("\n%s <%s:%d, %d> (%d instructions)\n",
		funcType, p.Source, p.LineDefined, p.LastLineDefined, len(p.Code))
	fmt.Printf("%d%s params, %d slots, %d upvalue, ",
		p.NumParams, varargFlag, p.MaxStackSize, len(p.Upvalues))
	fmt.Printf("%d locals, %d constant, %d function\n", len(p.LocVars), len(p.Constants), len(p.Protos))
}

func printCode(p *binchunk.Prototype) {
	for idx, code := range p.Code {
		line := "-"
		if len(p.LineInfo) > 0 {
			line = fmt.Sprintf("%d", p.LineInfo[idx])
		}
		i := vm.Instruction(code)
		fmt.Printf("\t%d\t[%s]\t%s \t", idx+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func printOperands(i vm.Instruction) {
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != vm.OpArgN {
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != vm.OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case vm.IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.BMode() == vm.OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == vm.OpArgU {
			fmt.Printf(" %d", bx)
		}
	case vm.IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case vm.IAx:
		ax := i.Ax()
		fmt.Printf("%d", ax)
	}
}

func printDetail(p *binchunk.Prototype) {
	fmt.Printf("Constants (%d):\n", len(p.Constants))
	for i, constant := range p.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(constant))
	}

	fmt.Printf("locals (%d):\n", len(p.LocVars))
	for i, locVar := range p.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(p.Upvalues))
	for i, upvalue := range p.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalueName(p, i), upvalue.Instack, upvalue.Idx)
	}
}

func constantToString(c interface{}) string {
	switch c.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", c)
	case float64:
		return fmt.Sprintf("%g", c)
	case int64:
		return fmt.Sprintf("%d", c)
	case string:
		return fmt.Sprintf("%q", c)
	default:
		return "?"
	}
}

func upvalueName(p *binchunk.Prototype, i int) string {
	if len(p.UpvalueNames) > 0 {
		return p.UpvalueNames[i]
	}
	return "-"
}

func printStack(s state.LuaState) {
	top := s.GetTop()
	for i := 1; i <= top; i++ {
		t := s.Type(i)
		switch t {
		case api.LuaTBoolean:
			fmt.Printf("[%t]", s.ToBoolean(i))
		case api.LuaTNumber:
			fmt.Printf("[%g]", s.ToNumber(i))
		case api.LuaTString:
			fmt.Printf("[%q]", s.ToString(i))
		default:
			fmt.Printf("[%s]", s.TypeName(t))
		}
	}
	fmt.Println()
}
