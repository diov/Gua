package state

func (s *luaState) PushNil()              { s.stack.push(nil) }
func (s *luaState) PushBoolean(b bool)    { s.stack.push(b) }
func (s *luaState) PushInteger(n int64)   { s.stack.push(n) }
func (s *luaState) PushNumber(n float64)  { s.stack.push(n) }
func (s *luaState) PushString(str string) { s.stack.push(str) }
