// cmd/input-test/main.go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"novampires-go/internal/engine/debug"
	"novampires-go/internal/engine/input"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

type Game struct {
	input    *input.Manager
	debugMgr *debug.Manager
}

func NewGame() *Game {
	im := input.New()
	dm := debug.New()

	// Add the debug window from the input system
	dm.AddWindow(im.CreateDebugWindow())
	keyBindEditor := input.NewKeyBindingEditorWindow(im)
	dm.AddWindow(keyBindEditor)

	// Set the display size for ImGui
	dm.SetDisplaySize(screenWidth, screenHeight)

	return &Game{
		input:    im,
		debugMgr: dm,
	}
}

func (g *Game) Update() error {
	// Update the debug manager
	g.debugMgr.Update()

	// Toggle debug if needed
	if g.input.JustPressed(input.ToggleDebug) {
		g.debugMgr.Toggle()
	}

	return nil
}

var RectX, RectY float64 = 320, 240

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the rectangle
	dx, dy := g.input.GetMovementVector()
	RectX += dx
	RectY += dy
	vector.DrawFilledRect(screen, float32(RectX), float32(RectY), 10, 10, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)

	// Draw the debug UI
	g.debugMgr.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Update ImGui display size if window is resized
	g.debugMgr.SetDisplaySize(float32(outsideWidth), float32(outsideHeight))
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Input Test")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
