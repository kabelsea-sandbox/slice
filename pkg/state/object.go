package state

// Object
type Object interface {
	CurrentState() interface{}
	SetCurrentState(state interface{})
}
