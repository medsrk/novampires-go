package input

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/engine/debug"
)

type DebugWindow struct {
	manager *Manager
	open    bool
}

func NewDebugWindow(manager *Manager) *DebugWindow {
	return &DebugWindow{
		manager: manager,
		open:    true,
	}
}

func (w *DebugWindow) Draw() {
	if !w.open {
		return
	}

	// Use the FixedWindow component
	debug.FixedWindow("Input Debug", 400, 500, func() {
		// Input states section
		debug.CollapsingSection("Input States", func() {
			actions := []Action{
				ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight,
				ActionAutoAttack, ActionUseAbility1, ActionUseAbility2, ActionUseAbility3,
				ActionPause, ActionToggleDebug,
			}

			for _, a := range actions {
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
		// Aim vector section
		debug.CollapsingSection("Aim Vector", func() {
			dx, dy := w.manager.GetAimVector()
			debug.AimInfo(dx, dy)
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
	})
}
