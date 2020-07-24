package main

import (
	"Gua/binchunk"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		file, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		prototype := binchunk.Undump(file)
		listProto(prototype)
	}
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
	for i, code := range p.Code {
		line := "-"
		if len(p.LineInfo) > 0 {
			line = fmt.Sprintf("%d", p.LineInfo[i])
		}
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", i+1, line, code)
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
