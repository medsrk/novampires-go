// internal/engine/sprite/debug.go
package sprite

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"novampires-go/internal/common"
	"unsafe"
)

// DebugWindow represents a debug window for sprite animation
type DebugWindow struct {
	open    bool
	openPtr unsafe.Pointer

	// Reference to the player controller
	controller interface{} // Using interface{} to avoid circular imports

	// Animation parameters for adjustment
	frameRate    float32
	frameRatePtr unsafe.Pointer

	scale    float32
	scalePtr unsafe.Pointer

	flipX    bool
	flipXPtr unsafe.Pointer

	// Animation preview
	previewSize    float32
	previewSizePtr unsafe.Pointer

	// Selected animation
	currentAnim string
	animations  []string
}

// NewDebugWindow creates a new sprite debug window
func NewDebugWindow(controller interface{}) *DebugWindow {
	w := &DebugWindow{
		open:        true,
		controller:  controller,
		frameRate:   100.0,
		scale:       1.0,
		flipX:       false,
		previewSize: 128.0,
		currentAnim: "idle",
		animations:  []string{"idle", "walk", "attack"},
	}

	w.openPtr = unsafe.Pointer(&w.open)
	w.frameRatePtr = unsafe.Pointer(&w.frameRate)
	w.scalePtr = unsafe.Pointer(&w.scale)
	w.flipXPtr = unsafe.Pointer(&w.flipX)
	w.previewSizePtr = unsafe.Pointer(&w.previewSize)

	return w
}

// Draw renders the debug window
func (w *DebugWindow) Draw() {
	if !w.open {
		return
	}

	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 400}, imgui.CondFirstUseEver)

	if imgui.BeginV("Sprite Debugger", (*bool)(w.openPtr), imgui.WindowFlagsNone) {
		// Animation selector
		imgui.Text("Animation:")
		imgui.SameLine()

		if imgui.BeginCombo("##animation", w.currentAnim) {
			for _, anim := range w.animations {
				isSelected := w.currentAnim == anim
				if imgui.SelectableBool(anim) {
					w.currentAnim = anim
					// Apply animation change to controller if it has the method
					if ctrl, ok := w.controller.(interface{ PlayAnimation(string) }); ok {
						ctrl.PlayAnimation(w.currentAnim)
					}
				}
				if isSelected {
					imgui.SetItemDefaultFocus()
				}
			}
			imgui.EndCombo()
		}

		// Animation parameters
		imgui.Separator()
		imgui.Text("Animation Parameters")

		if imgui.SliderFloat("Frame Rate", (*float32)(w.frameRatePtr), 30.0, 300.0) {
			// Apply frame rate change to controller if it has the method
			if ctrl, ok := w.controller.(interface{ SetFrameRate(float32) }); ok {
				ctrl.SetFrameRate(w.frameRate)
			}
		}

		if imgui.SliderFloat("Scale", (*float32)(w.scalePtr), 0.5, 3.0) {
			// Apply scale change to controller if it has the method
			if ctrl, ok := w.controller.(interface{ SetScale(float32) }); ok {
				ctrl.SetScale(w.scale)
			}
		}

		if imgui.Checkbox("Flip X", (*bool)(w.flipXPtr)) {
			// Apply flip change to controller if it has the method
			if ctrl, ok := w.controller.(interface{ SetFlipX(bool) }); ok {
				ctrl.SetFlipX(w.flipX)
			}
		}

		// Animation preview
		imgui.Separator()
		imgui.Text("Preview")

		imgui.SliderFloat("Preview Size", (*float32)(w.previewSizePtr), 64.0, 256.0)

		// Preview animation
		//previewPos := imgui.CursorPos()
		imgui.InvisibleButton("preview", imgui.Vec2{X: w.previewSize, Y: w.previewSize})

		// Display current animation info
		imgui.Separator()
		imgui.Text("Animation Info")

		// Display current frame if controller has the method
		if ctrl, ok := w.controller.(interface{ GetCurrentFrame() int }); ok {
			frame := ctrl.GetCurrentFrame()
			imgui.Text(fmt.Sprintf("Current Frame: %d", frame))
		}

		// Display animation duration if controller has the method
		if ctrl, ok := w.controller.(interface{ GetAnimationDuration() float32 }); ok {
			duration := ctrl.GetAnimationDuration()
			imgui.Text(fmt.Sprintf("Duration: %.2f ms", duration))
		}
	}
	imgui.End()
}

// Name returns the window name for the debug manager
func (w *DebugWindow) Name() string {
	return common.WindowSpriteDebug
}

// IsOpen returns whether the window is open
func (w *DebugWindow) IsOpen() bool {
	return w.open
}

// Toggle toggles the window's visibility
func (w *DebugWindow) Toggle() {
	w.open = !w.open
}

// Close closes the window
func (w *DebugWindow) Close() {
	w.open = false
}

func (a *Animation) CreateDebugWindow() *DebugWindow {
	return NewDebugWindow(a)
}
