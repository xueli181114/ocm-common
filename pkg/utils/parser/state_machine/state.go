package state_machine

import "fmt"

// State represent a single state of the state machine
// T - the type of the data attached to the state. This can be anything and is useful only to the user
// U - the type of the value that the state will receive to move to the next state
type State[T any, U any] struct {
	stateName string
	stateData T
	accept    Acceptor[U]

	last  bool
	isEof bool
	next  []*State[T, U]

	onStateTransition TransitionInterceptor[T, U]
	observers         []TransitionObserver[T, U]
}

type Acceptor[U any] func(value U) bool
type TransitionInterceptor[T any, U any] func(from, to *State[T, U], value U) error
type TransitionObserver[T any, U any] func(from, to *State[T, U], value U)

func (s *State[T, U]) Name() string {
	return s.stateName
}

// Data returns the custom data associated with this state
func (s *State[T, U]) Data() T {
	return s.stateData
}

func (s *State[T, U]) Move(value U) (*State[T, U], error) {
	for _, next := range s.next {
		if next.accept(value) {
			// valid Value
			if next.onStateTransition != nil {
				if err := next.onStateTransition(s, next, value); err != nil {
					return nil, err
				}
			}

			for _, observer := range s.observers {
				observer(s, next, value)
			}
			return next, nil
		}
	}

	return nil, fmt.Errorf("unexpected token `%v`", value)
}

func (s *State[T, U]) Eof() bool {
	// EOF has been reached. Check if the current Token can be the last one
	return s.last
}

func (s *State[T, U]) addNextState(next *State[T, U]) {
	n := next
	if n.isEof {
		// if the passed in next State is an Eof State, means this is a valid 'last' State
		// Just save the info and discard the 'next' State
		s.last = true
	} else {
		s.next = append(s.next, next)
	}
}

func newStartState[T any, U any]() *State[T, U] {
	return &State[T, U]{
		stateName: "START",
		accept:    func(value U) bool { return false },
	}
}

func newEndState[T any, U any]() *State[T, U] {
	return &State[T, U]{
		stateName: "END",
		isEof:     true,
	}
}
