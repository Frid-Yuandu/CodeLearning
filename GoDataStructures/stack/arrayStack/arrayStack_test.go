package arrayStack

import (
	"reflect"
	"testing"
)

func checkStackSize[T any](t *testing.T, s *ArrayStack[T], expected uint) bool {
	if size := s.Size(); size != expected {
		t.Errorf("s.Size() = %d, expected %d", size, expected)
		return false
	}
	return true
}

func checkStackTop[T any](t *testing.T, s *ArrayStack[T], expectedValue T, expectedOk bool) {
	value, ok := s.Top()
	if !reflect.DeepEqual(value, expectedValue) || ok != expectedOk {
		t.Errorf("s.Top() = %v, %v, expected %v, %v", value, ok, expectedValue, expectedOk)
	}
}

func checkStackPop[T any](t *testing.T, s *ArrayStack[T], expectedValue T, expectedOk bool) {
	value, ok := s.Pop()
	if !reflect.DeepEqual(value, expectedValue) || ok != expectedOk {
		t.Errorf("s.Pop() = %v, %v, expected %v, %v", value, ok, expectedValue, expectedOk)
	}
}

func checkStack[T any](t *testing.T, s *ArrayStack[T], expected []T) {
	size := s.Size()
	if !checkStackSize(t, s, uint(len(expected))) {
		return
	}
	for i := uint(0); i < size; i++ {
		if !reflect.DeepEqual(s.data[i], expected[i]) {
			t.Errorf("s.data[%d] = %v, expected %v", i, s.data[i], expected[i])
		}
	}
}

func checkStackClear[T any](t *testing.T, s *ArrayStack[T]) {
	beforeAddress := &s.data
	s.Clear()
	afterAddress := &s.data
	if beforeAddress != afterAddress {
		t.Errorf("beforeAddress = %p, afterAddress = %p", beforeAddress, afterAddress)
	}
	checkStack(t, s, []T{})
}

func TestArrayStack_Size(t *testing.T) {
	s := NewArrayStack[int]()
	checkStackSize(t, s, 0)

	s.Push(1)
	checkStackSize(t, s, 1)
	s.Push(2)
	checkStackSize(t, s, 2)

	s.Clear()
	checkStackSize(t, s, 0)
}

func TestArrayStack_Top(t *testing.T) {
	s := NewArrayStack[int]()
	checkStackTop(t, s, 0, false)

	s.Push(3)
	checkStackTop(t, s, 3, true)
	s.Push(4)
	checkStackTop(t, s, 4, true)
	s.Pop()
	checkStackTop(t, s, 3, true)
	s.Pop()
	checkStackTop(t, s, 0, false)
}

func TestArrayStack_Pop(t *testing.T) {
	s := NewArrayStack[int]()
	checkStackPop(t, s, 0, false)

	s.Push(5)
	s.Push(6)
	checkStackPop(t, s, 6, true)
	checkStackPop(t, s, 5, true)
	checkStackPop(t, s, 0, false)
}

func TestArrayStack(t *testing.T) {
	s := NewArrayStack[int]()
	checkStackTop(t, s, 0, false)
	checkStackPop(t, s, 0, false)

	// single element stack
	s.Push(1)
	checkStack(t, s, []int{1})

	// bigger stack
	s.Push(2)
	s.Push(3)
	checkStack(t, s, []int{1, 2, 3})

	// clear stack
	s.Push(4)
	s.Push(5)
	checkStack(t, s, []int{1, 2, 3, 4, 5})
	checkStackClear(t, s)
	checkStack(t, s, []int{})
}
