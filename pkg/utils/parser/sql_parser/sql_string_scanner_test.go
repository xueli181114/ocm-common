package sql_parser

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-common/pkg/utils/parser/string_scanner"
)

var _ = Describe("SQL String Scanner", func() {
	makeToken := func(tokenType int, value string, pos int) string_scanner.Token {
		return string_scanner.Token{
			TokenType: tokenType,
			Value:     value,
			Position:  pos,
		}
	}

	DescribeTable("scanning", func(value string, expectedTokens []string_scanner.Token) {
		scanner := NewSQLScanner()
		scanner.Init(value)
		var allTokens []string_scanner.Token
		for scanner.Next() {
			allTokens = append(allTokens, *scanner.Token())
		}
		Expect(allTokens).To(Equal(expectedTokens))
	},
		Entry("Simple select",
			"SELECT * FROM TABLE_NAME",
			[]string_scanner.Token{
				makeToken(LITERAL, "SELECT", 0),
				makeToken(LITERAL, "*", 7),
				makeToken(LITERAL, "FROM", 9),
				makeToken(LITERAL, "TABLE_NAME", 14),
			},
		),
		Entry("Select with quoted string",
			"SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'surname with spaces'",
			[]string_scanner.Token{
				makeToken(LITERAL, "SELECT", 0),
				makeToken(LITERAL, "*", 7),
				makeToken(LITERAL, "FROM", 9),
				makeToken(LITERAL, "ADDRESS_BOOK", 14),
				makeToken(LITERAL, "WHERE", 27),
				makeToken(LITERAL, "SURNAME", 33),
				makeToken(OP, "=", 41),
				makeToken(QUOTED_LITERAL, "'surname with spaces'", 43),
			},
		),
		Entry("Select with quoted string including a comma",
			"SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'surname with , comma'",
			[]string_scanner.Token{
				makeToken(LITERAL, "SELECT", 0),
				makeToken(LITERAL, "*", 7),
				makeToken(LITERAL, "FROM", 9),
				makeToken(LITERAL, "ADDRESS_BOOK", 14),
				makeToken(LITERAL, "WHERE", 27),
				makeToken(LITERAL, "SURNAME", 33),
				makeToken(OP, "=", 41),
				makeToken(QUOTED_LITERAL, "'surname with , comma'", 43),
			},
		),
		Entry("Select with quoted string including an open parenthesis",
			"SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'surname with ( parenthesis'",
			[]string_scanner.Token{
				makeToken(LITERAL, "SELECT", 0),
				makeToken(LITERAL, "*", 7),
				makeToken(LITERAL, "FROM", 9),
				makeToken(LITERAL, "ADDRESS_BOOK", 14),
				makeToken(LITERAL, "WHERE", 27),
				makeToken(LITERAL, "SURNAME", 33),
				makeToken(OP, "=", 41),
				makeToken(QUOTED_LITERAL, "'surname with ( parenthesis'", 43),
			},
		),
		Entry("Select with quoted string including escaped chars",
			`SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'surname with spaces and \'quote\''`,
			[]string_scanner.Token{
				makeToken(LITERAL, "SELECT", 0),
				makeToken(LITERAL, "*", 7),
				makeToken(LITERAL, "FROM", 9),
				makeToken(LITERAL, "ADDRESS_BOOK", 14),
				makeToken(LITERAL, "WHERE", 27),
				makeToken(LITERAL, "SURNAME", 33),
				makeToken(OP, "=", 41),
				makeToken(QUOTED_LITERAL, `'surname with spaces and \'quote\''`, 43),
			},
		),
		Entry("SQL with operators",
			`SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'Mouse' AND AGE > 3`,
			[]string_scanner.Token{
				makeToken(LITERAL, "SELECT", 0),
				makeToken(LITERAL, "*", 7),
				makeToken(LITERAL, "FROM", 9),
				makeToken(LITERAL, "ADDRESS_BOOK", 14),
				makeToken(LITERAL, "WHERE", 27),
				makeToken(LITERAL, "SURNAME", 33),
				makeToken(OP, "=", 41),
				makeToken(QUOTED_LITERAL, `'Mouse'`, 43),
				makeToken(LITERAL, "AND", 51),
				makeToken(LITERAL, "AGE", 55),
				makeToken(OP, ">", 59),
				makeToken(LITERAL, "3", 61),
			},
		),
		Entry("SQL with empty parenthesis",
			"name IN ()",
			[]string_scanner.Token{
				makeToken(LITERAL, "name", 0),
				makeToken(LITERAL, "IN", 5),
				makeToken(BRACE, "(", 8),
				makeToken(BRACE, ")", 9),
			}),
		Entry("LIST VALUES",
			"value1, 'value2', 'value3', value4",
			[]string_scanner.Token{
				makeToken(LITERAL, "value1", 0),
				makeToken(LITERAL, ",", 6),
				makeToken(QUOTED_LITERAL, "'value2'", 8),
				makeToken(LITERAL, ",", 16),
				makeToken(QUOTED_LITERAL, "'value3'", 18),
				makeToken(LITERAL, ",", 26),
				makeToken(LITERAL, "value4", 28),
			}),
		Entry("QUOTED STRING with special characters",
			`name = '@,\'""(){}/'`,
			[]string_scanner.Token{
				makeToken(LITERAL, "name", 0),
				makeToken(OP, "=", 5),
				makeToken(QUOTED_LITERAL, `'@,\'""(){}/'`, 7),
			}),
		Entry("SQL with JSONB",
			`select * from table where manifest->'data'->'manifest'->'metadata'->'labels'->>'foo' = 'bar'`,
			[]string_scanner.Token{
				makeToken(LITERAL, "select", 0),
				makeToken(LITERAL, "*", 7),
				makeToken(LITERAL, "from", 9),
				makeToken(LITERAL, "table", 14),
				makeToken(LITERAL, "where", 20),
				makeToken(LITERAL, "manifest", 26),
				makeToken(OP, "->", 34),
				makeToken(QUOTED_LITERAL, "'data'", 36),
				makeToken(OP, "->", 42),
				makeToken(QUOTED_LITERAL, "'manifest'", 44),
				makeToken(OP, "->", 54),
				makeToken(QUOTED_LITERAL, "'metadata'", 56),
				makeToken(OP, "->", 66),
				makeToken(QUOTED_LITERAL, "'labels'", 68),
				makeToken(OP, "->>", 76),
				makeToken(QUOTED_LITERAL, "'foo'", 79),
				makeToken(OP, "=", 85),
				makeToken(QUOTED_LITERAL, "'bar'", 87),
			},
		),
		Entry("SQL with JSONB contains token",
			`resources.payload -> 'data' -> 'manifests' @> '[{"metadata":{"labels":{"foo":"bar"}}}]'`,
			[]string_scanner.Token{
				makeToken(LITERAL, "resources.payload", 0),
				makeToken(OP, "->", 18),
				makeToken(QUOTED_LITERAL, "'data'", 21),
				makeToken(OP, "->", 28),
				makeToken(QUOTED_LITERAL, "'manifests'", 31),
				makeToken(OP, "@>", 43),
				makeToken(QUOTED_LITERAL, `'[{"metadata":{"labels":{"foo":"bar"}}}]'`, 46),
			},
		),
	)

	DescribeTable("peeking", func(scanner string_scanner.Scanner, wantedBool bool, wantedResult *string_scanner.Token) {
		got, gotVal := scanner.Peek()
		Expect(got).To(Equal(wantedBool))
		Expect(gotVal).To(Equal(wantedResult))
	},
		Entry("return true and token if pos < length of tokens -1",
			&scanner{
				pos: 1,
				tokens: []string_scanner.Token{
					{Value: "testToken1"},
					{Value: "testToken2"},
					{Value: "testToken3"},
				},
			},
			true,
			&string_scanner.Token{Value: "testToken3"},
		),
		Entry("return false and nil if pos < length of tokens -1",
			&scanner{
				pos: 2,
				tokens: []string_scanner.Token{
					{Value: "testToken1"},
					{Value: "testToken2"},
				},
			},
			false,
			nil,
		),
		Entry("return false and nil if pos == length of tokens -1",
			&scanner{
				pos: 1,
				tokens: []string_scanner.Token{
					{Value: "testToken1"},
					{Value: "testToken2"},
				},
			},
			false,
			nil,
		),
	)
})
