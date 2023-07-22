package arraylist

import (
	"fmt"
	"yo/list"
)

const oneByte = 8

type SingleArrayList[T comparable] struct {
	data []T
	len  uint
}

func (l *SingleArrayList[T]) Len() uint     { return l.len }
func (l *SingleArrayList[T]) IsEmpty() bool { return l.len == 0 }

func (l *SingleArrayList[T]) init() *SingleArrayList[T] {
	l.data = make([]T, 0, oneByte)
	l.len = 0
	return l
}

func (l *SingleArrayList[T]) lazyInit() {
	if l.data == nil {
		l.init()
	}
}

func (l *SingleArrayList[T]) Clear() {
	l.init()
}

func (l *SingleArrayList[T]) Clone() list.SingleList[T] {
	return &SingleArrayList[T]{
		data: l.data,
		len:  l.len,
	}
}

func (l *SingleArrayList[T]) AddAtHead(value T) {
	l.lazyInit()
	l.add(value, []T{}, l.data)
}

func (l *SingleArrayList[T]) AddAtTail(value T) {
	l.lazyInit()
	l.add(value, l.data)
}

func (l *SingleArrayList[T]) AddAt(index uint, value T) {
	l.lazyInit()
	if index > l.len {
		return
	} else if l.IsEmpty() || index == l.len {
		l.add(value, l.data)
	} else {
		l.add(value, l.data[:index-1], l.data[index-1:])
	}
}

// the first element of joint is the head, and the second element is the tail.
func (l *SingleArrayList[T]) add(value T, joint ...[]T) {
	result := append(joint[0], value)
	if len(joint) > 1 {
		result = append(result, joint[1]...)
	}
	l.data = result
	l.len++
}

func (l *SingleArrayList[T]) find(value T) (uint, bool) {
	for i := uint(0); i < l.len; i++ {
		if l.data[i] == value {
			return i, true
		}
	}
	return 0, false
}

func (l *SingleArrayList[T]) Find(value T) (uint, bool) {
	return l.find(value)
}

func (l *SingleArrayList[T]) Get(index uint) (value T, found bool) {
	if l.IsEmpty() || index >= l.len {
		return
	}
	return l.data[index], true
}

func (l *SingleArrayList[T]) Head() (value T, found bool) {
	l.lazyInit()
	return l.Get(0)
}

func (l *SingleArrayList[T]) Tail() (value T, found bool) {
	l.lazyInit()
	return l.Get(l.len - 1)
}

func (l *SingleArrayList[T]) Set(index uint, value T) bool {
	if _, ok := l.Get(index); !ok {
		return false
	}
	l.data[index] = value
	return true
}

func (l *SingleArrayList[T]) Remove(index uint) (value T, found bool) {
	if value, found = l.Get(index); !found {
		return
	}
	l.remove(index)
	return
}

func (l *SingleArrayList[T]) RemoveByValue(value T) (index uint, found bool) {
	if index, found = l.find(value); !found {
		return
	}
	l.remove(index)
	return
}

func (l *SingleArrayList[T]) remove(index uint) {
	if index == l.len-1 {
		l.data = l.data[:index]
	} else {
		l.data = append(l.data[:index], l.data[index+1:]...)
	}
	l.len--
}

func (l *SingleArrayList[T]) Show() (showed string) {
	if !l.IsEmpty() {
		showed += fmt.Sprintf("%v", l.data)
		for i := uint(1); i < l.len; i++ {
			showed += fmt.Sprintf(" -> %v", l.data[i])
		}
	}
	return
}
