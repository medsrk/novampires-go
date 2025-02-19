package input

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"sort"
)

type KeyBindingEditorWindow struct {
	manager        *Manager
	open           bool
	listening      bool
	selectedAction Action
	rebindMode     bool
	oldKey         ebiten.Key
	// Sort bindings by action and keep them stable
	stableBindings []KeyActionPair
	needsRefresh   bool
}

// KeyActionPair represents a key binding with its action
type KeyActionPair struct {
	Key    ebiten.Key
	Action Action
}

func NewKeyBindingEditorWindow(manager *Manager) *KeyBindingEditorWindow {
	return &KeyBindingEditorWindow{
		manager:        manager,
		open:           true,
		listening:      false,
		selectedAction: MoveUp,
		rebindMode:     false,
		oldKey:         ebiten.Key(0),
		stableBindings: nil,
		needsRefresh:   true,
	}
}

// GetKeyBindings returns a map of keyboard key bindings
func (m *Manager) GetKeyBindings() map[ebiten.Key]Action {
	bindings := make(map[ebiten.Key]Action)

	// Extract keyboard bindings from the manager's internal map
	for input, action := range m.bindings {
		if keyInput, ok := input.(KeyboardKey); ok {
			bindings[keyInput.Key] = action
		}
	}

	return bindings
}

func (w *KeyBindingEditorWindow) refreshBindings() {
	if !w.needsRefresh && w.stableBindings != nil {
		return
	}

	// Get current bindings
	bindings := w.manager.GetKeyBindings()

	// Convert map to sorted slice for stable display
	w.stableBindings = make([]KeyActionPair, 0, len(bindings))
	for key, action := range bindings {
		w.stableBindings = append(w.stableBindings, KeyActionPair{
			Key:    key,
			Action: action,
		})
	}

	// Sort by action enum value for stable ordering
	sort.Slice(w.stableBindings, func(i, j int) bool {
		return int(w.stableBindings[i].Action) < int(w.stableBindings[j].Action)
	})

	w.needsRefresh = false
}

func (w *KeyBindingEditorWindow) Draw() {
	if !w.open {
		return
	}

	// Refresh our bindings list if needed
	w.refreshBindings()

	// Basic fixed window setup
	imgui.SetNextWindowSize(imgui.Vec2{X: 350, Y: 450})

	if imgui.Begin("Key Bindings") {
		// Edit existing bindings section
		if imgui.CollapsingHeaderTreeNodeFlagsV("Edit Bindings", 0) {
			for i, binding := range w.stableBindings {
				actionStr := binding.Action.String()

				// Use numeric ID for each row
				imgui.PushIDInt(int32(i))

				// Display the binding info
				imgui.Text(fmt.Sprintf("%s:", actionStr))
				imgui.SameLine()

				// Display current key
				keyName := binding.Key.String()
				if len(keyName) > 3 && keyName[:3] == "Key" {
					keyName = keyName[3:] // Remove "Key" prefix
				}

				// Rebinding button - simple version without color changes
				buttonText := keyName
				if w.listening && w.rebindMode && w.oldKey == binding.Key {
					buttonText = "Press a key..."
				}

				if imgui.Button(buttonText) {
					if w.listening && w.rebindMode && w.oldKey == binding.Key {
						// Cancel if clicked again
						w.listening = false
						w.rebindMode = false
					} else {
						// Start listening for a new key
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

		imgui.Separator()

		// Add new bindings section
		if imgui.CollapsingHeaderTreeNodeFlagsV("Add New Binding", 0) {
			// Action selector
			imgui.Text("Select action to bind:")

			actions := []Action{
				MoveUp,
				MoveDown,
				MoveLeft,
				MoveRight,
				Attack,
				UseAbility1,
				UseAbility2,
				UseAbility3,
				Pause,
				ToggleDebug,
			}

			// Simple buttons for action selection
			for i, action := range actions {
				imgui.PushIDInt(int32(1000 + i))

				if imgui.Button(action.String()) {
					w.selectedAction = action
				}

				imgui.PopID()

				// 3 buttons per row
				if i%3 != 2 && i < len(actions)-1 {
					imgui.SameLine()
				}
			}

			imgui.Separator()

			// Show selected action
			imgui.Text(fmt.Sprintf("Selected: %s", w.selectedAction.String()))
			imgui.Separator()

			// Bind button - simple version without color changes
			bindButtonText := "Bind New Key"
			if w.listening && !w.rebindMode {
				bindButtonText = "Press a key to bind..."
			}

			if imgui.Button(bindButtonText) {
				if w.listening && !w.rebindMode {
					// Cancel if clicked again
					w.listening = false
				} else {
					// Start listening
					w.listening = true
					w.rebindMode = false
				}
			}
		}
	}
	imgui.End()

	// Check for key presses if we're listening
	if w.listening {
		w.listenForKeyPress()
	}
}

func (w *KeyBindingEditorWindow) listenForKeyPress() {
	for key := ebiten.KeyA; key <= ebiten.KeyMax; key++ {
		if ebiten.IsKeyPressed(key) {
			// Simple validation
			if int(key) > 0 && int(key) < int(ebiten.KeyMax) {
				if w.rebindMode {
					// Remove old binding
					w.manager.Unbind(KeyboardKey{Key: w.oldKey})
					// Add new binding
					w.manager.Bind(KeyboardKey{Key: key}, w.selectedAction)
				} else {
					// Just add a new binding
					w.manager.Bind(KeyboardKey{Key: key}, w.selectedAction)
				}
				// Mark for refresh
				w.needsRefresh = true
			}

			w.listening = false
			w.rebindMode = false
			break
		}
	}
}
