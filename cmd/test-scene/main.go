// cmd/test-scene/main.go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/debug"
	"novampires-go/internal/engine/input"
	"novampires-go/internal/game/player"
	"novampires-go/internal/game/rendering"
)

const (
	screenWidth  = 1600
	screenHeight = 900
)

type Game struct {
	input            *input.Manager
	debugMgr         *debug.Manager
	renderer         *rendering.Renderer
	playerController *player.Controller
	targets          []common.TargetInfo
	frameCount       int
}

func NewGame() *Game {
	// Initialize input manager
	im := input.New()

	// Initialize debug manager
	dm := debug.New()
	dm.AddWindow(im.CreateDebugWindow())
	keyBindEditor := input.NewKeyBindingEditorWindow(im)
	dm.AddWindow(keyBindEditor)
	dm.SetDisplaySize(screenWidth, screenHeight)

	// Initialize renderer
	renderCfg := rendering.DefaultRenderConfig()
	renderer := rendering.NewRenderer(renderCfg)

	// Create player controller at center of screen
	startPos := common.Vector2{
		X: float64(screenWidth) / 2,
		Y: float64(screenHeight) / 2,
	}
	controllerCfg := player.DefaultConfig()
	playerController := player.NewController(im, startPos, controllerCfg)

	// Create some target entities
	targets := createInitialTargets()

	return &Game{
		input:            im,
		debugMgr:         dm,
		renderer:         renderer,
		playerController: playerController,
		targets:          targets,
		frameCount:       0,
	}
}

func createInitialTargets() []common.TargetInfo {
	targets := make([]common.TargetInfo, 0, 4)

	// Add some targets at different positions
	enemyPositions := []common.Vector2{
		{X: screenWidth / 4, Y: screenHeight / 4},
		{X: screenWidth * 3 / 4, Y: screenHeight / 4},
		{X: screenWidth / 4, Y: screenHeight * 3 / 4},
		{X: screenWidth * 3 / 4, Y: screenHeight * 3 / 4},
	}

	for i, pos := range enemyPositions {
		// Initial velocity in circular pattern direction
		angle := float64(i) * math.Pi / 2 // Distribute angles
		velocity := common.Vector2{
			X: math.Cos(angle) * 2,
			Y: math.Sin(angle) * 2,
		}

		targets = append(targets, common.TargetInfo{
			ID:       uint64(i + 1),
			Pos:      pos,
			Vel:      velocity,
			Radius:   15,
			Priority: float64(i + 1), // Different priorities
		})
	}

	return targets
}

func (g *Game) Update() error {
	g.frameCount++

	// Update the debug manager
	g.debugMgr.Update()

	// Toggle debug if needed
	if g.input.JustPressed(input.ActionToggleDebug) {
		g.debugMgr.Toggle()
	}

	// Update player controller with current targets
	g.playerController.Update(g.targets)

	// Move targets in circular patterns
	for i := range g.targets {
		// Calculate circular movement
		centerX := screenWidth / 2
		centerY := screenHeight / 2
		radius := 300.0

		// Each target moves at slightly different speeds
		speed := 0.005 + float64(i)*0.002

		// Calculate position on circle
		angle := float64(g.frameCount)*speed + float64(i)*(math.Pi/2)
		g.targets[i].Pos = common.Vector2{
			X: float64(centerX) + math.Cos(angle)*radius,
			Y: float64(centerY) + math.Sin(angle)*radius,
		}

		// Calculate velocity (tangent to circle)
		g.targets[i].Vel = common.Vector2{
			X: -math.Sin(angle) * 5,
			Y: math.Cos(angle) * 5,
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with background color
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw a grid for reference
	g.drawGrid(screen)

	// Get player position and information
	playerPos := g.playerController.GetPosition()
	playerRot := g.playerController.GetRotation()
	aimDir := g.playerController.GetAimDirection()

	// Draw player
	g.renderer.DrawPlayerCharacter(screen, playerPos, playerRot, 20)

	// Draw aim line
	g.renderer.DrawAimLine(screen, playerPos, aimDir, 200)

	// Draw targets
	for _, target := range g.targets {
		// Draw the target circle
		g.renderer.DrawCircle(screen, target.Pos, target.Radius, color.RGBA{204, 0, 0, 255})

		// Draw velocity vector
		endPoint := target.Pos.Add(target.Vel.Scale(5))
		g.renderer.DrawLine(
			screen,
			target.Pos,
			endPoint,
			2,
			color.RGBA{255, 255, 100, 200},
		)

		// Draw health bar for each target
		healthBarWidth := target.Radius * 2
		healthBarHeight := 4.0
		healthBarPos := common.Vector2{
			X: target.Pos.X - healthBarWidth/2,
			Y: target.Pos.Y - target.Radius - 10,
		}

		// Use priority as "health" for visualization
		healthPercent := target.Priority / 4.0 // Assuming max priority is 4.0
		g.renderer.DrawHealthBar(screen, healthBarPos, healthBarWidth, healthBarHeight, healthPercent)
	}

	// Draw debug UI
	g.debugMgr.Draw(screen)
}

// Draw a reference grid
func (g *Game) drawGrid(screen *ebiten.Image) {
	gridColor := color.RGBA{50, 50, 60, 80}
	gridSpacing := 100

	// Vertical lines
	for x := gridSpacing; x < screenWidth; x += gridSpacing {
		vector.StrokeLine(
			screen,
			float32(x),
			0,
			float32(x),
			float32(screenHeight),
			1,
			gridColor,
			false,
		)
	}

	// Horizontal lines
	for y := gridSpacing; y < screenHeight; y += gridSpacing {
		vector.StrokeLine(
			screen,
			0,
			float32(y),
			float32(screenWidth),
			float32(y),
			1,
			gridColor,
			false,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Update ImGui display size if window is resized
	g.debugMgr.SetDisplaySize(float32(outsideWidth), float32(outsideHeight))
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("NoVampires Test Scene")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
