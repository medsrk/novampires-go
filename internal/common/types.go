package common

type TargetInfo struct {
	ID       uint64
	Pos      Vector2
	Vel      Vector2
	Radius   float64
	Priority float64 // higher is more important
}

type Rectangle struct {
	Pos  Vector2
	Size Vector2
}

func (r Rectangle) Contains(p Vector2) bool {
	return p.X >= r.Pos.X &&
		p.X <= r.Pos.X+r.Size.X &&
		p.Y >= r.Pos.Y &&
		p.Y <= r.Pos.Y+r.Size.Y
}

func (r Rectangle) Intersects(other Rectangle) bool {
	return r.Pos.X < other.Pos.X+other.Size.X &&
		r.Pos.X+r.Size.X > other.Pos.X &&
		r.Pos.Y < other.Pos.Y+other.Size.Y &&
		r.Pos.Y+r.Size.Y > other.Pos.Y
}

func (r Rectangle) Center() Vector2 {
	return Vector2{
		X: r.Pos.X + r.Size.X/2,
		Y: r.Pos.Y + r.Size.Y/2,
	}
}

// Action represents a game action that can be triggered by input
type Action uint8

// Add debug-specific actions
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

	// Debug window specific actions
	ActionTogglePlayerDebug
	ActionToggleInputDebug
	ActionToggleBindingEditor
)

var Actions = []Action{
	ActionMoveUp,
	ActionMoveDown,
	ActionMoveLeft,
	ActionMoveRight,
	ActionAutoAttack,
	ActionUseAbility1,
	ActionUseAbility2,
	ActionUseAbility3,
	ActionPause,
	ActionToggleDebug,
	ActionInteract,
	ActionMenu,

	ActionTogglePlayerDebug,
	ActionToggleInputDebug,
	ActionToggleBindingEditor,
}

func (a Action) String() string {
	switch a {
	case ActionMoveUp:
		return "Move Up"
	case ActionMoveDown:
		return "Move Down"
	case ActionMoveLeft:
		return "Move Left"
	case ActionMoveRight:
		return "Move Right"
	case ActionAutoAttack:
		return "Auto Attack"
	case ActionUseAbility1:
		return "Use Ability 1"
	case ActionUseAbility2:
		return "Use Ability 2"
	case ActionUseAbility3:
		return "Use Ability 3"
	case ActionPause:
		return "Pause"
	case ActionToggleDebug:
		return "Toggle Debug"
	case ActionInteract:
		return "Interact"
	case ActionMenu:
		return "Menu"
	case ActionTogglePlayerDebug:
		return "Toggle Player Debug"
	case ActionToggleInputDebug:
		return "Toggle Input Debug"
	case ActionToggleBindingEditor:
		return "Toggle Binding Editor"
	default:
		return "Unknown Action"
	}
}

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

// debug Window names for the debug manager
const (
	WindowPlayerDebug = "Player Debug"
	WindowInputDebug  = "Input Debug"
	WindowBindingEdit = "Binding Editor"
	WindowCameraDebug = "Camera Debug"
)

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

	// GetMousePositionScreen returns the current mouse position in screen coordinates
	GetMousePositionScreen() (int, int)

	// GetMousePositionWorld returns the current mouse position in world coordinates
	GetMousePositionWorld() (int, int)

	// GetGamepadAim returns the aim vector from the gamepad right stick
	GetGamepadAim() (float64, float64, bool)
}
