package input

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/engine/debug"
	"sort"
)

type KeyBindingEditorWindow struct {
	manager         *Manager
	open            bool
	listening       bool
	selectedAction  Action
	rebindMode      bool
	oldKey          ebiten.Key
	keyBindings     []KeyActionPair
	gamepadBindings []GamepadActionPair
	needsRefresh    bool
	showGamepad     bool
}

// KeyActionPair represents a keyboard binding
type KeyActionPair struct {
	Key    ebiten.Key
	Action Action
}

// GamepadActionPair represents a gamepad binding
type GamepadActionPair struct {
	GamepadID ebiten.GamepadID
	Button    ebiten.StandardGamepadButton
	Action    Action
}

func NewKeyBindingEditorWindow(manager *Manager) *KeyBindingEditorWindow {
	return &KeyBindingEditorWindow{
		manager:         manager,
		open:            true,
		listening:       false,
		selectedAction:  ActionMoveUp,
		rebindMode:      false,
		oldKey:          ebiten.Key(0),
		keyBindings:     nil,
		gamepadBindings: nil,
		needsRefresh:    true,
		showGamepad:     false,
	}
}

// GetKeyBindings returns keyboard bindings
func (m *Manager) GetKeyBindings() map[ebiten.Key]Action {
	bindings := make(map[ebiten.Key]Action)
	for input, action := range m.bindings {
		if keyInput, ok := input.(KeyboardKey); ok {
			bindings[keyInput.Key] = action
		}
	}
	return bindings
}

// GetGamepadBindings returns gamepad bindings
func (m *Manager) GetGamepadBindings() []GamepadActionPair {
	var bindings []GamepadActionPair
	for input, action := range m.bindings {
		if gamepadInput, ok := input.(GamepadButton); ok {
			bindings = append(bindings, GamepadActionPair{
				GamepadID: gamepadInput.GamepadID,
				Button:    gamepadInput.Button,
				Action:    action,
			})
		}
	}
	return bindings
}

func (w *KeyBindingEditorWindow) refreshBindings() {
	if !w.needsRefresh {
		return
	}

	// Get keyboard bindings
	keyBindings := w.manager.GetKeyBindings()
	w.keyBindings = make([]KeyActionPair, 0, len(keyBindings))
	for key, action := range keyBindings {
		w.keyBindings = append(w.keyBindings, KeyActionPair{
			Key:    key,
			Action: action,
		})
	}
	sort.Slice(w.keyBindings, func(i, j int) bool {
		return int(w.keyBindings[i].Action) < int(w.keyBindings[j].Action)
	})

	// Get gamepad bindings
	w.gamepadBindings = w.manager.GetGamepadBindings()
	sort.Slice(w.gamepadBindings, func(i, j int) bool {
		return int(w.gamepadBindings[i].Action) < int(w.gamepadBindings[j].Action)
	})

	w.needsRefresh = false
}

func (w *KeyBindingEditorWindow) Draw() {
	if !w.open {
		return
	}

	w.refreshBindings()
	imgui.SetNextWindowSize(imgui.Vec2{X: 450, Y: 500})

	if imgui.Begin("Input Bindings") {
		// Tab bar for keyboard/gamepad
		if imgui.BeginTabBar("##input_tabs") {
			if imgui.BeginTabItem("Keyboard") {
				w.showGamepad = false
				w.drawKeyboardBindings()
				imgui.EndTabItem()
			}

			if imgui.BeginTabItem("Gamepad") {
				w.showGamepad = true
				w.drawGamepadBindings()
				imgui.EndTabItem()
			}

			imgui.EndTabBar()
		}

		if w.listening && !w.showGamepad {
			w.listenForKeyPress()
		}
	}
	imgui.End()
}

func (w *KeyBindingEditorWindow) drawKeyboardBindings() {
	// Edit existing bindings section
	if imgui.CollapsingHeaderTreeNodeFlagsV("Edit Keyboard Bindings", imgui.TreeNodeFlagsDefaultOpen) {
		for i, binding := range w.keyBindings {
			actionStr := binding.Action.String()
			imgui.PushIDInt(int32(i))

			// Display binding info
			imgui.Text(fmt.Sprintf("%s:", actionStr))
			imgui.SameLine()

			// Display current key
			keyName := binding.Key.String()
			if len(keyName) > 3 && keyName[:3] == "Key" {
				keyName = keyName[3:]
			}

			// Rebinding button
			buttonText := keyName
			if w.listening && w.rebindMode && w.oldKey == binding.Key {
				buttonText = "Press any key..."
				debug.ColoredText(buttonText, debug.ColorJustPressed)
			} else {
				if imgui.Button(buttonText) {
					w.listening = true
					w.rebindMode = true
					w.oldKey = binding.Key
					w.selectedAction = binding.Action
				}
			}

			imgui.PopID()
			imgui.Separator()
		}
	}

	// Add new bindings section
	if imgui.CollapsingHeaderTreeNodeFlagsV("Add New Binding", 0) {
		// Action selector
		imgui.Text("Select action:")
		actions := []Action{
			ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight,
			ActionAttack, ActionUseAbility1, ActionUseAbility2, ActionUseAbility3,
			ActionPause, ActionToggleDebug,
		}

		for i, action := range actions {
			imgui.PushIDInt(int32(1000 + i))
			if imgui.Button(action.String()) {
				w.selectedAction = action
			}
			imgui.PopID()
			if i%3 != 2 && i < len(actions)-1 {
				imgui.SameLine()
			}
		}

		imgui.Separator()
		imgui.Text(fmt.Sprintf("Selected: %s", w.selectedAction.String()))

		// Bind button
		bindText := "Bind New Key"
		if w.listening && !w.rebindMode {
			bindText = "Press a key..."
			debug.ColoredText(bindText, debug.ColorJustPressed)
		} else {
			if imgui.Button(bindText) {
				w.listening = true
				w.rebindMode = false
			}
		}
	}
}

// Helper to get standard gamepad button names
func getGamepadButtonName(button ebiten.StandardGamepadButton) string {
	// W3C standard button indices
	switch button {
	case ebiten.StandardGamepadButtonRightBottom:
		return "Button 0 (Bottom right)"
	case ebiten.StandardGamepadButtonRightRight:
		return "Button 1 (Right right)"
	case ebiten.StandardGamepadButtonRightLeft:
		return "Button 2 (Left right)"
	case ebiten.StandardGamepadButtonRightTop:
		return "Button 3 (Top right)"
	case ebiten.StandardGamepadButtonLeftTop:
		return "Button 12 (Top left)"
	case ebiten.StandardGamepadButtonLeftRight:
		return "Button 15 (Right left)"
	case ebiten.StandardGamepadButtonLeftBottom:
		return "Button 13 (Bottom left)"
	case ebiten.StandardGamepadButtonLeftLeft:
		return "Button 14 (Left left)"
	case ebiten.StandardGamepadButtonCenterLeft:
		return "Button 8 (Left center)"
	case ebiten.StandardGamepadButtonCenterRight:
		return "Button 9 (Right center)"
	default:
		return fmt.Sprintf("Button %d", int(button))
	}
}

func (w *KeyBindingEditorWindow) drawGamepadBindings() {
	gamepads := ebiten.AppendGamepadIDs(nil)
	if len(gamepads) == 0 {
		imgui.Text("No gamepads connected")
		return
	}

	// Show current gamepad bindings
	if imgui.CollapsingHeaderTreeNodeFlagsV("Current Gamepad Bindings", imgui.TreeNodeFlagsDefaultOpen) {
		if len(w.gamepadBindings) == 0 {
			imgui.Text("No gamepad bindings configured")
		} else {
			for i, binding := range w.gamepadBindings {
				actionStr := binding.Action.String()
				buttonName := getGamepadButtonName(binding.Button)
				displayName := fmt.Sprintf("Gamepad %d: %s", binding.GamepadID, buttonName)

				imgui.PushIDInt(int32(2000 + i))
				imgui.Text(fmt.Sprintf("%s: %s", actionStr, displayName))
				imgui.SameLine()

				if imgui.Button("Remove") {
					w.manager.Unbind(GamepadButton{
						GamepadID: binding.GamepadID,
						Button:    binding.Button,
					})
					w.needsRefresh = true
				}

				imgui.PopID()
				imgui.Separator()
			}
		}
	}

	// Add new gamepad binding
	if imgui.CollapsingHeaderTreeNodeFlagsV("Add Gamepad Binding", 0) {
		// Action selector
		imgui.Text("Select action:")
		actions := []Action{
			ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight,
			ActionAttack, ActionUseAbility1, ActionUseAbility2, ActionUseAbility3,
			ActionPause, ActionToggleDebug,
		}

		for i, action := range actions {
			imgui.PushIDInt(int32(3000 + i))
			if imgui.Button(action.String()) {
				w.selectedAction = action
			}
			imgui.PopID()
			if i%3 != 2 && i < len(actions)-1 {
				imgui.SameLine()
			}
		}

		imgui.Separator()
		imgui.Text(fmt.Sprintf("Selected: %s", w.selectedAction.String()))
		imgui.Separator()

		// Button selection
		imgui.Text("Select button to bind:")
		standardButtons := []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonRightBottom,
			ebiten.StandardGamepadButtonRightRight,
			ebiten.StandardGamepadButtonRightLeft,
			ebiten.StandardGamepadButtonRightTop,
			ebiten.StandardGamepadButtonLeftTop,
			ebiten.StandardGamepadButtonLeftRight,
			ebiten.StandardGamepadButtonLeftBottom,
			ebiten.StandardGamepadButtonLeftLeft,
			ebiten.StandardGamepadButtonCenterLeft,
			ebiten.StandardGamepadButtonCenterRight,
		}

		for i, button := range standardButtons {
			buttonName := getGamepadButtonName(button)
			imgui.PushIDInt(int32(4000 + i))
			if imgui.Button(buttonName) {
				w.manager.Bind(GamepadButton{
					GamepadID: gamepads[0], // Use first gamepad
					Button:    button,
				}, w.selectedAction)
				w.needsRefresh = true
			}
			imgui.PopID()
			if i%2 != 1 && i < len(standardButtons)-1 {
				imgui.SameLine()
			}
		}
	}
}

func (w *KeyBindingEditorWindow) listenForKeyPress() {
	for key := ebiten.KeyA; key <= ebiten.KeyMax; key++ {
		if ebiten.IsKeyPressed(key) {
			if int(key) > 0 && int(key) < int(ebiten.KeyMax) {
				if w.rebindMode {
					w.manager.Unbind(KeyboardKey{Key: w.oldKey})
					w.manager.Bind(KeyboardKey{Key: key}, w.selectedAction)
				} else {
					w.manager.Bind(KeyboardKey{Key: key}, w.selectedAction)
				}
				w.needsRefresh = true
			}
			w.listening = false
			w.rebindMode = false
			break
		}
	}
}
