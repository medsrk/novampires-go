package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net/http"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/camera"
	"novampires-go/internal/engine/debug"
	"novampires-go/internal/engine/entity"
	"novampires-go/internal/engine/input"
	"novampires-go/internal/engine/rendering"
	"novampires-go/internal/game/config"
	"novampires-go/internal/game/scenes"
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

// Game represents the main game state and logic
type Game struct {
	// Core systems
	inputManager *input.Manager
	debugManager *debug.Manager
	camera       *camera.Camera
	renderer     *rendering.Renderer

	// Game state
	currentScene scene.TestScene
	showDebug    bool
}

func (g *Game) Update() error {
	// Update core systems
	g.inputManager.Update()

	// Check debug toggle
	if g.inputManager.JustPressed(common.ActionToggleDebug) {
		g.showDebug = !g.showDebug
		g.debugManager.Toggle()
	}

	// Only update debug manager if enabled
	if g.showDebug {
		g.debugManager.Update()
	}

	g.camera.Update()

	// Update current scene
	return g.currentScene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Begin frame
	g.renderer.BeginFrame(screen)

	// Draw current scene
	g.currentScene.Draw(screen)

	// Draw debug UI if enabled
	if g.showDebug {
		g.debugManager.Draw(screen)
	}

	// End frame
	g.renderer.EndFrame(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Update ImGui display size if window is resized
	g.debugManager.SetDisplaySize(float32(outsideWidth), float32(outsideHeight))
	return screenWidth, screenHeight
}

func main() {
	setupMemoryProfiling()

	// Set window properties
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("NoVampires Test Scene - Refactored")

	// Initialize core systems
	im := input.New()

	// Create config
	cfg := config.Default()

	// Create debug manager
	dm := debug.New(debug.Deps{InputManager: im})

	// Initialize camera with proper viewport size
	camConf := camera.DefaultConfig()
	camConf.ViewportSize = common.Vector2{X: float64(screenWidth), Y: float64(screenHeight)}
	cam := camera.NewWithConfig(camConf)
	cam.SetCenter(common.Vector2{X: float64(screenWidth) / 2, Y: float64(screenHeight) / 2})
	im.SetCamera(cam)

	// Create renderer
	renderConfig := rendering.DefaultRenderConfig()
	renderer := rendering.NewRenderer(renderConfig, cam)

	// Create renderer adapter for entity system
	rendererAdapter := entity.NewRendererAdapter(renderer)

	// Create game instance
	game := &Game{
		inputManager: im,
		debugManager: dm,
		camera:       cam,
		renderer:     renderer,
		showDebug:    cfg.Display.ShowDebugInfo,
	}

	// Add debug windows
	dm.AddWindow(im.CreateDebugWindow())
	keyBindEditor := input.NewKeyBindingEditorWindow(im)
	dm.AddWindow(keyBindEditor)
	dm.AddWindow(cam.CreateDebugWindow())

	// Create scene dependencies
	sceneDeps := scene.Dependencies{
		InputManager: im,
		Renderer:     rendererAdapter,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}

	// Create and set test scene
	game.currentScene = *scene.NewTestScene(sceneDeps)
	player := game.currentScene.GetPlayer()
	cam.SetTarget(&player.Position)

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
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
