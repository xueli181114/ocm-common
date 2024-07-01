package string_parser

import "github.com/openshift-online/ocm-common/pkg/utils/parser/state_machine"

type TokenDefinition = state_machine.StateDefinition[string, string]

// Grammar - When parsing strings, it is easier to speak about a grammar and tokens than dealing with 'states'. These objects are just
// aliases of the StateMachineDefinition and TransitionDefinition objects
type Grammar struct {
	Tokens      []TokenDefinition
	Transitions []TokenTransitions
}

type TokenTransitions struct {
	TokenName        string
	ValidTransitions []string
}

func (g *Grammar) toStateMachineDefinition() *state_machine.StateMachineDefinition[string, string] {
	ret := state_machine.StateMachineDefinition[string, string]{
		States:      nil,
		Transitions: nil,
	}

	ret.States = append(ret.States, g.Tokens...)

	for _, transition := range g.Transitions {
		ret.Transitions = append(ret.Transitions, *transition.toStateTransitions())
	}

	return &ret
}

func (tt *TokenTransitions) toStateTransitions() *state_machine.TransitionDefinition {
	return &state_machine.TransitionDefinition{
		StateName:        tt.TokenName,
		ValidTransitions: tt.ValidTransitions,
	}
}
