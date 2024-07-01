package sql_parser

import (
	"github.com/openshift-online/ocm-common/pkg/utils/parser/string_parser"
	"strings"
)

type SQLParserOption func(parser *sqlParser)

func WithValidColumns(validColumns ...string) SQLParserOption {
	return func(parser *sqlParser) {
		parser.validColumns = validColumns
	}
}

func WithColumnPrefix(columnPrefix string) SQLParserOption {
	return func(parser *sqlParser) {
		parser.columnPrefix = strings.Trim(columnPrefix, " ")
	}
}

func WithMaximumComplexity(maximumComplexity int) SQLParserOption {
	return func(parser *sqlParser) {
		parser.maximumComplexity = maximumComplexity
	}
}

func NewSQLParser(options ...SQLParserOption) SQLParser {
	parser := &sqlParser{
		maximumComplexity: defaultMaximumComplexity,
	}

	for _, option := range options {
		option(parser)
	}

	stringParser := string_parser.NewStringParserBuilder().
		WithGrammar(BasicSQLGrammar()).
		WithTransitionInterceptor(parser.transitionInterceptor).
		WithScanner(NewSQLScanner()).
		Build()

	parser.parser = stringParser
	return parser
}
