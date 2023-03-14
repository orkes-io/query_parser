package stack

type Stack[T any] struct {
	storage []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		storage: []T{},
	}
}

func (s *Stack[T]) Push(v T) {
	s.storage = append(s.storage, v)
}

func (s *Stack[T]) Pop() T {
	lastElement := s.storage[len(s.storage)-1]
	s.storage = s.storage[:len(s.storage)-1]
	return lastElement
}

func (s *Stack[T]) Top() T {
	return s.storage[len(s.storage)-1]
}

func (s *Stack[T]) Size() int {
	return len(s.storage)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.storage) == 0
}
