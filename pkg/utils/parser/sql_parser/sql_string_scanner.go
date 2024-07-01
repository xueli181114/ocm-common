package sql_parser

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/parser/string_scanner"
)

const (
	OP = iota
	BRACE
	LITERAL
	QUOTED_LITERAL
	NO_TOKEN
)

// scanner - This scanner is to be used to parse SQL Strings. It splits the provided string by whole words
// or sentences if it finds quotes. Nested round braces are supported too.
type scanner struct {
	tokens []string_scanner.Token
	pos    int
}

var _ string_scanner.Scanner = &scanner{}

// Init feeds the scanner with the text to be scanned
func (s *scanner) Init(txt string) {
	s.pos = -1
	s.tokens = nil

	var tokens []string_scanner.Token
	currentTokenType := NO_TOKEN

	quoted := false
	escaped := false

	sendCurrentTokens := func() {
		res := ""
		for _, token := range tokens {
			res += token.Value
		}
		if res != "" {
			s.tokens = append(s.tokens, string_scanner.Token{TokenType: currentTokenType, Value: res, Position: tokens[0].Position})
		}
		tokens = nil
		currentTokenType = NO_TOKEN
	}

	// extract all the tokens from the string
	for i, currentChar := range txt {
		switch {
		case currentChar == '\'' && quoted:
			tokens = append(tokens, string_scanner.Token{TokenType: QUOTED_LITERAL, Value: "'", Position: i})
			if !escaped {
				sendCurrentTokens()
				quoted = false
				currentTokenType = NO_TOKEN
			}
			escaped = false
		case currentChar == '\\' && quoted:
			escaped = true
			tokens = append(tokens, string_scanner.Token{TokenType: QUOTED_LITERAL, Value: "\\", Position: i})
		case quoted: // everything that is not "'" or '\' must be added to the current token if quoted is true
			tokens = append(tokens, string_scanner.Token{TokenType: LITERAL, Value: string(currentChar), Position: i})
		case currentChar == ' ':
			sendCurrentTokens()
		case currentChar == ',':
			sendCurrentTokens()
			s.tokens = append(s.tokens, string_scanner.Token{TokenType: LITERAL, Value: string(currentChar), Position: i})
		case currentChar == '\'':
			sendCurrentTokens()
			quoted = true
			currentTokenType = QUOTED_LITERAL
			tokens = append(tokens, string_scanner.Token{TokenType: OP, Value: "'", Position: i})
		case currentChar == '\\':
			if currentTokenType != NO_TOKEN && currentTokenType != LITERAL && currentTokenType != QUOTED_LITERAL {
				sendCurrentTokens()
			}
			currentTokenType = LITERAL
			tokens = append(tokens, string_scanner.Token{TokenType: LITERAL, Value: `\`, Position: i})
		case currentChar == '@', currentChar == '-', currentChar == '=', currentChar == '<', currentChar == '>':
			// found op Token
			if currentTokenType != NO_TOKEN && currentTokenType != OP {
				sendCurrentTokens()
			}
			tokens = append(tokens, string_scanner.Token{TokenType: OP, Value: string(currentChar), Position: i})
			currentTokenType = OP
		case currentChar == '(', currentChar == ')':
			sendCurrentTokens()
			s.tokens = append(s.tokens, string_scanner.Token{TokenType: BRACE, Value: string(currentChar), Position: i})
		default:
			if currentTokenType != NO_TOKEN && currentTokenType != LITERAL && currentTokenType != QUOTED_LITERAL {
				sendCurrentTokens()
			}
			currentTokenType = LITERAL
			tokens = append(tokens, string_scanner.Token{TokenType: LITERAL, Value: string(currentChar), Position: i})
		}
	}

	sendCurrentTokens()
}

// Next moves to the next token and return `true` if another token is present. Otherwise returns `false`
func (s *scanner) Next() bool {
	if s.pos < (len(s.tokens) - 1) {
		s.pos++
		return true
	}
	return false
}

// Peek looks if another token is present after the current position without moving the cursor
func (s *scanner) Peek() (bool, *string_scanner.Token) {
	if s.pos < (len(s.tokens) - 1) {
		ret := s.tokens[s.pos+1]
		return true, &ret
	}
	return false, nil
}

// Token returns the current token
func (s *scanner) Token() *string_scanner.Token {
	if s.pos < 0 || s.pos >= len(s.tokens) {
		panic(fmt.Errorf("invalid scanner position %d", s.pos))
	}
	ret := s.tokens[s.pos]
	return &ret
}

func NewSQLScanner() string_scanner.Scanner {
	return &scanner{
		pos: -1,
	}
}
