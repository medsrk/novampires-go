package player

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/debug"
)

type DebugWindow struct {
	manager *debug.Manager
	open    bool

	controller *Controller
}

func NewDebugWindow(manager *debug.Manager, controller *Controller) *DebugWindow {
	return &DebugWindow{
		manager:    manager,
		open:       true,
		controller: controller,
	}
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
