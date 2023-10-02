package jul

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "embed"
)

//go:embed prelude.ju
var Prelude string

type VM struct {
	stack          *Stack
	dictionary     *Dictionary
	stdin          io.Reader
	stdout, stderr io.Writer
	rrand          *rand.Rand
}

type Option func(vm *VM)

func WithStack(s *Stack) Option           { return func(vm *VM) { vm.stack = s } }
func WithDictionary(d *Dictionary) Option { return func(vm *VM) { vm.dictionary = d } }
func WithStdin(stdin io.Reader) Option    { return func(vm *VM) { vm.stdin = stdin } }
func WithStdout(stdout io.Writer) Option  { return func(vm *VM) { vm.stdout = stdout } }
func WithStderr(stderr io.Writer) Option  { return func(vm *VM) { vm.stderr = stderr } }
func WithRandomSeed(seed int64) Option {
	return func(vm *VM) { vm.rrand = rand.New(rand.NewSource(seed)) }
}

func NewVM(opts ...Option) *VM {
	vm := &VM{}
	for _, opt := range opts {
		opt(vm)
	}

	// Set defaults if needed
	if vm.stack == nil {
		vm.stack = NewStack(0)
	}
	if vm.dictionary == nil {
		vm.dictionary = NewDictionary()
	}
	if vm.stdin == nil {
		vm.stdin = os.Stdin
	}
	if vm.stdout == nil {
		vm.stdout = os.Stdout
	}
	if vm.stderr == nil {
		vm.stderr = os.Stdin
	}
	if vm.rrand == nil {
		vm.rrand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// Execute prelude
	err := vm.Execute(strings.NewReader(Prelude))
	if err != nil {
		panic(err)
	}

	return vm
}

func (vm *VM) Execute(r io.Reader) error {
	src := NewSource(r)
	for {
		tok, err := src.Next()
		if err != nil {
			return err
		}
		switch tok.Type {
		default:
			panic(fmt.Errorf("unreachable: unhandled token type %q", tok.Type))
		case TokenTypeEOF:
			return nil
		case TokenTypeComment:
			continue
		case TokenTypeFunctionCall:
			w := vm.dictionary.FindLatestDefinition(tok.Value)
			if w == nil {
				num, err := strconv.Atoi(tok.Value)
				if err != nil {
					return RuntimeError{Position: src.p, Cause: fmt.Errorf("unknown word %q", tok.Value)}
				}
				err = vm.stack.Push(CellInteger(int(num)))
				if err != nil {
					return RuntimeError{Position: src.p, Cause: err}
				}
				continue
			}
			err = w.Func(vm)
			if err != nil {
				return RuntimeError{Position: src.p, Cause: fmt.Errorf("%s: %w", w.Name, err)}
			}
		case TokenTypeQuotation:
			err = vm.stack.Push(CellQuotation(tok.Value))
			if err != nil {
				return RuntimeError{Position: src.p, Cause: err}
			}
		case TokenTypeLiteralText, TokenTypeLiteralTextWord:
			err = vm.stack.Push(CellText(tok.Value))
			if err != nil {
				return RuntimeError{Position: src.p, Cause: err}
			}
		}
	}
}

type RuntimeError struct {
	Position Position
	Cause    error
}

func (err RuntimeError) Error() string { return fmt.Sprintf("\n(at %s) %s", err.Position, err.Cause) }
