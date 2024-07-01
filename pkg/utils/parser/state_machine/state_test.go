package state_machine

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TokenFamily string

func makeAcceptorForString(s string) func(value string) bool {
	return func(value string) bool {
		return value == s
	}
}

func createTestStateMachine() *State[TokenFamily, string] {
	definition := StateMachineDefinition[TokenFamily, string]{
		States: []StateDefinition[TokenFamily, string]{
			{Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
			{Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGNED")},
			{Name: "IN PROGRESS", Acceptor: makeAcceptorForString("IN PROGRESS")},
			{Name: "WAITING FOR REVIEW", Acceptor: makeAcceptorForString("WAITING FOR REVIEW")},
			{Name: "REVIEWING", Acceptor: makeAcceptorForString("REVIEWING")},
			{Name: "WAITING FOR RELEASE", Acceptor: makeAcceptorForString("WAITING FOR RELEASE")},
			{Name: "WON'T DO", Acceptor: makeAcceptorForString("WON'T DO")},
			{Name: "DONE", Acceptor: makeAcceptorForString("DONE")},
		},
		Transitions: []TransitionDefinition{
			{StateName: StartState, ValidTransitions: []string{"NEW"}},
			{StateName: "NEW", ValidTransitions: []string{"ASSIGNED", "WON'T DO"}},
			{StateName: "ASSIGNED", ValidTransitions: []string{"IN PROGRESS", "WON'T DO"}},
			{StateName: "IN PROGRESS", ValidTransitions: []string{"WAITING FOR REVIEW", "ASSIGNED", "WON'T DO"}},
			{StateName: "WAITING FOR REVIEW", ValidTransitions: []string{"IN PROGRESS", "REVIEWING", "WON'T DO"}},
			{StateName: "REVIEWING", ValidTransitions: []string{"IN PROGRESS", "WON'T DO", "WAITING FOR RELEASE"}},
			{StateName: "WAITING FOR RELEASE", ValidTransitions: []string{"DONE", "WON'T DO", "IN PROGRESS"}},
			{StateName: "DONE", ValidTransitions: []string{EndState}},
			{StateName: "WON'T DO", ValidTransitions: []string{EndState}},
		},
	}

	return NewStateMachineBuilder[TokenFamily, string]().
		WithStateMachineDefinition(&definition).
		Build()
}

var _ = Describe("State Machine Tests", func() {
	DescribeTable("Valid Paths", func(path []string) {
		initialState := createTestStateMachine()
		currentState := initialState
		var err error
		for _, s := range path {
			currentState, err = currentState.Move(s)
			Expect(err).ToNot(HaveOccurred())
		}
		Expect(currentState.Eof()).To(BeTrue())
	},
		Entry("Simple, ends with: WON'T DO", []string{"NEW", "ASSIGNED", "WON'T DO"}),
		Entry("Pass through REVIEW, ends with: WON'T DO", []string{"NEW", "ASSIGNED", "IN PROGRESS", "WAITING FOR REVIEW", "REVIEWING", "WON'T DO"}),
		Entry("From NEW to DONE passing through REVIEW", []string{"NEW", "ASSIGNED", "IN PROGRESS", "WAITING FOR REVIEW", "REVIEWING", "WAITING FOR RELEASE", "DONE"}),
		Entry("From NEW to DONE reiterating IN PROGRESS State", []string{"NEW", "ASSIGNED", "IN PROGRESS", "WAITING FOR REVIEW", "IN PROGRESS", "WAITING FOR REVIEW", "REVIEWING", "WAITING FOR RELEASE", "DONE"}),
	)

	DescribeTable("Invalid Paths", func(path []string, expectedError string) {
		initialState := createTestStateMachine()
		currentState := initialState
		var err error
		for _, s := range path {
			currentState, err = currentState.Move(s)
			if err != nil {
				Expect(err.Error()).To(Equal(expectedError))
			}
		}
		if currentState != nil {
			Expect(currentState.Eof()).To(BeFalse()) // in this test we never end successfully
		}
	},
		Entry("From IN PROGRESS to DONE without REVIEW", []string{"NEW", "ASSIGNED", "IN PROGRESS", "DONE"}, "unexpected token `DONE`"),
		Entry("From REVIEW to ASSIGNED", []string{"NEW", "ASSIGNED", "IN PROGRESS", "WAITING FOR REVIEW", "ASSIGNED"}, "unexpected token `ASSIGNED`"),
		Entry("END STATE NOT REACHED", []string{"NEW", "ASSIGNED", "IN PROGRESS", "WAITING FOR REVIEW"}, ""),
	)
})
