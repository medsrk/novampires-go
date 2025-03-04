// internal/game/scene/test_scene.go
package scene

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/entity"
	"novampires-go/internal/game/player"
)

// Dependencies contains all external dependencies needed by scenes
type Dependencies struct {
	InputManager common.InputProvider
	Renderer     *entity.RendererAdapter
	ScreenWidth  int
	ScreenHeight int
}

// TestScene implements a test scene with moving targets
type TestScene struct {
	deps       Dependencies
	player     *player.Player
	targets    []common.TargetInfo
	frameCount int
}

// NewTestScene creates a new test scene
func NewTestScene(deps Dependencies) *TestScene {
	// Create initial targets
	targets := createInitialTargets(deps.ScreenWidth, deps.ScreenHeight)

	// Create player at center of screen
	initialPos := common.Vector2{
		X: float64(deps.ScreenWidth) / 2,
		Y: float64(deps.ScreenHeight) / 2,
	}
	player := player.NewPlayer(deps.InputManager, initialPos)

	return &TestScene{
		deps:       deps,
		player:     player,
		targets:    targets,
		frameCount: 0,
	}
}

// createInitialTargets creates initial target objects
func createInitialTargets(screenWidth, screenHeight int) []common.TargetInfo {
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

// Update updates the scene
func (s *TestScene) Update() error {
	s.frameCount++

	// Update player with current targets
	s.player.Update(s.targets)

	// Move targets in circular patterns
	centerX := s.deps.ScreenWidth / 2
	centerY := s.deps.ScreenHeight / 2
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

// Draw draws the scene
func (s *TestScene) Draw(screen *ebiten.Image) {
	// Draw background grid
	s.deps.Renderer.DrawGrid(screen)

	// Draw all targets
	for _, target := range s.targets {
		drawTarget(screen, s.deps.Renderer, target)
	}

	// Draw player
	s.player.Draw(screen, s.deps.Renderer)

	// Draw UI
	// This would be better handled by a proper UI system
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), 8, 8)
}

// Helper function to draw a target
func drawTarget(screen *ebiten.Image, renderer *entity.RendererAdapter, target common.TargetInfo) {
	// This would be better handled by a proper target entity
	// For now, we'll just use the renderer adapter

	// Draw target circle
	renderer.DrawCircle(
		screen,
		target.Pos,
		target.Radius,
		color.RGBA{204, 0, 0, 255},
	)

	// Draw velocity vector
	velEndPos := target.Pos.Add(target.Vel.Scale(5))
	renderer.DrawLine(
		screen,
		target.Pos,
		velEndPos,
		2.0,
		color.RGBA{255, 255, 100, 200},
	)

	// Draw health bar
	healthBarPos := common.Vector2{
		X: target.Pos.X - target.Radius,
		Y: target.Pos.Y - target.Radius - 10,
	}

	// Draw health bar using renderer
	renderer.DrawHealthBar(
		screen,
		healthBarPos,
		target.Radius*2,
		4.0,
		target.Priority/4.0,
	)
}

func (s *TestScene) GetPlayer() *player.Player {
	return s.player
}
