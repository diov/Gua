package state

import (
	"Gua/api"
	"fmt"
)

func (s *luaState) TypeName(lt LuaType) string {
	switch lt {
	case api.LuaTNone:
		return "no value"
	case api.LuaTNil:
		return "nil"
	case api.LuaTBoolean:
		return "boolean"
	case api.LuaTNumber:
		return "number"
	case api.LuaTString:
		return "string"
	case api.LuaTTable:
		return "table"
	case api.LuaTFunction:
		return "function"
	case api.LuaTLightUserData:
		return "light userdata"
	case api.LuaTThread:
		return "thread"
	default:
		return "userdata"
	}
}

func (s *luaState) Type(idx int) LuaType {
	if s.stack.isValid(idx) {
		v := s.stack.get(idx)
		return typeOf(v)
	}
	return api.LuaTNone
}

func (s *luaState) IsNone(idx int) bool {
	return s.Type(idx) == api.LuaTNone
}

func (s *luaState) IsNil(idx int) bool {
	return s.Type(idx) == api.LuaTNil
}

func (s *luaState) IsNoneOrNil(idx int) bool {
	return s.Type(idx) <= api.LuaTNil
}

func (s *luaState) IsBoolean(idx int) bool {
	return s.Type(idx) == api.LuaTBoolean
}

func (s *luaState) IsInteger(idx int) bool {
	v := s.stack.get(idx)
	_, ok := v.(int64)
	return ok
}

func (s *luaState) IsNumber(idx int) bool {
	_, ok := s.ToNumberX(idx)
	return ok
}

func (s *luaState) IsString(idx int) bool {
	t := s.Type(idx)
	return t == api.LuaTString || t == api.LuaTNumber
}

func (s *luaState) ToBoolean(idx int) bool {
	v := s.stack.get(idx)
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	default:
		return true
	}
}

func (s *luaState) ToInteger(idx int) int64 {
	i, _ := s.ToIntegerX(idx)
	return i
}

func (s *luaState) ToIntegerX(idx int) (int64, bool) {
	v := s.stack.get(idx)
	i, ok := v.(int64)
	return i, ok
}

func (s *luaState) ToNumber(idx int) float64 {
	n, _ := s.ToNumberX(idx)
	return n
}

func (s *luaState) ToNumberX(idx int) (float64, bool) {
	v := s.stack.get(idx)
	switch t := v.(type) {
	case float64:
		return t, true
	case int64:
		return float64(t), true
	default:
		return 0, false
	}
}
func (s *luaState) ToString(idx int) string {
	str, _ := s.ToStringX(idx)
	return str
}

func (s *luaState) ToStringX(idx int) (string, bool) {
	v := s.stack.get(idx)
	switch t := v.(type) {
	case string:
		return t, true
	case int64, float64:
		str := fmt.Sprintf("%v", t)
		s.stack.set(idx, str)
		return str, true
	default:
		return "", false
	}
}
