package state_machine

const (
	StartState = "__$$_START_$$__"
	EndState   = "__$$_END_$$__"
)

type StateDefinition[T any, U any] struct {
	Name        string
	StateData   T
	Acceptor    func(value U) bool
	OnIntercept func(data T, value U) error
}

type StateMachineDefinition[T any, U any] struct {
	States      []StateDefinition[T, U]
	Transitions []TransitionDefinition
}

type TransitionDefinition struct {
	StateName        string
	ValidTransitions []string
}

type StateMachineBuilder[T any, U any] interface {
	WithTransitionInterceptor(handler TransitionInterceptor[T, U]) StateMachineBuilder[T, U]
	WithTransitionObserver(observer TransitionObserver[T, U]) StateMachineBuilder[T, U]
	Build() *State[T, U]
}

type StateMachineConfigurator[T any, U any] interface {
	WithStateMachineDefinition(definition *StateMachineDefinition[T, U]) StateMachineBuilder[T, U]
}

type stateMachineBuilder[T any, U any] struct {
	definition            *StateMachineDefinition[T, U]
	transitionInterceptor TransitionInterceptor[T, U]
	observers             []TransitionObserver[T, U]
}

func (smb *stateMachineBuilder[T, U]) WithStateMachineDefinition(definition *StateMachineDefinition[T, U]) StateMachineBuilder[T, U] {
	smb.definition = definition
	return smb
}

func (smb *stateMachineBuilder[T, U]) WithTransitionInterceptor(handler TransitionInterceptor[T, U]) StateMachineBuilder[T, U] {
	smb.transitionInterceptor = handler
	return smb
}

func (smb *stateMachineBuilder[T, U]) WithTransitionObserver(observer TransitionObserver[T, U]) StateMachineBuilder[T, U] {
	smb.observers = append(smb.observers, observer)
	return smb
}

func (smb *stateMachineBuilder[T, U]) Build() *State[T, U] {
	stateMap := make(map[string]*State[T, U])

	stateMap[StartState] = newStartState[T, U]()
	stateMap[EndState] = newEndState[T, U]()

	// build all the tokens
	for _, t := range smb.definition.States {
		stateMap[t.Name] = NewStateBuilder[T, U](t.Name).
			Data(t.StateData).
			Accept(t.Acceptor).
			WithTransitionInterceptor(smb.transitionInterceptor).
			WithTransitionObserver(smb.observers...).
			Build()
	}

	// add all the transitions
	for _, transition := range smb.definition.Transitions {
		currentState := stateMap[transition.StateName]
		for _, targetStateName := range transition.ValidTransitions {
			targetState := stateMap[targetStateName]
			currentState.addNextState(targetState)
		}
	}

	return stateMap[StartState]
}

func NewStateMachineBuilder[T any, U any]() StateMachineConfigurator[T, U] {
	return &stateMachineBuilder[T, U]{}
}
