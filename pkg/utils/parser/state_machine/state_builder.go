package state_machine

var _ StateBuilder[string, string] = &stateBuilder[string, string]{}

// StateBuilder - builder of State objects
type StateBuilder[T any, U any] interface {
	Data(stateData T) StateBuilder[T, U]
	Accept(acceptor func(value U) bool) StateBuilder[T, U]
	WithTransitionInterceptor(handler TransitionInterceptor[T, U]) StateBuilder[T, U]
	WithTransitionObserver(observers ...TransitionObserver[T, U]) StateBuilder[T, U]
	Build() *State[T, U]
}

type stateBuilder[T any, U any] struct {
	s *State[T, U]
}

func (sb *stateBuilder[T, U]) Data(stateData T) StateBuilder[T, U] {
	sb.s.stateData = stateData
	return sb
}

func (sb *stateBuilder[T, U]) Accept(acceptor func(value U) bool) StateBuilder[T, U] {
	sb.s.accept = acceptor
	return sb
}

func (sb *stateBuilder[T, U]) WithTransitionInterceptor(handler TransitionInterceptor[T, U]) StateBuilder[T, U] {
	sb.s.onStateTransition = handler
	return sb
}

func (sb *stateBuilder[T, U]) WithTransitionObserver(observers ...TransitionObserver[T, U]) StateBuilder[T, U] {
	if observers != nil {
		sb.s.observers = append(sb.s.observers, observers...)
	}

	return sb
}

func (sb *stateBuilder[T, U]) Build() *State[T, U] {
	return sb.s
}

func NewStateBuilder[T any, U any](stateName string) StateBuilder[T, U] {
	return &stateBuilder[T, U]{
		s: &State[T, U]{
			last:      false,
			stateName: stateName,
		},
	}
}
