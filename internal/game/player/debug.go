// internal/game/player/debug.go
package player

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
				frameDuration := int(1000.0 / w.frameRate)

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
