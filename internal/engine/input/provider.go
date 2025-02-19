package input

type Action uint8

const (
	ActionMoveUp Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionAutoAttack
	ActionUseAbility1
	ActionUseAbility2
	ActionUseAbility3
	ActionPause
	ActionToggleDebug

	ActionInteract
	ActionMenu
)

// ActionState represents the state of an input action
type ActionState struct {
	// Whether the action is currently active
	Active bool

	// Whether the action was just activated this frame
	JustPressed bool

	// Whether the action was just released this frame
	JustReleased bool

	// For analog inputs, value between 0-1
	Value float64
}

func (a *Action) String() string {
	return [...]string{"ActionMoveUp", "ActionMoveDown", "ActionMoveLeft", "ActionMoveRight", "ActionAutoAttack", "ActionUseAbility1", "ActionUseAbility2", "ActionUseAbility3", "ActionPause", "Toggle Debug"}[*a]
}

// InputProvider defines the interface for accessing input state
type InputProvider interface {
	// GetActionState returns the current state of an action
	GetActionState(action Action) ActionState

	// IsPressed returns whether an action is currently active
	IsPressed(action Action) bool

	// JustPressed returns whether an action was just activated this frame
	JustPressed(action Action) bool

	// JustReleased returns whether an action was just released this frame
	JustReleased(action Action) bool

	// GetMovementVector returns the normalized movement vector from input
	GetMovementVector() (float64, float64)

	// GetAimVector returns the normalized aim vector from input
	GetAimVector() (float64, float64)
}
