package input

import (
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
	debug.FixedWindow("Input Debug", 300, 400, func() {
		// Input states section
		debug.CollapsingSection("Input States", func() {
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
	})
}
