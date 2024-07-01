package string_parser

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/parser/state_machine"
	"github.com/openshift-online/ocm-common/pkg/utils/parser/string_scanner"
)

type StringParser struct {
	stateMachineStart *state_machine.State[string, string]
	scanner           string_scanner.Scanner
}

func (p *StringParser) Parse(sql string) error {
	p.reset()

	state := p.stateMachineStart

	scanner := p.scanner
	scanner.Init(sql)

	for scanner.Next() {
		if next, err := state.Move(scanner.Token().Value); err != nil {
			return fmt.Errorf("[%d] error parsing the filter: %v", scanner.Token().Position+1, err)
		} else {
			state = next
		}
	}

	if !state.Eof() {
		return fmt.Errorf(`EOF encountered while parsing string`)
	}

	return nil
}

func (p *StringParser) reset() {
}
