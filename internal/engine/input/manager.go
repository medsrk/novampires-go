package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

// InputID represents any type of input (keyboard, gamepad, etc)
type InputID interface {
	// This is just a marker interface
	isInputID()
}

// KeyboardKey wraps ebiten.Key
type KeyboardKey struct {
	Key ebiten.Key
}

func (k KeyboardKey) isInputID() {}

// GamepadButton wraps both the gamepad ID and button
type GamepadButton struct {
	GamepadID ebiten.GamepadID
	Button    ebiten.StandardGamepadButton
}

func (g GamepadButton) isInputID() {}

type GamepadAxis struct {
	GamepadID ebiten.GamepadID
	Axis      ebiten.StandardGamepadAxis
}

func (g GamepadAxis) isInputID() {}

// Manager handles mapping between physical inputs and game actions
type Manager struct {
	bindings   map[InputID]Action
	axisValues map[GamepadAxis]float64 // cache axis values
}

func New() *Manager {
	m := &Manager{
		bindings:   make(map[InputID]Action),
		axisValues: make(map[GamepadAxis]float64),
	}

	// Set default keyboard bindings
	defaultKeys := map[ebiten.Key]Action{
		ebiten.KeyW:      ActionMoveUp,
		ebiten.KeyS:      ActionMoveDown,
		ebiten.KeyA:      ActionMoveLeft,
		ebiten.KeyD:      ActionMoveRight,
		ebiten.KeySpace:  ActionAttack,
		ebiten.Key1:      ActionUseAbility1,
		ebiten.Key2:      ActionUseAbility2,
		ebiten.Key3:      ActionUseAbility3,
		ebiten.KeyEscape: ActionPause,
		ebiten.KeyUp:     ActionMoveUp,
		ebiten.KeyDown:   ActionMoveDown,
		ebiten.KeyLeft:   ActionMoveLeft,
		ebiten.KeyRight:  ActionMoveRight,
		ebiten.KeyF1:     ActionToggleDebug,
	}

	defaultGamepadButtons := map[ebiten.StandardGamepadButton]Action{
		ebiten.StandardGamepadButtonLeftTop:    ActionMoveUp,    // D-pad Up
		ebiten.StandardGamepadButtonLeftRight:  ActionMoveRight, // D-pad Right
		ebiten.StandardGamepadButtonLeftBottom: ActionMoveDown,  // D-pad Down
		ebiten.StandardGamepadButtonLeftLeft:   ActionMoveLeft,  // D-pad Left

		ebiten.StandardGamepadButtonRightBottom: ActionAttack,      // A/Cross
		ebiten.StandardGamepadButtonRightRight:  ActionUseAbility1, // B/Circle
		ebiten.StandardGamepadButtonRightLeft:   ActionUseAbility2, // X/Square
		ebiten.StandardGamepadButtonRightTop:    ActionUseAbility3, // Y/Triangle

		ebiten.StandardGamepadButtonCenterRight: ActionPause,       // Start
		ebiten.StandardGamepadButtonCenterLeft:  ActionToggleDebug, // Select
	}

	for k, v := range defaultKeys {
		m.Bind(KeyboardKey{Key: k}, v)
	}

	for b, v := range defaultGamepadButtons {
		m.Bind(GamepadButton{Button: b}, v)
	}

	return m
}

func (m *Manager) Rebind(oldInput InputID, newInput InputID) {
	action := m.bindings[oldInput]
	m.Unbind(oldInput)
	m.Bind(newInput, action)
}

func (m *Manager) Bind(input InputID, action Action) {
	m.bindings[input] = action
}

func (m *Manager) Unbind(input InputID) {
	delete(m.bindings, input)
}

func (m *Manager) isInputActive(id InputID) bool {
	switch v := id.(type) {
	case KeyboardKey:
		return ebiten.IsKeyPressed(v.Key)
	case GamepadButton:
		return ebiten.IsStandardGamepadButtonPressed(v.GamepadID, v.Button)
	default:
		return false
	}
}

func (m *Manager) isInputJustPressed(id InputID) bool {
	switch v := id.(type) {
	case KeyboardKey:
		return inpututil.IsKeyJustPressed(v.Key)
	case GamepadButton:
		return inpututil.IsStandardGamepadButtonJustPressed(v.GamepadID, v.Button)
	default:
		return false
	}
}

func (m *Manager) isInputJustReleased(id InputID) bool {
	switch v := id.(type) {
	case KeyboardKey:
		return inpututil.IsKeyJustReleased(v.Key)
	case GamepadButton:
		return inpututil.IsStandardGamepadButtonJustReleased(v.GamepadID, v.Button)
	default:
		return false
	}
}

func (m *Manager) IsPressed(action Action) bool {
	for input, a := range m.bindings {
		if a == action && m.isInputActive(input) {
			return true
		}
	}
	return false
}

func (m *Manager) JustPressed(action Action) bool {
	for input, a := range m.bindings {
		if a == action && m.isInputJustPressed(input) {
			return true
		}
	}
	return false
}

func (m *Manager) JustReleased(action Action) bool {
	for input, a := range m.bindings {
		if a == action && m.isInputJustReleased(input) {
			return true
		}
	}
	return false
}

func (m *Manager) GetMovementVector() (float64, float64) {
	dx, dy := 0.0, 0.0

	// Digital input (keyboard/d-pad)
	if m.IsPressed(ActionMoveUp) {
		dy--
	}
	if m.IsPressed(ActionMoveDown) {
		dy++
	}
	if m.IsPressed(ActionMoveLeft) {
		dx--
	}
	if m.IsPressed(ActionMoveRight) {
		dx++
	}

	// If no digital input, check stick
	if dx == 0 && dy == 0 {
		// Get first connected gamepad
		ids := ebiten.AppendGamepadIDs(nil)
		if len(ids) > 0 {
			dx = ebiten.StandardGamepadAxisValue(ids[0], ebiten.StandardGamepadAxisLeftStickHorizontal)
			dy = ebiten.StandardGamepadAxisValue(ids[0], ebiten.StandardGamepadAxisLeftStickVertical)

			// Apply deadzone
			deadzone := 0.2 // Typical value, could be made configurable
			if math.Abs(dx) < deadzone {
				dx = 0
			}
			if math.Abs(dy) < deadzone {
				dy = 0
			}
		}
	}

	// Normalize for diagonal movement
	if dx != 0 && dy != 0 {
		length := math.Sqrt(dx*dx + dy*dy)
		if length > 1 {
			dx /= length
			dy /= length
		}
	}

	return dx, dy
}

func (m *Manager) GetAimVector() (float64, float64) {
	return 0, 0
}

// GetAllBindings returns all input bindings, including gamepad bindings
func (m *Manager) GetAllBindings() map[InputID]Action {
	// Return a copy of all bindings
	allBindings := make(map[InputID]Action)
	for k, v := range m.bindings {
		allBindings[k] = v
	}
	return allBindings
}

func (m *Manager) GetActionState(action Action) ActionState {
	return ActionState{
		Active:       m.IsPressed(action),
		JustPressed:  m.JustPressed(action),
		JustReleased: m.JustReleased(action),
	}
}

func (m *Manager) CreateDebugWindow() *DebugWindow {
	return NewDebugWindow(m)
}
