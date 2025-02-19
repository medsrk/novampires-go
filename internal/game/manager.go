// internal/game/manager.go
package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/debug"
	"novampires-go/internal/engine/input"
	"novampires-go/internal/game/config"
	"novampires-go/internal/game/player"
	"novampires-go/internal/game/rendering"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// Dependencies contains all external dependencies
type Dependencies struct {
	InputManager *input.Manager
	DebugManager *debug.Manager
	// Will expand with audio, assets, etc. as needed
}

// Game represents the main game state and logic
type Game struct {
	// Core dependencies
	dependencies Dependencies
	config       config.Config

	// Game state
	state GameState

	// Core systems
	renderer         *rendering.Renderer
	playerController *player.Controller

	// Debugging
	showDebug bool
}

// NewGame creates a new game instance
func NewGame(deps Dependencies, cfg config.Config) *Game {
	// Create renderer
	renderCfg := rendering.DefaultRenderConfig()
	renderer := rendering.NewRenderer(renderCfg)

	// Create player controller at center of screen
	startPos := common.Vector2{
		X: float64(cfg.Display.Width) / 2,
		Y: float64(cfg.Display.Height) / 2,
	}
	controllerCfg := player.DefaultConfig()
	// Apply auto-aim strength from gameplay settings
	//controllerCfg.AimAssistStrength = cfg.Gameplay.AutoAimStrength

	// Initialize player controller
	playerController := player.NewController(
		deps.InputManager,
		startPos,
		controllerCfg,
	)

	return &Game{
		dependencies:     deps,
		config:           cfg,
		state:            StateMainMenu,
		renderer:         renderer,
		playerController: playerController,
		showDebug:        cfg.Display.ShowDebugInfo,
	}
}

// Update handles game logic updates
func (g *Game) Update() error {
	// Update debug manager first
	g.dependencies.DebugManager.Update()

	// Check debug toggle
	if g.dependencies.InputManager.JustPressed(input.ActionToggleDebug) {
		g.showDebug = !g.showDebug
		g.dependencies.DebugManager.Toggle()
	}

	// Handle state-specific updates
	switch g.state {
	case StateMainMenu:
		g.updateMainMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		g.updatePaused()
	case StateGameOver:
		g.updateGameOver()
	}

	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMainMenu:
		g.drawMainMenu(screen)
	case StatePlaying:
		g.drawPlaying(screen)
	case StatePaused:
		g.drawPaused(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	}

	// Draw debug UI if enabled
	if g.showDebug {
		g.dependencies.DebugManager.Draw(screen)
	}
}

// Layout implements ebiten.Game
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Update ImGui display size if window is resized
	g.dependencies.DebugManager.SetDisplaySize(
		float32(outsideWidth),
		float32(outsideHeight),
	)
	return g.config.Display.Width, g.config.Display.Height
}

// State-specific update methods
func (g *Game) updateMainMenu() {
	// Check for game start
	if g.dependencies.InputManager.JustPressed(input.ActionInteract) {
		g.state = StatePlaying
	}
}

func (g *Game) updatePlaying() {
	// Update player controller with empty target list for now
	// Will be populated with enemies later
	g.playerController.Update([]common.TargetInfo{})

	// Check for pause
	if g.dependencies.InputManager.JustPressed(input.ActionMenu) {
		g.state = StatePaused
	}
}

func (g *Game) updatePaused() {
	// Check for unpause
	if g.dependencies.InputManager.JustPressed(input.ActionMenu) {
		g.state = StatePlaying
	}
}

func (g *Game) updateGameOver() {
	// Check for return to main menu
	if g.dependencies.InputManager.JustPressed(input.ActionInteract) {
		g.state = StateMainMenu
	}
}

// State-specific drawing methods
func (g *Game) drawMainMenu(screen *ebiten.Image) {
	// Draw main menu
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Draw game
	g.renderer.Draw(screen)
}

func (g *Game) drawPaused(screen *ebiten.Image) {
	// Draw pause menu
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	// Draw game over screen
}
