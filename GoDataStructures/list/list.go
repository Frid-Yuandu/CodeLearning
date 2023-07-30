package list

type List[T comparable] interface {
	IsEmpty() bool
	Len() uint
	Clear()
	Clone() List[T]

	AddAtHead(T)
	AddAtTail(T)
	AddAt(index uint, value T)

	Find(value T) (index uint, found bool)
	Get(index uint) (value T, found bool)
	Head() (value T, found bool)
	Tail() (value T, found bool)

	Set(index uint, value T) bool

	Remove(uint) (value T, succeed bool)
	RemoveValue(T) (index uint, succeed bool)
}
