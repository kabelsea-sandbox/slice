package state

import (
	"github.com/pkg/errors"
)

var (
	// ErrTransitionNotAllowed
	ErrTransitionNotAllowed = errors.New("transition not allowed")
)

// TransitionFunc
type TransitionFunc func(toState interface{}) error

// New create new state machine
func New(storage Object, declarations []Transition, fn TransitionFunc) *Machine {
	return &Machine{
		fn:          fn,
		transitions: declarations,
		object:      storage,
	}
}

// Transition
type Transition struct {
	Name string
	From interface{}
	To   interface{}
}

// Machine
type Machine struct {
	object      Object
	fn          TransitionFunc
	transitions []Transition
}

// ToState
func (m *Machine) ToState(new interface{}) (err error) {
	current := m.object.CurrentState()
	allowed := false
	for _, v := range m.transitions {
		if v.From == current && v.To == new {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.Wrapf(ErrTransitionNotAllowed, "%s -> %s", current, new)
	}
	if err = m.fn(new); err != nil {
		return err
	}
	m.object.SetCurrentState(new)
	return nil
}
