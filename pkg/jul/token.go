package jul

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Token represents a logical part of the source code.
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

func (t Token) String() string { return fmt.Sprintf("%s (%s) %q", t.Type, t.Position, t.Value) }

type Position struct{ Line, Column int }

func (p Position) String() string { return strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Column) }

type TokenType string

const (
	TokenTypeEOF             TokenType = "EOF"                // Reached EOF
	TokenTypeFunctionCall    TokenType = "function-call"      // Function call (incl. literal number)
	TokenTypeQuotation       TokenType = "anonymous-function" // Raw function body (= quotation)
	TokenTypeLiteralText     TokenType = "literal-text"       // Literal text string
	TokenTypeLiteralTextWord TokenType = "literal-text-word"  // Text without spaces
	TokenTypeComment         TokenType = "comment"            // Code comments and remarks
)

const (
	MarkAnonymousFunctionStart = '['
	MarkAnonymousFunctionEnd   = ']'
	MarkLiteralTextQuote       = '"'
	MarkLiteralTextWordStart   = '*'
	MarkCommentStart           = '('
	MarkCommentEnd             = ')'
)

// Source reads source code tokens from the underlying reader.
type Source struct {
	r io.Reader
	p Position
}

func NewSource(r io.Reader) *Source { return &Source{r: r, p: Position{1, 1}} }

// Next returns the next token from the source.
// When EOF is reached, a token of type EOF is returned,
// errors may be due to illegal syntax error or a failed read on the underlying reader.
func (src *Source) Next() (Token, error) {
	for {
		start := src.p
		c, err := src.read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{Position: start, Type: TokenTypeEOF}, nil
			}
			return Token{}, err
		}
		switch {
		case isSpace(c):
			// Skip whitespace
			continue
		case c == MarkAnonymousFunctionStart:
			// Tokenize quotation
			v, err := src.readEnclosed(MarkAnonymousFunctionStart, MarkAnonymousFunctionEnd)
			if err != nil {
				return Token{}, err
			}
			return Token{Position: start, Type: TokenTypeQuotation, Value: string(v)}, nil
		case c == MarkCommentStart:
			// Tokenize comment
			v, err := src.readEnclosed(MarkCommentStart, MarkCommentEnd)
			if err != nil {
				return Token{}, err
			}
			return Token{Position: start, Type: TokenTypeComment, Value: string(v)}, nil
		case c == MarkLiteralTextQuote:
			// Tokenize literal text
			var v []byte
			isEscaped := false
			for {
				c, err := src.read()
				if errors.Is(err, io.EOF) || (c == MarkLiteralTextQuote && !isEscaped) {
					break
				} else if err != nil {
					return Token{}, err
				}
				if isEscaped {
					switch c {
					case 'n':
						c = '\n'
					case 't':
						c = '\t'
					}
					isEscaped = false
				} else if !isEscaped && c == '\\' {
					isEscaped = true
					continue
				}
				v = append(v, c)
			}
			return Token{Position: start, Type: TokenTypeLiteralText, Value: string(v)}, nil
		case c == MarkLiteralTextWordStart:
			// Tokenize literal word
			rest, err := src.readWhile(isNotSpace)
			if err != nil {
				return Token{}, err
			}
			return Token{Position: start, Type: TokenTypeLiteralTextWord, Value: string(rest)}, nil
		case isPrintable(c) && c != MarkAnonymousFunctionEnd:
			// Tokenize call
			v := []byte{c}
			rest, err := src.readWhile(isNotSpace)
			if err != nil {
				return Token{}, err
			}
			v = append(v, rest...)
			return Token{Position: start, Type: TokenTypeFunctionCall, Value: string(v)}, nil
		default:
			return Token{}, syntaxError{Position: src.p, Message: fmt.Sprintf("unexpected character %q", c)}
		}
	}
}

func (src *Source) Tokens() ([]Token, error) {
	var out []Token
	for {
		tok, err := src.Next()
		if err != nil {
			return out, err
		}
		out = append(out, tok)
		if tok.Type == TokenTypeEOF {
			break
		}
	}
	return out, nil
}

func (src *Source) read() (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(src.r, buf)
	if err != nil {
		return 0, err
	}
	c := buf[0]
	if c == '\n' {
		src.p.Line++
		src.p.Column = 1
	} else {
		src.p.Column++
	}
	return c, nil
}

func (src *Source) readWhile(predicate func(byte) bool) ([]byte, error) {
	var out []byte
	for {
		c, err := src.read()
		if errors.Is(err, io.EOF) || !predicate(c) {
			break
		} else if err != nil {
			return out, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (src *Source) readEnclosed(markStart, markEnd byte) ([]byte, error) {
	depth := 1
	var v []byte
	missingClosingErrMsg := fmt.Sprintf("missing closing character: %q", markEnd)
	for {
		c, err := src.read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return v, syntaxError{Position: src.p, Message: missingClosingErrMsg}
			}
			return v, err
		}
		if c == markStart {
			depth++
		} else if c == markEnd {
			depth--
		}

		if c == markEnd && depth == 0 {
			break
		}
		v = append(v, c)
	}
	if depth > 0 {
		return v, syntaxError{Position: src.p, Message: missingClosingErrMsg}
	}
	return v, nil
}

type syntaxError struct {
	Position Position
	Message  string
}

func (err syntaxError) Error() string { return fmt.Sprintf("%s (%s)", err.Message, err.Position) }

func isSpace(c byte) bool     { return c == ' ' || c == '\n' || c == '\t' }
func isNotSpace(c byte) bool  { return !isSpace(c) }
func isPrintable(c byte) bool { return c >= 33 && c <= 126 }
