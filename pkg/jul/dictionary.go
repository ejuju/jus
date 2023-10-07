package jul

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ejuju/jus/pkg/jutp"
)

type Dictionary struct{ words []*Definition }

func NewDictionary() *Dictionary {
	return &Dictionary{words: Builtins}
}

func (d *Dictionary) FindLatestDefinition(name string) *Definition {
	for i := len(d.words) - 1; i >= 0; i-- {
		if w := d.words[i]; w.Name == name {
			return w
		}
	}
	return nil
}

func (d *Dictionary) Define(w *Definition) error {
	if w := d.FindLatestDefinition(w.Name); w != nil {
		return fmt.Errorf("already defined word: %q", w.Name)
	}
	d.words = append(d.words, w)
	return nil
}

type Definition struct {
	Name string
	Func func(vm *VM) error
}

var Builtins = []*Definition{
	{Name: "drop", Func: func(vm *VM) error { return vm.stack.Drop() }},
	{Name: "pick", Func: func(vm *VM) error { return vm.stack.Pick() }},
	{Name: "swap", Func: func(vm *VM) error { return vm.stack.Swap() }},
	{Name: "rot", Func: func(vm *VM) error { return vm.stack.Rot() }},
	{Name: "is-equal", Func: func(vm *VM) error { return vm.stack.IsEqual() }},
	{Name: "is-greater", Func: func(vm *VM) error { return vm.stack.IsGreater() }},
	{Name: "is-smaller", Func: func(vm *VM) error { return vm.stack.IsSmaller() }},
	{Name: "add", Func: func(vm *VM) error { return vm.stack.Add() }},
	{Name: "subtract", Func: func(vm *VM) error { return vm.stack.Subtract() }},
	{Name: "multiply", Func: func(vm *VM) error { return vm.stack.Multiply() }},
	{Name: "divide", Func: func(vm *VM) error { return vm.stack.Divide() }},
	{Name: "modulo", Func: func(vm *VM) error { return vm.stack.Modulo() }},
	{Name: "to-integer", Func: func(vm *VM) error { return vm.stack.ToInteger() }},
	{Name: "to-text", Func: func(vm *VM) error { return vm.stack.ToText() }},
	{Name: "invert", Func: func(vm *VM) error { return vm.stack.Invert() }},
	{
		Name: "do",
		Func: func(vm *VM) error {
			cellA, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			quotation, ok := cellA.(CellQuotation)
			if !ok {
				return fmt.Errorf("got (A) %T instead of quotation", cellA)
			}
			return vm.Execute(strings.NewReader(string(quotation)))
		},
	},
	{
		Name: "define",
		Func: func(vm *VM) error {
			// Pop function body
			cellB, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("pop (B) function body (quotation): %w", err)
			}
			quotation, ok := cellB.(CellQuotation)
			if !ok {
				return fmt.Errorf("got (B) %T instead of quotation", cellB)
			}

			// Pop keyword
			cellA, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("pop (A) keyword (text): %w", err)
			}
			keyword, ok := cellA.(CellText)
			if !ok {
				return fmt.Errorf("got (A) %T instead of text", cellA)
			}

			// Add word to dictionary
			return vm.dictionary.Define(&Definition{
				Name: string(keyword),
				Func: func(vm *VM) error {
					return vm.Execute(strings.NewReader(string(quotation)))
				},
			})
		},
	},
	{
		Name: "if",
		Func: func(vm *VM) error {
			// Pop falsy callback
			cellC, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("pop (C) falsy callback (quotation): %w", err)
			}
			callbackIfFalse, ok := cellC.(CellQuotation)
			if !ok {
				return fmt.Errorf("got (C) %T instead of quotation", cellC)
			}

			// Pop truthy callback
			cellB, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("pop (B) truthy callback (quotation): %w", err)
			}
			callbackIfTrue, ok := cellB.(CellQuotation)
			if !ok {
				return fmt.Errorf("got (B) %T instead of quotation", cellB)
			}

			// Pop boolean
			cellA, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("pop (A) boolean: %w", err)
			}
			boolean, ok := cellA.(CellBoolean)
			if !ok {
				return fmt.Errorf("got (A) %T instead of boolean", cellA)
			}

			// Execute callback depending on boolean
			if boolean {
				return vm.Execute(strings.NewReader(string(callbackIfTrue)))
			} else {
				return vm.Execute(strings.NewReader(string(callbackIfFalse)))
			}
		},
	},
	{
		Name: "repeat",
		Func: func(vm *VM) error {
			// Pop callback
			cellA, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("pop (B) callback (quotation): %w", err)
			}
			callback, ok := cellA.(CellQuotation)
			if !ok {
				return fmt.Errorf("got (B) %T instead of quotation", cellA)
			}

			// Execute callback in a loop while the callback returns true
			for i := 0; true; i++ {
				// Push iteration count
				err = vm.stack.Push(CellInteger(i))
				if err != nil {
					return fmt.Errorf("push iteration count before callback: %w", err)
				}

				// Execute callback
				err = vm.Execute(strings.NewReader(string(callback)))
				if err != nil {
					return fmt.Errorf("executing callback (%d): %w", i, err)
				}

				// Pop boolean and exit loop if set to false
				cellAFromCallback, err := vm.stack.Pop()
				if err != nil {
					return fmt.Errorf("pop (A) boolean: %w", err)
				}
				boolean, ok := cellAFromCallback.(CellBoolean)
				if !ok {
					return fmt.Errorf("got (A) %T instead of boolean", cellAFromCallback)
				}
				if !boolean {
					break
				}
			}
			return nil
		},
	},
	{
		Name: "write",
		Func: func(vm *VM) error {
			cellA, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			switch a := (cellA).(type) {
			case CellText:
				return vm.ui.Write(string(a))
			case CellBoolean:
				return vm.ui.Write(strconv.FormatBool(bool(a)))
			case CellInteger:
				return vm.ui.Write(strconv.Itoa(int(a)))
			case CellFloat:
				return vm.ui.Write(strconv.FormatFloat(float64(a), 'f', 2, 64))
			}
			return newInvalidTypeError(cellA)
		},
	},
	{
		Name: "read",
		Func: func(vm *VM) error {
			line, err := vm.ui.Read()
			if err != nil {
				return err
			}
			return vm.stack.Push(CellText(line))
		},
	},
	{
		Name: "random",
		Func: func(vm *VM) error {
			cellA, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			switch a := cellA.(type) {
			default:
				return newInvalidTypeError(a)
			case CellInteger:
				if a <= 0 {
					return fmt.Errorf("can't generate integer in range [0;%d[", a)
				}
				return vm.stack.Push(CellInteger(vm.rrand.Intn(int(a))))
			case CellFloat:
				return vm.stack.Push(CellInteger(vm.rrand.Float64() * float64(a)))
			}
		},
	},
	{
		Name: "now",
		Func: func(vm *VM) error { return vm.stack.Push(CellTime(time.Now())) },
	},
	{
		Name: "wait",
		Func: func(vm *VM) error {
			cellA, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			switch a := cellA.(type) {
			default:
				return newInvalidTypeError(a)
			case CellInteger:
				time.Sleep(time.Duration(a) * time.Millisecond)
			case CellFloat:
				time.Sleep(time.Duration(float64(a) * float64(time.Second)))
			case CellTime:
				time.Sleep(time.Until(time.Time(a)))
			}
			return nil
		},
	},
	{
		Name: "retrieve",
		Func: func(vm *VM) error {
			if vm.conn == nil {
				return errors.New("not connected to server")
			}
			cellA, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			switch a := cellA.(type) {
			case CellText:
				_, err = jutp.Write(vm.conn, jutp.Message(a))
				return err
			}
			return newInvalidTypeError(cellA)
		},
	},
}
