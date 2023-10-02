package jul

import (
	"strings"
	"testing"
)

func TestPrelude(t *testing.T) {
	t.Run("random: generates random number between 0 and 1 (exclusive)", func(t *testing.T) {
		vm := NewVM()
		err := vm.Execute(strings.NewReader("1 random"))
		if err != nil {
			panic(err)
		}
		c, err := vm.stack.Pop()
		if err != nil {
			panic(err)
		}
		result, ok := c.(CellInteger)
		if !ok {
			t.Fatalf("got type %T", c)
		}
		if result != 0 {
			t.Fatalf("got random number %d instead of %d", result, 0)
		}
	})

	t.Run("can generate random number in range (random-between)", func(t *testing.T) {
		vm := NewVM(WithRandomSeed(1))
		err := vm.Execute(strings.NewReader("1 2 random-between"))
		if err != nil {
			panic(err)
		}
		c, err := vm.stack.Pop()
		if err != nil {
			panic(err)
		}
		result, ok := c.(CellInteger)
		if !ok {
			t.Fatalf("got type %T", c)
		}
		if result != 1 {
			t.Fatalf("got random number %d instead of %d", result, 1)
		}
	})
}
