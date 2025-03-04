package player

import (
	"novampires-go/internal/common"
	"novampires-go/internal/engine/debug"
	"unsafe"
)

type DebugWindow struct {
	manager *debug.Manager
	open    bool
	player  *Player

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

func NewDebugWindow(manager *debug.Manager, player Player) *DebugWindow {
	return &DebugWindow{}
}

func (w *DebugWindow) Draw() {

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
