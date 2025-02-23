package input

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/debug"
	"unsafe"
)

type DebugWindow struct {
	manager *Manager
	open    bool
	openPtr unsafe.Pointer
}

func NewDebugWindow(manager *Manager) *DebugWindow {
	w := &DebugWindow{
		manager: manager,
		open:    true,
	}

	w.openPtr = unsafe.Pointer(&w.open)

	return w
}

func (w *DebugWindow) Draw() {
	if !w.open {
		return
	}

	visible := imgui.BeginV("Input Window", (*bool)(w.openPtr), imgui.WindowFlagsNone)
	if visible {
		// Input states section
		debug.CollapsingSection("Input States", func() {

			for _, a := range common.Actions {
				debug.InputActionState(
					a.String(),
					w.manager.IsPressed(a),
					w.manager.JustPressed(a),
					w.manager.JustReleased(a),
				)
			}
		})

		// Movement vector section
		debug.CollapsingSection("Movement Vector", func() {
			dx, dy := w.manager.GetMovementVector()
			debug.MovementInfo(dx, dy)
		})

		// Mouse position section
		debug.CollapsingSection("Mouse Position", func() {
			screenX, screenY := w.manager.GetMousePositionScreen()
			worldX, worldY := w.manager.GetMousePositionWorld()
			debug.MousePosition(screenX, screenY, worldX, worldY)
		})

		// Connected gamepads section
		debug.CollapsingSection("Connected Devices", func() {
			// Show connected gamepads
			gamepads := ebiten.AppendGamepadIDs(nil)
			if len(gamepads) == 0 {
				imgui.Text("No gamepads connected")
			} else {
				imgui.Text(fmt.Sprintf("%d gamepad(s) connected:", len(gamepads)))
				for _, id := range gamepads {
					name := ebiten.GamepadName(id)
					if name == "" {
						name = fmt.Sprintf("Gamepad %d", id)
					}
					imgui.Text(fmt.Sprintf("- %s", name))
				}
			}
		})
	}
	imgui.End()

}

func (w *DebugWindow) Name() string {
	return common.WindowInputDebug
}

func (w *DebugWindow) Toggle() {
	w.open = !w.open
}

func (w *DebugWindow) IsOpen() bool {
	return w.open
}

func (w *DebugWindow) SetOpen(open bool) {
	w.open = open
}

func (w *DebugWindow) Close() {
	w.open = false
}
