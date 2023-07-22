package list

type SingleList[T comparable] interface {
	IsEmpty() bool
	Len() uint
	Clear()
	Clone() SingleList[T]

	AddAt(uint, T)

	Find(T) (uint, bool)
	Get(uint) (T, bool)
	Head() (T, bool)
	Tail() (T, bool)

	Set(uint, T) bool

	Remove(uint) (T, bool)
}
