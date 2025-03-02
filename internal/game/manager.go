package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/camera"
	"novampires-go/internal/engine/debug"
	"novampires-go/internal/engine/input"
	"novampires-go/internal/engine/rendering"
	"novampires-go/internal/game/config"
	"novampires-go/internal/game/player"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// Scene represents a game scene (menu, level, etc)
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// Dependencies contains all external dependencies
type Dependencies struct {
	InputManager *input.Manager
	DebugManager *debug.Manager
	Camera       *camera.Camera
	Renderer     *rendering.Renderer
	Config       config.Config
	ScreenWidth  int
	ScreenHeight int
}

// Game represents the main game state and logic
type Game struct {
	deps Dependencies

	// Game state
	state       GameState
	activeScene Scene

	// Core systems
	playerController *player.Controller

	// State flags
	showDebug bool
}

// NewGame creates a new game instance
func NewGame(deps Dependencies) *Game {
	// First create the game instance with initial state
	g := &Game{
		deps:      deps,
		state:     StatePlaying, // Start in playing state for test scene
		showDebug: deps.Config.Display.ShowDebugInfo,
	}

	// Create player at center of screen
	startPos := common.Vector2{
		X: float64(deps.ScreenWidth) / 2,
		Y: float64(deps.ScreenHeight) / 2,
	}
	controllerCfg := player.DefaultConfig()
	g.playerController = player.NewController(deps.InputManager, startPos, controllerCfg)
	// Set player sprite
	playerSpritesheet, _, err := ebitenutil.NewImageFromFile("assets/doux.png")
	if err != nil {
		log.Printf("Failed to load player spritesheet: %v", err)
	}
	g.playerController.SetSpriteSheet(playerSpritesheet)
	g.playerController.SetScale(1.0)

	eyeSpritesheet, _, err := ebitenutil.NewImageFromFile("assets/doux-eyes.png")
	if err != nil {
		log.Printf("Failed to load eye spritesheet: %v", err)
	}
	g.playerController.SetEyeSpriteSheet(eyeSpritesheet)

	// Set up all debug windows
	deps.DebugManager.AddWindow(deps.InputManager.CreateDebugWindow())  // Input debug window
	keyBindEditor := input.NewKeyBindingEditorWindow(deps.InputManager) // Key binding editor
	deps.DebugManager.AddWindow(keyBindEditor)
	deps.DebugManager.AddWindow(deps.Camera.CreateDebugWindow()) // Camera debug window

	// Add player debug window last (after player controller is created)
	pdw := player.NewDebugWindow(deps.DebugManager, g.playerController)
	deps.DebugManager.AddWindow(pdw)

	// Tell camera to follow the player
	deps.Camera.SetTarget(g.playerController.GetPositionPtr())

	return g
}

// SetScene changes the active scene
func (g *Game) SetScene(scene Scene) {
	g.activeScene = scene
}

// Update handles game logic updates
func (g *Game) Update() error {
	// Update input manager
	g.deps.InputManager.Update()
	// Update debug manager first
	g.deps.DebugManager.Update()

	// Check debug toggle
	if g.deps.InputManager.JustPressed(common.ActionToggleDebug) {
		g.showDebug = !g.showDebug
		g.deps.DebugManager.Toggle()
	}

	// Update camera
	g.deps.Camera.Update()

	// Handle state-specific updates
	switch g.state {
	case StateMainMenu:
		return g.updateMainMenu()
	case StatePlaying:
		return g.updatePlaying()
	case StatePaused:
		return g.updatePaused()
	case StateGameOver:
		return g.updateGameOver()
	}

	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	if g.activeScene != nil {
		g.activeScene.Draw(screen)
	}

	// Draw debug UI if enabled
	if g.showDebug {
		g.deps.DebugManager.Draw(screen)
	}
}

// Layout implements ebiten.Game
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Update ImGui display size if window is resized
	g.deps.DebugManager.SetDisplaySize(float32(outsideWidth), float32(outsideHeight))
	return g.deps.ScreenWidth, g.deps.ScreenHeight
}

// State-specific update methods
func (g *Game) updateMainMenu() error {
	// Check for game start
	if g.deps.InputManager.JustPressed(common.ActionInteract) {
		g.state = StatePlaying
	}
	return nil
}

func (g *Game) updatePlaying() error {
	if g.activeScene != nil {
		return g.activeScene.Update()
	}

	// Check for pause
	if g.deps.InputManager.JustPressed(common.ActionMenu) {
		g.state = StatePaused
	}
	return nil
}

func (g *Game) updatePaused() error {
	// Check for unpause
	if g.deps.InputManager.JustPressed(common.ActionMenu) {
		g.state = StatePlaying
	}
	return nil
}

func (g *Game) updateGameOver() error {
	// Check for return to main menu
	if g.deps.InputManager.JustPressed(common.ActionInteract) {
		g.state = StateMainMenu
	}
	return nil
}

// Helper methods to access core systems
func (g *Game) GetPlayerController() *player.Controller {
	return g.playerController
}

func (g *Game) GetDependencies() *Dependencies {
	return &g.deps
}
