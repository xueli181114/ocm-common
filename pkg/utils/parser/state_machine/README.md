## The State Machine

This package provides a lightweight, generic state machine framework for enforcing well-defined state
transitions during program execution. The framework is designed for broad applicability, including:

* Flow validation: Ensuring valid sequences of operations in a process.
* Language validation (e.g., string syntax): Verifying adherence to grammar rules.

### Configuring the state machine

The State machine configuration centers around defining a set of states and their valid transitions. 
Each state holds an acceptor object responsible for validating the current data against the transition's criteria. 
This allows for user-defined, data-driven state transitions, even for states with multiple outgoing paths.

Let's start with an example: we want have an API that manages our ticketing system and, based on the inputs it receives, 
it move the ticket to the next status.

Our ticket workflow is very simple:
```
NEW => ASSIGNED => PROGRESSING => REVIEW => COMPLETED  
                        ^           |  
                        |===========|  
```
The ticket starts in status NEW, then moves to ASSIGNED, then to PROGRESSING. After the developer has finished, 
the ticket moves to REVIEW. Here, based on the received input, it can move to COMPLETED or back to PROGRESSING.

The transitions will be as follows:

* -> NEW: from start state to 'NEW' when the 'NEW' input is received
* NEW->ASSIGNED: to ASSIGNED when receive the 'ASSIGN' input
* ASSIGNED->PROGRESSING: to PROGRESSING when receive the 'PROGRESS' input
* PROGRESSING->REVIEW: to REVIEW when receive the 'REVIEW' input
* REVIEW->APPROVED: to complete when receive the 'APPROVED' input
* REVIEW->REJECTED: to progressing when receive the 'REJECTED' input
* REJECTED->END
* APPROVED->END

Let's start with state definition:
```go
definition := state_machine.StateMachineDefinition[string, string]{
    States: []state_machine.StateDefinition[string, string]{
        {Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
        {Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGN")},
        {Name: "PROGRESSING", Acceptor: makeAcceptorForString("PROGRESS")},
        {Name: "REVIEW", Acceptor: makeAcceptorForString("REVIEW")},
        {Name: "COMPLETED", Acceptor: makeAcceptorForString("APPROVED")},
    },
}
```

The `makeAcceptorForString` is an utility function that returns an acceptor that recognize a simple string.
Complex acceptors can be implemented, for example to accept regular expressions.

The one in the example is implemented as follows:
```go
func makeAcceptorForString(s string) func(value string) bool {
	return func(value string) bool {
		return value == s
	}
}
```
Now that we have defined the statuses, we need to define the transitions:
```go
definition := state_machine.StateMachineDefinition[string, string]{
    States: []state_machine.StateDefinition[string, string]{
        {Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
        {Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGN")},
        {Name: "PROGRESSING", Acceptor: makeAcceptorForString("PROGRESS")},
        {Name: "REVIEW", Acceptor: makeAcceptorForString("REVIEW")},
        {Name: "COMPLETED", Acceptor: makeAcceptorForString("APPROVED")},
    },
    Transitions: []state_machine.TransitionDefinition{
        {StateName: state_machine.StartState, ValidTransitions: []string{"NEW"}},
        {StateName: "NEW", ValidTransitions: []string{"ASSIGNED"}},
        {StateName: "ASSIGNED", ValidTransitions: []string{"PROGRESSING"}},
        {StateName: "PROGRESSING", ValidTransitions: []string{"REVIEW"}},
        {StateName: "REVIEW", ValidTransitions: []string{"COMPLETED", "PROGRESSING"}},
        {StateName: "COMPLETED", ValidTransitions: []string{state_machine.EndState}},
    },
}
```

Now that we have both the definitions and the transitions, we can get the state machine:
```go
stateMachine := return NewStateMachineBuilder[string, string]().
    WithStateMachineDefinition(&definition).
    Build()
```

Let's try to parse some flow:

```go
receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "APPROVED"}
currentState := stateMachine
var err error
for _, s := range receivedEvents {
    currentState, err = currentState.Move(s)
    if err != nil {
        fmt.Printf("%v\n", err)
        return
    }
}

if currentState.Eof() {
    fmt.Println("Ticket completed")
} else {
    fmt.Println("Ticket still needs some work")
}
```
The output of this code (Example 1) will be 'Ticket completed'

Let's try with an invalid flow (Example 2). 
```go
receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "ASSIGN"}
...
```

Finally, let's try a valid flow with some recursion.
```go
receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW"}
...
```
This time we won't get any error, but our output will be "Ticket still needs some work" (Example 3)

This time the output will be 'unexpected token `ASSIGN`'.

If needed, you can observe the status transitions by registering an observer:

```go
stateMachine := state_machine.NewStateMachineBuilder[string, string]().
	WithStateMachineDefinition(&definition).
	WithTransitionObserver(func(from, to *state_machine.State[string, string], value string) {
		fmt.Printf("%s => %s (received value: %s) \n", from.Name(), to.Name(), value)
	}).
	Build()
```
With this change, the output will be (Example 4):
```
NEW => ASSIGNED (received value: ASSIGN) 
ASSIGNED => PROGRESSING (received value: PROGRESS) 
PROGRESSING => REVIEW (received value: REVIEW) 
REVIEW => PROGRESSING (received value: PROGRESS) 
PROGRESSING => REVIEW (received value: REVIEW) 
REVIEW => PROGRESSING (received value: PROGRESS) 
PROGRESSING => REVIEW (received value: REVIEW) 
Ticket still needs some work
```

An interceptor can be used to abort the parsing if some event occurs (Example 5):

```go
reviewCount := 0

stateMachine := state_machine.NewStateMachineBuilder[string, string]().
    WithStateMachineDefinition(&definition).
    WithTransitionObserver(func(from, to *state_machine.State[string, string], value string) {
        fmt.Printf("%s => %s (received value: %s) \n", from.Name(), to.Name(), value)
    }).
    WithTransitionInterceptor(func(from, to *state_machine.State[string, string], value string) error {
        if to.Name() == "REVIEW" {
            reviewCount++
            if reviewCount > 2 {
                return fmt.Errorf("Too many reviews. Aborting.")
            }
        }

        return nil
    }).
    Build()
```

This time the parsing is aborted if the statemachine passes through the REVIEW state more than 2 times. The output will be:
```
NEW => ASSIGNED (received value: ASSIGN) 
ASSIGNED => PROGRESSING (received value: PROGRESS) 
PROGRESSING => REVIEW (received value: REVIEW) 
REVIEW => PROGRESSING (received value: PROGRESS) 
PROGRESSING => REVIEW (received value: REVIEW) 
REVIEW => PROGRESSING (received value: PROGRESS) 
Too many reviews. Aborting.
```
## Examples
### Example 1
```go
package main

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/state_machine"
)

func makeAcceptorForString(s string) func(value string) bool {
	return func(value string) bool {
		return value == s
	}
}

func main() {
	definition := state_machine.StateMachineDefinition[string, string]{
		States: []state_machine.StateDefinition[string, string]{
			{Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
			{Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGN")},
			{Name: "PROGRESSING", Acceptor: makeAcceptorForString("PROGRESS")},
			{Name: "REVIEW", Acceptor: makeAcceptorForString("REVIEW")},
			{Name: "COMPLETED", Acceptor: makeAcceptorForString("APPROVED")},
		},
		Transitions: []state_machine.TransitionDefinition{
			{StateName: state_machine.StartState, ValidTransitions: []string{"NEW"}},
			{StateName: "NEW", ValidTransitions: []string{"ASSIGNED"}},
			{StateName: "ASSIGNED", ValidTransitions: []string{"PROGRESSING"}},
			{StateName: "PROGRESSING", ValidTransitions: []string{"REVIEW"}},
			{StateName: "REVIEW", ValidTransitions: []string{"COMPLETED", "PROGRESSING"}},
			{StateName: "COMPLETED", ValidTransitions: []string{state_machine.EndState}},
		},
	}

	stateMachine := state_machine.NewStateMachineBuilder[string, string]().
		WithStateMachineDefinition(&definition).
		Build()

	receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "APPROVED"}
	currentState := stateMachine
	var err error
	for _, s := range receivedEvents {
		currentState, err = currentState.Move(s)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	if currentState.Eof() {
		fmt.Println("Ticket completed")
	} else {
		fmt.Println("Ticket still needs some work")
	}
}
```
### Example 2
```go
package main

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/state_machine"
)

func makeAcceptorForString(s string) func(value string) bool {
	return func(value string) bool {
		return value == s
	}
}

func main() {
	definition := state_machine.StateMachineDefinition[string, string]{
		States: []state_machine.StateDefinition[string, string]{
			{Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
			{Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGN")},
			{Name: "PROGRESSING", Acceptor: makeAcceptorForString("PROGRESS")},
			{Name: "REVIEW", Acceptor: makeAcceptorForString("REVIEW")},
			{Name: "COMPLETED", Acceptor: makeAcceptorForString("APPROVED")},
		},
		Transitions: []state_machine.TransitionDefinition{
			{StateName: state_machine.StartState, ValidTransitions: []string{"NEW"}},
			{StateName: "NEW", ValidTransitions: []string{"ASSIGNED"}},
			{StateName: "ASSIGNED", ValidTransitions: []string{"PROGRESSING"}},
			{StateName: "PROGRESSING", ValidTransitions: []string{"REVIEW"}},
			{StateName: "REVIEW", ValidTransitions: []string{"COMPLETED", "PROGRESSING"}},
			{StateName: "COMPLETED", ValidTransitions: []string{state_machine.EndState}},
		},
	}

	stateMachine := state_machine.NewStateMachineBuilder[string, string]().
		WithStateMachineDefinition(&definition).
		Build()

	receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "ASSIGN"}
	currentState := stateMachine
	var err error
	for _, s := range receivedEvents {
		currentState, err = currentState.Move(s)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	if currentState.Eof() {
		fmt.Println("Ticket completed")
	} else {
		fmt.Println("Ticket still needs some work")
	}
}
```
### Example 3
```go
package main

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/state_machine"
)

func makeAcceptorForString(s string) func(value string) bool {
	return func(value string) bool {
		return value == s
	}
}

func main() {
	definition := state_machine.StateMachineDefinition[string, string]{
		States: []state_machine.StateDefinition[string, string]{
			{Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
			{Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGN")},
			{Name: "PROGRESSING", Acceptor: makeAcceptorForString("PROGRESS")},
			{Name: "REVIEW", Acceptor: makeAcceptorForString("REVIEW")},
			{Name: "COMPLETED", Acceptor: makeAcceptorForString("APPROVED")},
		},
		Transitions: []state_machine.TransitionDefinition{
			{StateName: state_machine.StartState, ValidTransitions: []string{"NEW"}},
			{StateName: "NEW", ValidTransitions: []string{"ASSIGNED"}},
			{StateName: "ASSIGNED", ValidTransitions: []string{"PROGRESSING"}},
			{StateName: "PROGRESSING", ValidTransitions: []string{"REVIEW"}},
			{StateName: "REVIEW", ValidTransitions: []string{"COMPLETED", "PROGRESSING"}},
			{StateName: "COMPLETED", ValidTransitions: []string{state_machine.EndState}},
		},
	}

	stateMachine := state_machine.NewStateMachineBuilder[string, string]().
		WithStateMachineDefinition(&definition).
		Build()

	receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW"}
	currentState := stateMachine
	var err error
	for _, s := range receivedEvents {
		currentState, err = currentState.Move(s)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	if currentState.Eof() {
		fmt.Println("Ticket completed")
	} else {
		fmt.Println("Ticket still needs some work")
	}
}
```
### Example 4
```go
package main

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/state_machine"
)

func makeAcceptorForString(s string) func(value string) bool {
	return func(value string) bool {
		return value == s
	}
}

func main() {
	definition := state_machine.StateMachineDefinition[string, string]{
		States: []state_machine.StateDefinition[string, string]{
			{Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
			{Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGN")},
			{Name: "PROGRESSING", Acceptor: makeAcceptorForString("PROGRESS")},
			{Name: "REVIEW", Acceptor: makeAcceptorForString("REVIEW")},
			{Name: "COMPLETED", Acceptor: makeAcceptorForString("APPROVED")},
		},
		Transitions: []state_machine.TransitionDefinition{
			{StateName: state_machine.StartState, ValidTransitions: []string{"NEW"}},
			{StateName: "NEW", ValidTransitions: []string{"ASSIGNED"}},
			{StateName: "ASSIGNED", ValidTransitions: []string{"PROGRESSING"}},
			{StateName: "PROGRESSING", ValidTransitions: []string{"REVIEW"}},
			{StateName: "REVIEW", ValidTransitions: []string{"COMPLETED", "PROGRESSING"}},
			{StateName: "COMPLETED", ValidTransitions: []string{state_machine.EndState}},
		},
	}

	stateMachine := state_machine.NewStateMachineBuilder[string, string]().
		WithStateMachineDefinition(&definition).
		WithTransitionObserver(func(from, to *state_machine.State[string, string], value string) {
			fmt.Printf("%s => %s (received value: %s) \n", from.Name(), to.Name(), value)
		}).
		Build()

	receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW"}
	currentState := stateMachine
	var err error
	for _, s := range receivedEvents {
		currentState, err = currentState.Move(s)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	if currentState.Eof() {
		fmt.Println("Ticket completed")
	} else {
		fmt.Println("Ticket still needs some work")
	}
}
```
### Example 5
```go
package main

import (
	"fmt"
	"github.com/openshift-online/ocm-common/pkg/utils/state_machine"
)

func makeAcceptorForString(s string) func(value string) bool {
	return func(value string) bool {
		return value == s
	}
}

func main() {
	definition := state_machine.StateMachineDefinition[string, string]{
		States: []state_machine.StateDefinition[string, string]{
			{Name: "NEW", Acceptor: makeAcceptorForString("NEW")},
			{Name: "ASSIGNED", Acceptor: makeAcceptorForString("ASSIGN")},
			{Name: "PROGRESSING", Acceptor: makeAcceptorForString("PROGRESS")},
			{Name: "REVIEW", Acceptor: makeAcceptorForString("REVIEW")},
			{Name: "COMPLETED", Acceptor: makeAcceptorForString("APPROVED")},
		},
		Transitions: []state_machine.TransitionDefinition{
			{StateName: state_machine.StartState, ValidTransitions: []string{"NEW"}},
			{StateName: "NEW", ValidTransitions: []string{"ASSIGNED"}},
			{StateName: "ASSIGNED", ValidTransitions: []string{"PROGRESSING"}},
			{StateName: "PROGRESSING", ValidTransitions: []string{"REVIEW"}},
			{StateName: "REVIEW", ValidTransitions: []string{"COMPLETED", "PROGRESSING"}},
			{StateName: "COMPLETED", ValidTransitions: []string{state_machine.EndState}},
		},
	}

	reviewCount := 0

	stateMachine := state_machine.NewStateMachineBuilder[string, string]().
		WithStateMachineDefinition(&definition).
		WithTransitionObserver(func(from, to *state_machine.State[string, string], value string) {
			fmt.Printf("%s => %s (received value: %s) \n", from.Name(), to.Name(), value)
		}).
		WithTransitionInterceptor(func(from, to *state_machine.State[string, string], value string) error {
			if to.Name() == "REVIEW" {
				reviewCount++
				if reviewCount > 2 {
					return fmt.Errorf("Too many reviews. Aborting.")
				}
			}

			return nil
		}).
		Build()

	receivedEvents := []string{"NEW", "ASSIGN", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW", "PROGRESS", "REVIEW"}
	currentState := stateMachine
	var err error
	for _, s := range receivedEvents {
		currentState, err = currentState.Move(s)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	if currentState.Eof() {
		fmt.Println("Ticket completed")
	} else {
		fmt.Println("Ticket still needs some work")
	}
}
```
