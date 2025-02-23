// cmd/test-scene/main.go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/camera"
	"novampires-go/internal/engine/debug"
	"novampires-go/internal/engine/input"
	"novampires-go/internal/game"
	"novampires-go/internal/game/config"
	"novampires-go/internal/game/rendering"
)

const (
	screenWidth  = 1600
	screenHeight = 900
)

// TestScene implements the game.Scene interface
type TestScene struct {
	game       *game.Game
	targets    []common.TargetInfo
	frameCount int
}

func NewTestScene(g *game.Game) *TestScene {
	return &TestScene{
		game:       g,
		targets:    createInitialTargets(),
		frameCount: 0,
	}
}

func createInitialTargets() []common.TargetInfo {
	targets := make([]common.TargetInfo, 0, 4)

	// Add targets in a circular pattern around center
	centerX := screenWidth / 2
	centerY := screenHeight / 2
	radius := 300.0

	for i := 0; i < 4; i++ {
		// Distribute evenly around circle
		angle := float64(i) * math.Pi / 2
		pos := common.Vector2{
			X: float64(centerX) + math.Cos(angle)*radius,
			Y: float64(centerY) + math.Sin(angle)*radius,
		}

		// Initial velocity tangent to circle
		velocity := common.Vector2{
			X: -math.Sin(angle) * 5,
			Y: math.Cos(angle) * 5,
		}

		targets = append(targets, common.TargetInfo{
			ID:       uint64(i + 1),
			Pos:      pos,
			Vel:      velocity,
			Radius:   15,
			Priority: float64(i + 1),
		})
	}

	return targets
}

func (s *TestScene) Update() error {
	s.frameCount++

	playerController := s.game.GetPlayerController()

	// Update player controller with current targets
	playerController.Update(s.targets)

	// Move targets in circular patterns
	centerX := screenWidth / 2
	centerY := screenHeight / 2
	radius := 300.0

	for i := range s.targets {
		// Each target moves at slightly different speeds
		speed := 0.005 + float64(i)*0.002
		angle := float64(s.frameCount)*speed + float64(i)*(math.Pi/2)

		// Update position
		s.targets[i].Pos = common.Vector2{
			X: float64(centerX) + math.Cos(angle)*radius,
			Y: float64(centerY) + math.Sin(angle)*radius,
		}

		// Update velocity (tangent to circle)
		s.targets[i].Vel = common.Vector2{
			X: -math.Sin(angle) * 5,
			Y: math.Cos(angle) * 5,
		}
	}

	return nil
}

func (s *TestScene) Draw(screen *ebiten.Image) {
	deps := s.game.GetDependencies()
	playerController := s.game.GetPlayerController()

	// Clear the screen
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Get player information
	playerPos := playerController.GetPosition()
	playerRot := playerController.GetRotation()
	aimDir := playerController.GetAimDirection() // World-space aim direction

	// Transform player position to screen coordinates
	screenPos := deps.Camera.WorldToScreen(playerPos)

	// Draw background grid
	deps.Renderer.DrawGrid(screen, deps.Camera)

	// Draw player
	deps.Renderer.DrawPlayerCharacter(screen, screenPos, playerRot, 20)

	// Draw aim line:  CORRECTED
	deps.Renderer.DrawAimLine(screen, screenPos, aimDir, 200)

	// ... (rest of the drawing code) ...
	// Draw targets
	for _, target := range s.targets {
		// Transform target position
		targetScreenPos := deps.Camera.WorldToScreen(target.Pos)

		// Draw target circle
		deps.Renderer.DrawCircle(screen, targetScreenPos, target.Radius, color.RGBA{204, 0, 0, 255})

		// Draw velocity vector
		velEndPos := deps.Camera.WorldToScreen(target.Pos.Add(target.Vel.Scale(5)))
		deps.Renderer.DrawLine(
			screen,
			targetScreenPos,
			velEndPos,
			2,
			color.RGBA{255, 255, 100, 200},
		)

		// Draw health bar
		healthBarWidth := target.Radius * 2
		healthBarHeight := 4.0
		healthBarPos := common.Vector2{
			X: targetScreenPos.X - healthBarWidth/2,
			Y: targetScreenPos.Y - target.Radius - 10,
		}

		deps.Renderer.DrawHealthBar(screen, healthBarPos, healthBarWidth, healthBarHeight, target.Priority/4.0)
	}

	// Draw FPS counter
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("NoVampires Test Scene")
	// Initialize core systems
	im := input.New()
	dm := debug.New(debug.Deps{InputManager: im})

	// Initialize camera
	cam := camera.New()
	im.SetCamera(cam)

	// Create game dependencies
	deps := game.Dependencies{
		InputManager: im,
		DebugManager: dm,
		Camera:       cam,
		Renderer:     rendering.NewRenderer(rendering.DefaultRenderConfig()),
		Config:       config.Default(),
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}

	// Create game instance
	g := game.NewGame(deps)

	// Create and set test scene
	scene := NewTestScene(g)
	g.SetScene(scene)

	// Set up window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("NoVampires Test Scene")

	// Run the game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
