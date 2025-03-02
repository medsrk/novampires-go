// internal/game/player/debug.go
package player

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/debug"
	"unsafe"
)

type DebugWindow struct {
	manager    *debug.Manager
	open       bool
	controller *Controller

	// Additional debug controls
	scale    float32
	scalePtr unsafe.Pointer
	flipX    bool
	flipXPtr unsafe.Pointer

	// Animation selection
	animations   []string
	currentAnim  string
	frameRate    float32
	frameRatePtr unsafe.Pointer
}

func NewDebugWindow(manager *debug.Manager, controller *Controller) *DebugWindow {
	// Prepare animation list from controller
	animations := []string{"idle", "walk", "attack"}

	w := &DebugWindow{
		manager:     manager,
		open:        true,
		controller:  controller,
		scale:       float32(controller.spriteScale),
		flipX:       controller.flipX,
		animations:  animations,
		currentAnim: controller.currentAnim,
		frameRate:   100.0, // Default frame rate
	}

	w.scalePtr = unsafe.Pointer(&w.scale)
	w.flipXPtr = unsafe.Pointer(&w.flipX)
	w.frameRatePtr = unsafe.Pointer(&w.frameRate)

	return w
}

func (w *DebugWindow) Draw() {
	if !w.open {
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		w.open = true
	}

	isOpen := w.open
	if imgui.BeginV("Player Debug", &isOpen, imgui.WindowFlagsNone) {
		// Player state section
		c := w.controller
		debug.CollapsingSection("Player State", func() {
			debug.LabeledValue("Position", c.pos.String(), nil)
			debug.LabeledValue("Velocity", c.vel.String(), nil)
			debug.LabeledValue("Rotation", fmt.Sprintf("%.2f", c.rotation), nil)
			debug.LabeledValue("Using Gamepad", fmt.Sprintf("%v", c.usingGamepad), nil)
			debug.LabeledValue("LastMousePos", fmt.Sprintf("X: %v, Y: %v", c.lastMouseX, c.lastMouseY), &imgui.Vec4{X: 0, Y: 1, Z: 0, W: 1})
			debug.LabeledValue("lastAimDir", fmt.Sprintf("X: %v, Y: %v", c.lastAimDx, c.lastAimDy), &imgui.Vec4{X: 0, Y: 1, Z: 0, W: 1})
		})

		// Sprite controls section
		debug.CollapsingSection("Sprite Controls", func() {
			// Sprite scale slider
			if imgui.SliderFloat("Sprite Scale", (*float32)(w.scalePtr), 0.5, 10.0) {
				w.controller.spriteScale = float64(w.scale)
			}

			// Flip X checkbox
			if imgui.Checkbox("Flip Horizontally", (*bool)(w.flipXPtr)) {
				w.controller.flipX = w.flipX
			}

			// Animation selection
			imgui.Text("Animation:")
			imgui.SameLine()

			if imgui.BeginCombo("##animation", w.currentAnim) {
				for _, anim := range w.animations {
					isSelected := w.currentAnim == anim
					if imgui.SelectableBool(anim) {
						w.currentAnim = anim
						w.controller.PlayAnimation(anim)
					}
					if isSelected {
						imgui.SetItemDefaultFocus()
					}
				}
				imgui.EndCombo()
			}

			// Frame rate slider
			if imgui.SliderFloat("Frame Rate", (*float32)(w.frameRatePtr), 30.0, 300.0) {
				// Convert frame rate to frame duration
				frameDuration := int(100.0 / w.frameRate)

				// Update all animations
				for _, anim := range w.controller.animations {
					for i := range anim.Frames {
						anim.Frames[i].Duration = frameDuration
					}
				}
			}

			// Current animation info
			if anim, exists := w.controller.animations[w.controller.currentAnim]; exists {
				imgui.Text(fmt.Sprintf("Current Frame: %d/%d", anim.GetCurrentFrameInt()+1, len(anim.Frames)))
			}

			// Is animation reversed
			imgui.Text(fmt.Sprintf("Reversed: %v", w.controller.IsReversed()))

			// Add this to the DebugWindow Draw method in player/debug.go

			// Eye Controller Debug Section
			debug.CollapsingSection("Eye Controller", func() {
				// Get the eye controller from the player controller
				eyeController := w.controller.eyeController

				// Show current eye direction
				directionNames := map[Direction]string{
					LookingCenter:    "Center",
					LookingRight:     "Right",
					LookingUp:        "Up",
					LookingDown:      "Down",
					LookingUpRight:   "Up-Right",
					LookingDownRight: "Down-Right",
				}

				dirName := "Unknown"
				if name, exists := directionNames[eyeController.direction]; exists {
					dirName = name
				}

				debug.LabeledValue("Eye Direction", dirName, nil)
				debug.LabeledValue("Eye Flipped", fmt.Sprintf("%v", eyeController.flipX), nil)

				// Blinking state
				debug.LabeledValue("Is Blinking", fmt.Sprintf("%v", eyeController.isBlinking), nil)
				debug.LabeledValue("Blink Timer", fmt.Sprintf("%d / %d", eyeController.blinkTimer, eyeController.blinkInterval), nil)

				// Eye position
				debug.LabeledValue("Eye Position", eyeController.position.String(), nil)

				// Manual blink trigger button
				if imgui.Button("Trigger Blink") {
					eyeController.TriggerBlink()
				}

				// Testing buttons for setting eye direction manually
				imgui.Separator()
				imgui.Text("Test Eye Directions:")

				// Basic directions
				if imgui.Button("Center") {
					eyeController.direction = LookingCenter
					eyeController.flipX = false
				}

				imgui.SameLine()
				if imgui.Button("Right") {
					eyeController.direction = LookingRight
					eyeController.flipX = false
				}

				imgui.SameLine()
				if imgui.Button("Left") {
					eyeController.direction = LookingRight
					eyeController.flipX = true
				}

				// Row 2
				if imgui.Button("Up") {
					eyeController.direction = LookingUp
					eyeController.flipX = false
				}

				imgui.SameLine()
				if imgui.Button("Down") {
					eyeController.direction = LookingDown
					eyeController.flipX = false
				}

				// Diagonal rows
				if imgui.Button("Up-Right") {
					eyeController.direction = LookingUpRight
					eyeController.flipX = false
				}

				imgui.SameLine()
				if imgui.Button("Up-Left") {
					eyeController.direction = LookingUpRight
					eyeController.flipX = true
				}

				if imgui.Button("Down-Right") {
					eyeController.direction = LookingDownRight
					eyeController.flipX = false
				}

				imgui.SameLine()
				if imgui.Button("Down-Left") {
					eyeController.direction = LookingDownRight
					eyeController.flipX = true
				}

				// Add manual aim angle test
				imgui.Separator()

				// Add buttons for preset angles to test
				imgui.Text("Test Specific Angles:")

				angleTests := []struct {
					name  string
					angle float64
				}{
					{"0° (Right)", 0},
					{"45° (Down-Right)", math.Pi / 4},
					{"90° (Down)", math.Pi / 2},
					{"135° (Down-Left)", 3 * math.Pi / 4},
					{"180° (Left)", math.Pi},
					{"225° (Up-Left)", 5 * math.Pi / 4},
					{"270° (Up)", 3 * math.Pi / 2},
					{"315° (Up-Right)", 7 * math.Pi / 4},
				}

				// Display in two rows
				for i, test := range angleTests {
					if i > 0 && i%4 == 0 {
						imgui.Spacing()
					}

					if imgui.Button(test.name) {
						// Create a directional vector from this angle
						testDir := common.Vector2{
							X: math.Cos(test.angle),
							Y: math.Sin(test.angle),
						}

						// Apply the direction to the eye controller
						eyeController.UpdateLookDirection(testDir)
					}

					if i%4 != 3 && i < len(angleTests)-1 {
						imgui.SameLine()
					}
				}
			})
		})
	}
	imgui.End()

	w.open = isOpen
}

func (w *DebugWindow) Name() string {
	return common.WindowPlayerDebug
}

func (w *DebugWindow) IsOpen() bool {
	return w.open
}

func (w *DebugWindow) Toggle() {
	w.open = !w.open
}

func (w *DebugWindow) Close() {
	w.open = false
}
