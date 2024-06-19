package string_scanner

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Simple String Scanner", func() {
	DescribeTable("Scanning", func(value string, expectedTokens []Token) {
		scanner := NewSimpleScanner()
		scanner.Init(value)
		allTokens := []Token{}
		for scanner.Next() {
			allTokens = append(allTokens, *scanner.Token())
		}
		Expect(allTokens).To(Equal(expectedTokens))
	},
		Entry("Empty string", "", []Token{}),
		Entry("Testing 1 token", "a", []Token{{TokenType: ALPHA, Value: "a", Position: 0}}),
		Entry("Testing 5 tokens", "ab(1.", []Token{
			{TokenType: ALPHA, Value: "a", Position: 0},
			{TokenType: ALPHA, Value: "b", Position: 1},
			{TokenType: SYMBOL, Value: "(", Position: 2},
			{TokenType: DIGIT, Value: "1", Position: 3},
			{TokenType: DECIMALPOINT, Value: ".", Position: 4},
		}),
	)

	DescribeTable("Peek", func(scanner Scanner, returnedBool bool, returnedValue *Token) {
		got, gotVal := scanner.Peek()
		Expect(got).To(Equal(returnedBool))
		Expect(gotVal).To(Equal(returnedValue))
	},
		Entry("return true and update token if pos < length of value -1",
			&simpleStringScanner{
				pos:   1,
				value: "testValue",
			},
			true,
			&Token{
				TokenType: 0,
				Value:     "s",
				Position:  2,
			},
		),
		Entry("return false and nil if pos > length of value -1",
			&simpleStringScanner{
				pos:   10,
				value: "testValue",
			},
			false, nil,
		),
		Entry("return false and nil if pos == length of value -1",
			&simpleStringScanner{
				pos:   8,
				value: "testValue",
			},
			false, nil,
		),
	)
})
