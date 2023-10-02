package jul

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Stack struct{ cells []any }

func NewStack(capacity int) *Stack {
	if capacity <= 0 {
		capacity = 4096
	}
	return &Stack{cells: make([]any, 0, capacity)}
}

var (
	ErrStackOverflow  = errors.New("stack overflow")
	ErrStackUnderflow = errors.New("stack underflow")
)

type (
	CellBoolean   bool
	CellInteger   int
	CellFloat     float64
	CellText      string
	CellQuotation string
)

func (s *Stack) Push(c any) error {
	if len(s.cells) == cap(s.cells) {
		return ErrStackOverflow
	}
	switch v := c.(type) {
	default:
		return newInvalidTypeError(v)
	case
		CellBoolean,
		CellInteger,
		CellFloat,
		CellText,
		CellQuotation:
	}
	s.cells = append(s.cells, c)
	return nil
}

func (s *Stack) Peek(i int) (any, error) {
	if len(s.cells)-i == 0 {
		return nil, ErrStackUnderflow
	}
	return s.cells[len(s.cells)-1-i], nil
}

func (s *Stack) Pop() (any, error) {
	c, err := s.Peek(0)
	if err != nil {
		return nil, err
	}
	s.cells = s.cells[:len(s.cells)-1]
	return c, nil
}

func newInvalidTypeError(v any) error {
	return fmt.Errorf("invalid type %T", v)
}

func newTypeMismatchError(a, b any) error {
	return fmt.Errorf("type mismatch between %T A and %T B", a, b)
}

func (s *Stack) Drop() error {
	_, err := s.Pop()
	return err
}

func (s *Stack) Pick() error {
	cell1, err := s.Pop()
	if err != nil {
		return err
	}
	i, ok := cell1.(CellInteger)
	if !ok {
		return newInvalidTypeError(cell1)
	}
	cellN, err := s.Peek(int(i))
	if err != nil {
		return err
	}
	return s.Push(cellN)
}

func (s *Stack) Swap() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	_ = s.Push(cellB)
	_ = s.Push(cellA)
	return nil
}

func (s *Stack) Rot() error {
	cellC, err := s.Pop()
	if err != nil {
		return err
	}
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	_ = s.Push(cellC)
	_ = s.Push(cellA)
	_ = s.Push(cellB)
	return nil
}

func (s *Stack) IsEqual() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellBoolean:
		if b, ok := cellB.(CellBoolean); ok {
			return s.Push(CellBoolean(a == b))
		}
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellBoolean(a == b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellBoolean(a == b))
		}
	case CellText:
		if b, ok := cellB.(CellText); ok {
			return s.Push(CellBoolean(a == b))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) IsGreater() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellBoolean(a > b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellBoolean(a > b))
		}
	case CellText:
		if b, ok := cellB.(CellText); ok {
			return s.Push(CellBoolean(a > b))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) IsSmaller() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellBoolean(a < b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellBoolean(a < b))
		}
	case CellText:
		if b, ok := cellB.(CellText); ok {
			return s.Push(CellBoolean(a < b))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) Add() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellInteger(a + b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellFloat(a + b))
		}
	case CellText:
		if b, ok := cellB.(CellText); ok {
			return s.Push(CellText(a + b))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) Subtract() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellInteger(a - b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellFloat(a - b))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) Multiply() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellInteger(a * b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellFloat(a * b))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) Divide() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellInteger(a / b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellFloat(a / b))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) Modulo() error {
	cellB, err := s.Pop()
	if err != nil {
		return err
	}
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		if b, ok := cellB.(CellInteger); ok {
			return s.Push(CellInteger(a % b))
		}
	case CellFloat:
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellFloat(math.Mod(float64(a), float64(b))))
		}
		if b, ok := cellB.(CellFloat); ok {
			return s.Push(CellFloat(math.Mod(float64(a), float64(b))))
		}
	}
	return newTypeMismatchError(cellA, cellB)
}

func (s *Stack) ToText() error {
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellText:
		return s.Push(a)
	case CellInteger:
		return s.Push(CellText(strconv.FormatInt(int64(a), 10)))
	case CellFloat:
		return s.Push(CellText(strconv.FormatFloat(float64(a), 'f', 5, 64)))
	}
}

func (s *Stack) ToInteger() error {
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellInteger:
		return s.Push(a)
	case CellText:
		n, err := strconv.Atoi(string(a))
		if err != nil {
			return s.Push(CellText(err.Error()))
		}
		return s.Push(CellInteger(n))
	case CellFloat:
		return s.Push(CellInteger(a))
	}
}

func (s *Stack) Invert() error {
	cellA, err := s.Pop()
	if err != nil {
		return err
	}
	switch a := cellA.(type) {
	default:
		return newInvalidTypeError(a)
	case CellBoolean:
		return s.Push(!a)
	}
}
