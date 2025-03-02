package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"math"
	"net/http"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/camera"
	"novampires-go/internal/engine/debug"
	"novampires-go/internal/engine/input"
	"novampires-go/internal/engine/rendering"
	"novampires-go/internal/game"
	"novampires-go/internal/game/config"
	"runtime"
	"time"

	_ "github.com/ebitengine/purego"
	_ "image/png"
	_ "net/http/pprof"
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

// Update the Draw method in TestScene to use the player's scale
func (s *TestScene) Draw(screen *ebiten.Image) {
	deps := s.game.GetDependencies()
	playerController := s.game.GetPlayerController()

	deps.Renderer.BeginFrame(screen)

	// Draw background grid
	deps.Renderer.DrawGrid(screen)

	// Draw all targets
	for _, target := range s.targets {
		deps.Renderer.DrawCircle(screen, target.Pos, target.Radius, color.RGBA{204, 0, 0, 255})
		velEndPos := target.Pos.Add(target.Vel.Scale(5))
		deps.Renderer.DrawLine(screen, target.Pos, velEndPos, 2, color.RGBA{255, 255, 100, 200})

		healthBarPos := common.Vector2{
			X: target.Pos.X - target.Radius,
			Y: target.Pos.Y - target.Radius - 10,
		}
		deps.Renderer.DrawHealthBar(screen, healthBarPos, target.Radius*2, 4.0, target.Priority/4.0)
	}

	// Draw player with layered sprites (body + eyes)
	playerPos := playerController.GetPosition()
	playerRot := playerController.GetRotation()
	playerSprite := playerController.GetSprite()
	playerEyeSprite := playerController.GetEyeSprite() // Get the eye sprite
	playerEyePos := playerController.GetEyePosition()  // Get the eye position offset
	playerFlip := playerController.GetFlipX()
	spriteScale := playerController.GetScale()
	aimDir := playerController.GetAimDirection()

	// flip sprite depending on aim direction
	if aimDir.X < 0 {
		playerFlip = true
	} else {
		playerFlip = false
	}

	// if face left, but moving right, reverse the animation
	if playerFlip && playerController.GetVelocity().X > 0 {
		playerController.ReverseAnimation()
	}

	// Use DrawLayeredPlayerSprite to draw both body and eyes
	deps.Renderer.DrawLayeredPlayerSprite(
		screen,
		playerSprite,
		playerEyeSprite,
		playerPos,
		playerEyePos,
		playerRot,
		spriteScale,
		playerFlip,
	)

	deps.Renderer.DrawAimLine(screen, playerPos, aimDir, 200)

	// Draw UI
	deps.Renderer.DrawUIText(fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), common.Vector2{X: 8, Y: 8}, color.RGBA{255, 255, 255, 255})

	deps.Renderer.EndFrame(screen)
}

func main() {
	setupMemoryProfiling()
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
		Renderer:     rendering.NewRenderer(rendering.DefaultRenderConfig(), cam),
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

func setupMemoryProfiling() {
	// Add memory profiling
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Force garbage collection more frequently during development
	go func() {
		for {
			time.Sleep(5 * time.Second)
			runtime.GC()
		}
	}()
}
