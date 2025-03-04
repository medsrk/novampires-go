package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/entity"
	"novampires-go/internal/engine/sprite"
)

// Player represents the player character built on the entity system
type Player struct {
	*entity.Entity
	input         *entity.PlayerInput
	eyeController *entity.EyeController
}

// NewPlayer creates a new player instance
func NewPlayer(inputManager common.InputProvider, initialPos common.Vector2) *Player {
	// Create base entity
	baseEntity := entity.NewEntity(1, initialPos)

	// Create player input component with default config
	playerInput := entity.NewPlayerInput(inputManager, entity.DefaultPlayerInputConfig(), baseEntity)
	baseEntity.SetInput(playerInput)

	// Create sprite component
	spriteComponent := entity.NewSpriteComponent()
	baseEntity.SetSprite(spriteComponent)

	// Create eye controller
	eyeController := entity.NewEyeController()

	// Create player instance
	player := &Player{
		Entity:        baseEntity,
		input:         playerInput,
		eyeController: eyeController,
	}

	// Load player sprites
	player.loadSprites()

	return player
}

// loadSprites loads the player's sprites and sets up animations
func (p *Player) loadSprites() {
	// Load body spritesheet
	playerSpritesheet, _, err := ebitenutil.NewImageFromFile("assets/doux.png")
	if err != nil {
		log.Printf("Failed to load player spritesheet: %v", err)
		return
	}

	// Get sprite component
	spriteComponent := p.GetSprite()
	if spriteComponent == nil {
		return
	}

	// Set main sprite sheet
	spriteComponent.SetSpriteSheet(playerSpritesheet)

	// Set up animations
	// Create idle animation (first row of spritesheet)
	idleFrames := []sprite.FrameData{
		{SrcX: 0, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
		{SrcX: 96, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
		{SrcX: 192, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
		{SrcX: 288, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
	}
	spriteComponent.AddAnimation("idle", idleFrames, true)

	// Create walk animation (second row of spritesheet)
	walkFrames := []sprite.FrameData{
		{SrcX: 384, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 480, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 576, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 672, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 768, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 864, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
	}
	spriteComponent.AddAnimation("walk", walkFrames, true)

	// Set default animation
	spriteComponent.PlayAnimation("idle")

	// Load eye spritesheet
	eyeSpritesheet, _, err := ebitenutil.NewImageFromFile("assets/doux-eyes.png")
	if err != nil {
		log.Printf("Failed to load eye spritesheet: %v", err)
	} else {
		// Set up eye controller
		p.eyeController.SetSpriteSheet(eyeSpritesheet)
		p.eyeController.SetPosition(common.Vector2{X: 0, Y: 0})

		// Link eye controller to sprite component
		spriteComponent.SetSecondarySpriteSheet(eyeSpritesheet, p.eyeController)
		spriteComponent.SetSecondaryOffset(common.Vector2{X: 0, Y: 0})
	}

	// Set scale
	spriteComponent.SetScale(1.0)
}

// Update updates the player state
func (p *Player) Update(targets []common.TargetInfo) {
	// Update targets in player input
	p.input.UpdateTargets(targets)

	// Update eye direction based on aim direction
	p.eyeController.UpdateLookDirection(p.input.GetAimDirection())

	// Update base entity
	p.Entity.Update()
}

// Draw draws the player
func (p *Player) Draw(screen *ebiten.Image, renderer entity.Renderer) {
	// Draw entity (will use sprite component if available)
	p.Entity.Draw(screen, renderer)

	// Draw aim line if needed
	renderer.DrawAimLine(screen, p.GetPosition(), p.input.GetAimDirection(), 200)
}

// TriggerEyeBlink triggers a blink animation
func (p *Player) TriggerEyeBlink() {
	p.eyeController.TriggerBlink()
}

// GetAimDirection returns the direction the player is aiming
func (p *Player) GetAimDirection() common.Vector2 {
	return p.input.GetAimDirection()
}

// IsAutoAimEnabled returns whether auto-aim is enabled
func (p *Player) IsAutoAimEnabled() bool {
	return p.input.IsAutoAimEnabled()
}

// SetAutoAim sets whether auto-aim is enabled
func (p *Player) SetAutoAim(enabled bool) {
	p.input.SetAutoAim(enabled)
}

func (p *Player) GetEyePosition() common.Vector2 {
	// Get base eye position from eye controller
	currentAnim := p.GetSprite().GetCurrentAnimation()
	currentFrame := p.GetSprite().GetCurrentFrame()
	eyePos := p.eyeController.GetPosition(currentAnim, currentFrame)

	// The eye controller's GetPosition now handles all animation-specific adjustments
	return eyePos
}
