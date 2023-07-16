package arrayStack

import (
	"sync"
)

type SafeArrayStack[T any] struct {
	data []T
	top  uint64
	lock sync.RWMutex
}

func (s *SafeArrayStack[T]) Empty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.top == 0
}

func (s *SafeArrayStack[T]) Size() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.top
}

func (s *SafeArrayStack[T]) Push(value T) {
	s.lock.Lock()
	s.lazyInit()
	s.data = append(s.data, value)
	s.top++
	s.lock.Unlock()
}

func (s *SafeArrayStack[T]) Top() (value T, found bool) {
	if !s.Empty() {
		s.lock.RLock()
		value, found = s.data[s.top-1], true
		s.lock.RUnlock()
	}
	return
}

func (s *SafeArrayStack[T]) safePop() {
	var zeroValue T
	s.data[s.top-1] = zeroValue
	s.data = s.data[:s.top-1]
	s.top--
}

func (s *SafeArrayStack[T]) Pop() (value T, found bool) {
	if !s.Empty() {
		s.lock.Lock()
		value, found = s.data[s.top-1], true
		s.safePop()
		s.lock.Unlock()
	} else {
		s.lazyInit()
	}
	return
}

func (s *SafeArrayStack[T]) Clear() {
	s.lock.Lock()
	s.init()
	s.lock.Unlock()
}

func (s *SafeArrayStack[T]) Clone() *SafeArrayStack[T] {
	return &SafeArrayStack[T]{data: s.data, top: s.top}
}

func (s *SafeArrayStack[T]) init() {
	s.data = make([]T, 0)
	s.top = 0
}

func (s *SafeArrayStack[T]) lazyInit() {
	if s.data == nil && s.top == 0 {
		s.init()
	}
}

func NewSafeArrayStack[T any]() *SafeArrayStack[T] { return new(SafeArrayStack[T]) }
