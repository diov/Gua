package state

func (s *luaState) GetTop() int {
	return s.stack.top
}

func (s *luaState) AbsIndex(idx int) int {
	return s.stack.absIndex(idx)
}

func (s *luaState) CheckStack(n int) bool {
	s.stack.check(n)
	return true
}

func (s *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		s.stack.pop()
	}
}

func (s *luaState) Copy(fromIdx, toIdx int) {
	v := s.stack.get(fromIdx)
	s.stack.set(toIdx, v)
}

func (s *luaState) PushValue(idx int) {
	v := s.stack.get(idx)
	s.stack.push(v)
}

func (s *luaState) Replace(idx int) {
	v := s.stack.pop()
	s.stack.set(idx, v)
}

func (s *luaState) Insert(idx int) {
	s.Rotate(idx, 1)
}

func (s *luaState) Remove(idx int) {
	s.Rotate(idx, -1)
	s.Pop(1)
}

func (s *luaState) Rotate(idx, n int) {
	t := s.stack.top - 1
	p := s.stack.absIndex(idx)
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	s.stack.reverse(p, m)
	s.stack.reverse(m+1, t)
	s.stack.reverse(p, t)
}

func (s *luaState) SetTop(idx int) {
	newTop := s.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}
	n := s.stack.top - newTop
	if n >= 0 {
		for i := 0; i < n; i++ {
			s.stack.pop()
		}
	} else {
		for i := 0; i > n; i-- {
			s.stack.push(nil)
		}
	}
}
