package state

import (
	"Gua/api"
	"fmt"
)

type luaValue interface{}

func typeOf(v luaValue) LuaType {
	switch v.(type) {
	case nil:
		return api.LuaTNil
	case bool:
		return api.LuaTBoolean
	case int64, float64:
		return api.LuaTNumber
	case string:
		return api.LuaTString
	default:
		fmt.Print(v)
		panic("To be done!")
	}
}
