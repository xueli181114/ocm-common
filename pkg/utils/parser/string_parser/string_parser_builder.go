package string_parser

import (
	"github.com/openshift-online/ocm-common/pkg/utils/parser/state_machine"
	"github.com/openshift-online/ocm-common/pkg/utils/parser/string_scanner"
)

type StringParserBuilder struct {
	grammar     Grammar
	scanner     string_scanner.Scanner
	interceptor state_machine.TransitionInterceptor[string, string]
	observers   []state_machine.TransitionObserver[string, string]
}

func (spb *StringParserBuilder) WithScanner(scanner string_scanner.Scanner) *StringParserBuilder {
	spb.scanner = scanner
	return spb
}

func (spb *StringParserBuilder) WithGrammar(grammar Grammar) *StringParserBuilder {
	spb.grammar = grammar
	return spb
}

func (spb *StringParserBuilder) WithTransitionInterceptor(interceptor state_machine.TransitionInterceptor[string, string]) *StringParserBuilder {
	spb.interceptor = interceptor
	return spb
}

func (spb *StringParserBuilder) WithTransitionObserver(observer state_machine.TransitionObserver[string, string]) *StringParserBuilder {
	spb.observers = append(spb.observers, observer)
	return spb
}

func (spb *StringParserBuilder) Build() *StringParser {
	builder := state_machine.NewStateMachineBuilder[string, string]().
		WithStateMachineDefinition(spb.grammar.toStateMachineDefinition()).
		WithTransitionInterceptor(spb.interceptor)

	for _, observer := range spb.observers {
		builder = builder.WithTransitionObserver(observer)
	}
	return &StringParser{
		stateMachineStart: builder.Build(),
		scanner:           spb.scanner,
	}
}

func NewStringParserBuilder() *StringParserBuilder {
	return &StringParserBuilder{
		scanner: string_scanner.NewSimpleScanner(), // defaults to char by char scanner
	}
}
