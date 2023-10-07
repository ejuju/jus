package jul

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	_ "embed"
)

//go:embed prelude.ju
var Prelude string

type VM struct {
	stack      *Stack
	dictionary *Dictionary
	rrand      *rand.Rand
	ui         UI
	conn       *net.TCPConn
}

type Option func(vm *VM)

func WithStack(s *Stack) Option                  { return func(vm *VM) { vm.stack = s } }
func WithDictionary(d *Dictionary) Option        { return func(vm *VM) { vm.dictionary = d } }
func WithUI(ui UI) Option                        { return func(vm *VM) { vm.ui = ui } }
func WithServerConnection(c *net.TCPConn) Option { return func(vm *VM) { vm.conn = c } }
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
	if vm.ui == nil {
		vm.ui = NewDefaultUI(nil, nil)
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

func RunCLI() {
	// Start in REPL or file mode
	from := os.Stdin
	if len(os.Args) > 1 {
		var err error
		from, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()
	}

	// Execute code from stdin or file
	vm := NewVM()
	err := vm.Execute(from)
	if err != nil {
		log.Println(err)
	}
}
