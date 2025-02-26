package game

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

	game *Game
}

func NewDebugWindow(manager *debug.Manager, game *Game) *DebugWindow {
	return &DebugWindow{
		manager: manager,
		open:    true,
		game:    game,
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
		g := w.game
		debug.CollapsingSection("Game State", func() {
			debug.LabeledValue("Game Active Scene", fmt.Sprintf("%v", g.activeScene), nil)
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
