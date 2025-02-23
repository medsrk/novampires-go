package debug

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/common"
	"strings"
)

// ColoredText displays text with the specified color
func ColoredText(text string, color imgui.Vec4) {
	imgui.PushStyleColorVec4(imgui.ColText, color)
	imgui.Text(text)
	imgui.PopStyleColor()
}

// LabeledValue displays a label and value with optional coloring
func LabeledValue(label string, value string, valueColor *imgui.Vec4) {
	imgui.Text(label)
	imgui.SameLine()

	if valueColor != nil {
		ColoredText(value, *valueColor)
	} else {
		imgui.Text(value)
	}
}

// State colors
var (
	ColorActive       = imgui.Vec4{X: 0, Y: 1, Z: 0, W: 1}       // Green
	ColorInactive     = imgui.Vec4{X: 0.5, Y: 0.5, Z: 0.5, W: 1} // Gray
	ColorJustPressed  = imgui.Vec4{X: 1, Y: 1, Z: 0, W: 1}       // Yellow
	ColorJustReleased = imgui.Vec4{X: 1, Y: 0.5, Z: 0, W: 1}     // Orange
)

// InputActionState displays the state of an input action with colored indicators
func InputActionState(name string, isActive bool, isJustPressed bool, isJustReleased bool) {
	imgui.Text(name)
	imgui.SameLine()

	// Active state
	activeColor := ColorInactive
	if isActive {
		activeColor = ColorActive
	}
	ColoredText("Active", activeColor)

	// Just pressed state
	imgui.SameLine()
	pressedColor := ColorInactive
	if isJustPressed {
		pressedColor = ColorJustPressed
	}
	ColoredText("Just Pressed", pressedColor)

	// Just released state
	imgui.SameLine()
	releasedColor := ColorInactive
	if isJustReleased {
		releasedColor = ColorJustReleased
	}
	ColoredText("Just Released", releasedColor)

	imgui.Separator()
}

// FixedWindow creates a fixed-size, non-resizable window
func FixedWindow(title string, width, height float32, drawContents func()) {
	imgui.SetNextWindowSize(imgui.Vec2{X: width, Y: height})
	flags := imgui.WindowFlagsNoResize

	if imgui.BeginV(title, nil, imgui.WindowFlags(flags)) {
		drawContents()
	}
	imgui.End()
}

// CollapsingSection creates a collapsible section
func CollapsingSection(label string, drawContents func()) {
	if imgui.CollapsingHeaderTreeNodeFlagsV(label, 0) {
		drawContents()
	}
	imgui.Separator()
}

// GetDirectionText converts a movement vector to a text direction
func GetDirectionText(dx, dy float64) string {
	if dx == 0 && dy == 0 {
		return "Stationary"
	} else if dx > 0 && dy < 0 {
		return "Up-Right"
	} else if dx < 0 && dy < 0 {
		return "Up-Left"
	} else if dx > 0 && dy > 0 {
		return "Down-Right"
	} else if dx < 0 && dy > 0 {
		return "Down-Left"
	} else if dx == 0 && dy < 0 {
		return "Up"
	} else if dx == 0 && dy > 0 {
		return "Down"
	} else if dx > 0 && dy == 0 {
		return "Right"
	} else if dx < 0 && dy == 0 {
		return "Left"
	}
	return "Unknown"
}

// MovementInfo displays movement vector information
func MovementInfo(dx, dy float64) {
	LabeledValue("Vector:", fmt.Sprintf("(%.2f, %.2f)", dx, dy), nil)
	imgui.Separator()

	strength := float32(dx*dx + dy*dy)
	LabeledValue("Strength:", fmt.Sprintf("%.2f", strength), nil)
	imgui.Separator()

	direction := GetDirectionText(dx, dy)
	LabeledValue("Direction:", direction, nil)
}

// MousePosition displays mouse position information
func MousePosition(screenX, screenY, worldX, worldY int) {
	LabeledValue("Screen:", fmt.Sprintf("(%d, %d)", screenX, screenY), nil)
	LabeledValue("World:", fmt.Sprintf("(%d, %d)", worldX, worldY), nil)
}

// PlayerInfo displays information about the player
func PlayerInfo(pos common.Vector2, rot float64, speed float64) {
	LabeledValue("Position:", fmt.Sprintf("(%.2f, %.2f)", pos.X, pos.Y), nil)
	LabeledValue("Rotation:", fmt.Sprintf("%.2f", rot), nil)
	LabeledValue("Speed:", fmt.Sprintf("%.2f", speed), nil)
}

// KeyBindingEditor displays and allows editing of a key binding
func KeyBindingEditor(label string, currentKey ebiten.Key, onRebind func(ebiten.Key)) {
	imgui.PushIDStrStr(label, currentKey.String())
	defer imgui.PopID()

	// Display current binding
	keyName := currentKey.String()
	if strings.HasPrefix(keyName, "Key") {
		keyName = keyName[3:] // Remove "Key" prefix for cleaner display
	}

	imgui.Text(fmt.Sprintf("%s:", label))
	imgui.SameLine()

	// Change binding button
	if imgui.Button(keyName) {
		// Set this component to listening mode
		// This would be tracked in a global state
		// For simplicity, we'll use a static variable here
		startListeningForKey(onRebind)
	}

	// If this component is in listening mode, display instructions
	if isListeningForKey() {
		imgui.SameLine()
		ColoredText("Press any key...", imgui.Vec4{X: 1, Y: 1, Z: 0, W: 1})
	}
}

// Global state for key binding listening -
// In a real implementation, you'd want to track this per component
var (
	listeningForKey   bool
	keyRebindCallback func(ebiten.Key)
)

func startListeningForKey(callback func(ebiten.Key)) {
	listeningForKey = true
	keyRebindCallback = callback
}

func isListeningForKey() bool {
	return listeningForKey
}

// This needs to be called from your main update loop
func UpdateKeyBindingListeners() {
	if !listeningForKey {
		return
	}

	// Check for any key press
	for key := ebiten.KeyA; key <= ebiten.KeyMax; key++ {
		if ebiten.IsKeyPressed(key) {
			if keyRebindCallback != nil {
				keyRebindCallback(key)
			}
			listeningForKey = false
			keyRebindCallback = nil
			break
		}
	}
}

// ActionBindingList displays and allows editing of action bindings
func ActionBindingList(bindings map[ebiten.Key]interface{}, getActionName func(interface{}) string, onRebind func(oldKey, newKey ebiten.Key)) {
	for key, action := range bindings {
		actionName := getActionName(action)

		KeyBindingEditor(actionName, key, func(newKey ebiten.Key) {
			onRebind(key, newKey)
		})
	}
}
