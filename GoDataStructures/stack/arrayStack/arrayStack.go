package arrayStack

type ArrayStack[T any] struct {
	data []T
	top  uint
}

func (s *ArrayStack[T]) Size() uint { return s.top }

func (s *ArrayStack[T]) Empty() bool { return s.top == 0 }

func (s *ArrayStack[T]) Push(value T) {
	s.lazyInit()
	s.data = append(s.data, value)
	s.top++
}

func (s *ArrayStack[T]) Top() (T, bool) {
	var zeroValue T
	if s.Empty() {
		return zeroValue, false
	}
	return s.data[s.top-1], true
}

func (s *ArrayStack[T]) Pop() (T, bool) {
	var zeroValue T
	if s.Empty() {
		return zeroValue, false
	}
	return s.popAndReturn(), true
}

func (s *ArrayStack[T]) popAndReturn() T {
	res := s.data[s.top-1]
	s.safePop()
	return res
}

func (s *ArrayStack[T]) safePop() {
	var zeroValue T
	s.data[s.top-1] = zeroValue // avoid memory leaks
	s.data = s.data[:s.top-1]
	s.top--
}

func (s *ArrayStack[T]) Clear() {
	s.init()
}

func (s *ArrayStack[T]) init() *ArrayStack[T] {
	s.data = make([]T, 0, 10)
	s.top = 0
	return s
}

func (s *ArrayStack[T]) lazyInit() {
	if s.data == nil && s.top == 0 {
		s.init()
	}
}

func (s *ArrayStack[T]) Clone() *ArrayStack[T] {
	return &ArrayStack[T]{data: s.data, top: s.top}
}

func NewArrayStack[T any]() *ArrayStack[T] { return new(ArrayStack[T]).init() }
