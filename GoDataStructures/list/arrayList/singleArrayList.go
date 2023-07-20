package arraylist

const oneByte = 8

type List[T any] struct {
	data []T
	len  uint
}

func (l *List[T]) Size() uint  { return l.len }
func (l *List[T]) Empty() bool { return l.len == 0 }

func (l *List[T]) init() *List[T] {
	l.data = make([]T, 0, oneByte)
	l.len = 0
	return l
}

func (l *List[T]) lazyInit() {
	if l.data == nil {
		l.init()
	}
}

func (l *List[T]) AddAtHead(value T) {
	l.lazyInit()
	l.data = append([]T{value}, l.data...)
	l.len++
}

func (l *List[T]) AddAtTail(value T) {
	l.lazyInit()
	l.data = append(l.data, value)
	l.len++
}

func (l *List[T]) AddAtIndex(index uint, value T) {
	// TODO: implement me
	panic("implement me")
}

func (l *List[T]) find(value T) uint {
	// TODO: implement me
	panic("implement me")
}

func (l *List[T]) RemoveByIndex(index uint) {
	// TODO: implement me
	panic("implement me")
}

func (l *List[T]) RemoveByValue(index uint) {
	// TODO: implement me
	panic("implement me")
}

func (l *List[T]) Head() (value T, found bool) {
	if l.Empty() {
		l.init()
		return value, false
	}
	return l.data[0], true
}

func (l *List[T]) Tail() (value T, found bool) {
	if l.Empty() {
		l.init()
		return value, false
	}
	return l.data[l.len-1], true
}

func (l *List[T]) Show() string {
	// TODO: implement me
	panic("implement me")
}
