package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/camera"
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

// ComboKey represents a key combination (e.g., Ctrl+P)
type ComboKey struct {
	Modifier ebiten.Key
	Key      ebiten.Key
}

func (c ComboKey) isInputID() {}

// Config holds all configurable input parameters
type Config struct {
	Deadzone float64 // Deadzone for analog sticks
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Deadzone: 0.2,
	}
}

// Manager handles mapping between physical inputs and game actions
type Manager struct {
	bindings     map[InputID]common.Action
	axisValues   map[GamepadAxis]float64
	playerPos    common.Vector2
	usingGamepad bool
	lastMouseX   int
	lastMouseY   int
	config       *Config
	camera       *camera.Camera
}

// New creates a new input manager with default bindings
func New() *Manager { return NewWithConfig(DefaultConfig()) }

func NewWithConfig(config *Config) *Manager {
	m := &Manager{
		bindings:   make(map[InputID]common.Action),
		axisValues: make(map[GamepadAxis]float64),
		config:     config,
	}

	m.setupDefaultBindings()
	return m
}

func (m *Manager) setupDefaultBindings() {
	// Default keyboard bindings
	defaultKeys := map[ebiten.Key]common.Action{
		ebiten.KeyW:      common.ActionMoveUp,
		ebiten.KeyS:      common.ActionMoveDown,
		ebiten.KeyA:      common.ActionMoveLeft,
		ebiten.KeyD:      common.ActionMoveRight,
		ebiten.KeySpace:  common.ActionAutoAttack,
		ebiten.Key1:      common.ActionUseAbility1,
		ebiten.Key2:      common.ActionUseAbility2,
		ebiten.Key3:      common.ActionUseAbility3,
		ebiten.KeyEscape: common.ActionPause,
		ebiten.KeyUp:     common.ActionMoveUp,
		ebiten.KeyDown:   common.ActionMoveDown,
		ebiten.KeyLeft:   common.ActionMoveLeft,
		ebiten.KeyRight:  common.ActionMoveRight,
		ebiten.KeyF1:     common.ActionToggleDebug,
	}

	defaultGamepadButtons := map[ebiten.StandardGamepadButton]common.Action{
		ebiten.StandardGamepadButtonLeftTop:    common.ActionMoveUp,
		ebiten.StandardGamepadButtonLeftRight:  common.ActionMoveRight,
		ebiten.StandardGamepadButtonLeftBottom: common.ActionMoveDown,
		ebiten.StandardGamepadButtonLeftLeft:   common.ActionMoveLeft,

		ebiten.StandardGamepadButtonRightBottom: common.ActionAutoAttack,
		ebiten.StandardGamepadButtonRightRight:  common.ActionUseAbility1,
		ebiten.StandardGamepadButtonRightLeft:   common.ActionUseAbility2,
		ebiten.StandardGamepadButtonRightTop:    common.ActionUseAbility3,

		ebiten.StandardGamepadButtonCenterRight: common.ActionPause,
		ebiten.StandardGamepadButtonCenterLeft:  common.ActionToggleDebug,
	}

	defaultDebugBindings := map[ComboKey]common.Action{
		{Modifier: ebiten.KeyControl, Key: ebiten.KeyP}: common.ActionTogglePlayerDebug,
		{Modifier: ebiten.KeyControl, Key: ebiten.KeyI}: common.ActionToggleInputDebug,
		{Modifier: ebiten.KeyControl, Key: ebiten.KeyB}: common.ActionToggleBindingEditor,
	}

	for k, v := range defaultKeys {
		m.Bind(KeyboardKey{Key: k}, v)
	}

	for b, v := range defaultGamepadButtons {
		m.Bind(GamepadButton{Button: b}, v)
	}

	for c, v := range defaultDebugBindings {
		m.Bind(c, v)
	}
}

func (m *Manager) updateGamepadState() {
	ids := ebiten.AppendGamepadIDs(nil)
	wasUsingGamepad := m.usingGamepad

	if len(ids) > 0 {
		// Check for any gamepad input
		for _, id := range ids {
			// Check buttons
			for b := ebiten.StandardGamepadButtonLeftTop; b <= ebiten.StandardGamepadButtonMax; b++ {
				if ebiten.IsStandardGamepadButtonPressed(id, b) {
					m.usingGamepad = true
					return
				}
			}

			// Check sticks
			dx := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
			dy := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
			if math.Abs(dx) >= m.config.Deadzone || math.Abs(dy) >= m.config.Deadzone {
				m.usingGamepad = true
				return
			}
		}
	}

	// Check if mouse moved
	x, y := ebiten.CursorPosition()
	if x != m.lastMouseX || y != m.lastMouseY {
		m.usingGamepad = false
		m.lastMouseX = x
		m.lastMouseY = y
	}

	// Check keyboard input
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			m.usingGamepad = false
			return
		}
	}

	// Preserve previous state if no input detected
	m.usingGamepad = wasUsingGamepad
}

func (m *Manager) Update() error {
	m.updateGamepadState()
	return nil
}

func (m *Manager) Rebind(oldInput InputID, newInput InputID) {
	binding := m.bindings[oldInput]
	m.Unbind(oldInput)
	m.Bind(newInput, binding)
}

func (m *Manager) Bind(input InputID, binding common.Action) {
	m.bindings[input] = binding
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
	case ComboKey:
		return ebiten.IsKeyPressed(v.Modifier) && ebiten.IsKeyPressed(v.Key)
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
	case ComboKey:
		// For combo keys, detect just pressed when either key is just pressed while the other is held
		return (ebiten.IsKeyPressed(v.Modifier) && inpututil.IsKeyJustPressed(v.Key)) ||
			(ebiten.IsKeyPressed(v.Key) && inpututil.IsKeyJustPressed(v.Modifier))
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
	case ComboKey:
		// Both keys must have been pressed in the previous frame
		wasPressedPreviousFrame := inpututil.IsKeyJustReleased(v.Key) || inpututil.IsKeyJustReleased(v.Modifier)
		wasComboActive := inpututil.IsKeyJustReleased(v.Key) && ebiten.IsKeyPressed(v.Modifier) ||
			inpututil.IsKeyJustReleased(v.Modifier) && ebiten.IsKeyPressed(v.Key)
		return wasPressedPreviousFrame && wasComboActive
	default:
		return false
	}
}

func (m *Manager) IsPressed(action common.Action) bool {
	for input, a := range m.bindings {
		if a == action && m.isInputActive(input) {
			return true
		}
	}
	return false
}

func (m *Manager) JustPressed(action common.Action) bool {
	for input, a := range m.bindings {
		if a == action && m.isInputJustPressed(input) {
			return true
		}
	}
	return false
}

func (m *Manager) JustReleased(action common.Action) bool {
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
	if m.IsPressed(common.ActionMoveUp) {
		dy--
	}
	if m.IsPressed(common.ActionMoveDown) {
		dy++
	}
	if m.IsPressed(common.ActionMoveLeft) {
		dx--
	}
	if m.IsPressed(common.ActionMoveRight) {
		dx++
	}

	// If using digital input, normalize to get full magnitude
	if dx != 0 || dy != 0 {
		return normalizeVector(dx, dy)
	}

	// Check analog stick
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) > 0 {
		dx = ebiten.StandardGamepadAxisValue(ids[0], ebiten.StandardGamepadAxisLeftStickHorizontal)
		dy = ebiten.StandardGamepadAxisValue(ids[0], ebiten.StandardGamepadAxisLeftStickVertical)

		// Apply deadzone with smooth transition
		magnitude := math.Sqrt(dx*dx + dy*dy)
		if magnitude < m.config.Deadzone {
			return 0, 0
		}

		// Smooth out the deadzone transition
		adjustedMagnitude := (magnitude - m.config.Deadzone) / (1 - m.config.Deadzone)
		dx = dx / magnitude * adjustedMagnitude
		dy = dy / magnitude * adjustedMagnitude
	}

	return dx, dy
}

func (m *Manager) GetMousePositionWorld() (int, int) {
	if m.camera == nil {
		return m.GetMousePosition()
	}

	// Get screen mouse position
	screenX, screenY := m.GetMousePosition()
	screenPos := common.Vector2{X: float64(screenX), Y: float64(screenY)}

	// Convert to world position using camera
	worldPos := m.camera.ScreenToWorld(screenPos)

	return int(worldPos.X), int(worldPos.Y)
}

func (m *Manager) GetMousePosition() (int, int) {
	return ebiten.CursorPosition()
}

func (m *Manager) SetCamera(camera *camera.Camera) {
	m.camera = camera
}

func (m *Manager) GetGamepadAim() (float64, float64, bool) {
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) > 0 {
		dx := ebiten.StandardGamepadAxisValue(ids[0], ebiten.StandardGamepadAxisRightStickHorizontal)
		dy := ebiten.StandardGamepadAxisValue(ids[0], ebiten.StandardGamepadAxisRightStickVertical)

		if math.Abs(dx) >= m.config.Deadzone || math.Abs(dy) >= m.config.Deadzone {
			return dx, dy, true
		}
	}
	return 0, 0, false
}

// GetAllBindings returns all input bindings, including gamepad bindings
func (m *Manager) GetAllBindings() map[InputID]common.Action {
	// Return a copy of all bindings
	allBindings := make(map[InputID]common.Action)
	for k, v := range m.bindings {
		allBindings[k] = v
	}
	return allBindings
}

func (m *Manager) GetActionState(action common.Action) common.ActionState {
	return common.ActionState{
		Active:       m.IsPressed(action),
		JustPressed:  m.JustPressed(action),
		JustReleased: m.JustReleased(action),
	}
}

// IsUsingGamepad returns whether the last input came from a gamepad
func (m *Manager) IsUsingGamepad() bool {
	return m.usingGamepad
}

// normalizeVector normalizes a 2D vector if its length is greater than 1
func normalizeVector(x, y float64) (float64, float64) {
	if x != 0 && y != 0 {
		length := math.Sqrt(x*x + y*y)
		if length > 1 {
			return x / length, y / length
		}
	}
	return x, y
}

func (m *Manager) GetCurrentKey() string {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			return k.String()
		}
	}
	return ""
}

func (m *Manager) CreateDebugWindow() *DebugWindow {
	return NewDebugWindow(m)
}
