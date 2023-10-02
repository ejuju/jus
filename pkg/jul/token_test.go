package jul

import (
	"reflect"
	"strings"
	"testing"
)

func TestToken(t *testing.T) {
	tests := []struct {
		desc   string
		input  string
		output []Token
	}{
		{
			desc:  string(TokenTypeQuotation),
			input: "[dup write]",
			output: []Token{
				{Position: Position{Line: 1, Column: 1}, Type: TokenTypeQuotation, Value: "dup write"},
				{Position: Position{Line: 1, Column: 12}, Type: TokenTypeEOF},
			},
		},
		{
			desc:  string(TokenTypeLiteralText),
			input: `"Hello\n\t world!"`,
			output: []Token{
				{Position: Position{Line: 1, Column: 1}, Type: TokenTypeLiteralText, Value: "Hello\n\t world!"},
				{Position: Position{Line: 1, Column: 19}, Type: TokenTypeEOF},
			},
		},
		{
			desc:  string(TokenTypeLiteralTextWord),
			input: "*single-word",
			output: []Token{
				{Position: Position{Line: 1, Column: 1}, Type: TokenTypeLiteralTextWord, Value: "single-word"},
				{Position: Position{Line: 1, Column: 13}, Type: TokenTypeEOF},
			},
		},
		{
			desc:  string(TokenTypeFunctionCall),
			input: "execute-my-function",
			output: []Token{
				{Position: Position{Line: 1, Column: 1}, Type: TokenTypeFunctionCall, Value: "execute-my-function"},
				{Position: Position{Line: 1, Column: 20}, Type: TokenTypeEOF},
			},
		},
		{
			desc:  string(TokenTypeComment),
			input: "(I'm a comment)",
			output: []Token{
				{Position: Position{Line: 1, Column: 1}, Type: TokenTypeComment, Value: "I'm a comment"},
				{Position: Position{Line: 1, Column: 16}, Type: TokenTypeEOF},
			},
		},
		{
			desc:   string(TokenTypeEOF),
			input:  "",
			output: []Token{{Position: Position{Line: 1, Column: 1}, Type: TokenTypeEOF}},
		},
		{
			desc:   string(TokenTypeEOF) + " with leading spaces",
			input:  "\n\t ",
			output: []Token{{Position: Position{Line: 2, Column: 3}, Type: TokenTypeEOF}},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			toks, err := NewSource(strings.NewReader(test.input)).Tokens()
			if err != nil {
				panic(err)
			}
			if !reflect.DeepEqual(toks, test.output) {
				t.Fatalf("got %+v instead of %+v", toks, test.output)
			}
		})
	}
}
