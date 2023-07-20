package arrayStack

import (
	"fmt"
	"sync"
	"testing"
)

func TestSafeArrayStack_DataRace(t *testing.T) {
	var wg sync.WaitGroup
	s := NewSafeArrayStack[int]()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			s.Push(i)
		}(i)
	}

	for i := 0; i < 999; i++ {
		go func() {
			defer wg.Done()
			s.Pop()
		}()
	}

	go func() {
		defer wg.Done()
		s.Clear()
	}()

	wg.Wait()
}

func TestSafeArrayStack_CloneReference(t *testing.T) {
	s := NewSafeArrayStack[int]()
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}
	s2 := s.Clone()
	if &s2.data == & s.data {
		t.Error("these stack have the same reference")
	}else {
		fmt.Printf("address of s: %p, address of s2: %p\n",
			&s.data, &s2.data)
	}
}

func BenchmarkSafeArrayStack_Push(b *testing.B) {
	s := NewSafeArrayStack[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
}

func BenchmarkSafeArrayStack_Pop(b *testing.B) {
	s := NewSafeArrayStack[int]()
	for i := 0; i < 1_000_000; i++ {
		s.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Pop()
	}
}

func BenchmarkSafeArrayStack_Clone(b *testing.B) {
	s := NewSafeArrayStack[int]()
	for i := 0; i < 1_000_000; i++ {
		s.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Clone()
	}
}

func BenchmarkSafeArrayStack_Clear(b *testing.B) {
	s := NewSafeArrayStack[int]()
	for i := 0; i < 1_000_000; i++ {
		s.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_s := s.Clone()
		_s.Clear()
	}
}
