package collection

import (
	"errors"
)

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(elem string) {
	*s = append(*s, elem)
}

func (s *Stack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", errStackEmpty
	}

	index := len(*s) - 1
	elem := (*s)[index]
	*s = (*s)[:index]
	return elem, nil
}

func (s *Stack) Peek() (string, error) {
	if s.IsEmpty() {
		return "", errStackEmpty
	}

	index := len(*s) - 1
	elem := (*s)[index]
	return elem, nil
}

var (
	errStackEmpty = errors.New("stack is empty")
)
