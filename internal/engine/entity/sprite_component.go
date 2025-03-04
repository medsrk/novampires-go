package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/sprite"
	"time"
)

// SpriteComponent handles sprite rendering and animation for an entity
type SpriteComponent struct {
	// Sprite sheet and current sprite
	spriteSheet *ebiten.Image
	sprite      *ebiten.Image

	// Animation management
	animations     map[string]*sprite.Animation
	currentAnim    string
	lastUpdateTime time.Time

	// Rendering properties
	scale float64
	flipX bool

	// Secondary sprite layers (e.g., eyes)
	secondarySprite      *ebiten.Image
	secondarySpriteSheet *ebiten.Image
	secondaryOffset      common.Vector2
	secondaryController  AnimationController
}

// NewSpriteComponent creates a new sprite component
func NewSpriteComponent() *SpriteComponent {
	return &SpriteComponent{
		animations:     make(map[string]*sprite.Animation),
		scale:          1.0,
		lastUpdateTime: time.Now(),
	}
}

// Update updates the sprite and animations
func (s *SpriteComponent) Update(entity *Entity) {
	currentTime := time.Now()
	deltaTime := currentTime.Sub(s.lastUpdateTime)
	s.lastUpdateTime = currentTime

	// Update animation
	s.updateAnimation(deltaTime)

	// Update secondary animation controller if available
	if s.secondaryController != nil {
		s.secondaryController.Update(deltaTime)
	}
}

// updateAnimation updates the current animation frame
func (s *SpriteComponent) updateAnimation(deltaTime time.Duration) {
	// Skip if no animations or sprite sheet
	if s.spriteSheet == nil || len(s.animations) == 0 {
		return
	}

	// Update the current animation
	if anim, exists := s.animations[s.currentAnim]; exists && anim != nil {
		anim.Update(deltaTime)

		// Update the sprite frame based on current animation
		frame := anim.GetCurrentFrame()

		// Create a subimage from the sprite sheet
		rect := image.Rect(
			frame.SrcX,
			frame.SrcY,
			frame.SrcX+frame.SrcWidth,
			frame.SrcY+frame.SrcHeight,
		)

		s.sprite = s.spriteSheet.SubImage(rect).(*ebiten.Image)
	}
}

// SetSpriteSheet sets the main sprite sheet
func (s *SpriteComponent) SetSpriteSheet(sheet *ebiten.Image) {
	s.spriteSheet = sheet

	// Setup default animations if needed
	if s.spriteSheet != nil && len(s.animations) == 0 {
		s.setupDefaultAnimations()
	}
}

// SetupDefaultAnimations sets up standard animations if none exist
func (s *SpriteComponent) setupDefaultAnimations() {
	// This is a placeholder - specific animation setup should be done in specific entity implementations
	// Create a default idle animation if sheet is present but no animations are defined
	idleFrames := []sprite.FrameData{
		{SrcX: 0, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
	}
	s.animations["idle"] = sprite.NewAnimation(idleFrames, true)
	s.currentAnim = "idle"
}

// SetSecondarySpriteSheet sets a secondary sprite sheet (e.g., for eyes)
func (s *SpriteComponent) SetSecondarySpriteSheet(sheet *ebiten.Image, controller AnimationController) {
	s.secondarySpriteSheet = sheet
	s.secondaryController = controller
}

// SetSecondaryOffset sets the offset for the secondary sprite
func (s *SpriteComponent) SetSecondaryOffset(offset common.Vector2) {
	s.secondaryOffset = offset
}

// AddAnimation adds an animation to the sprite
func (s *SpriteComponent) AddAnimation(name string, frames []sprite.FrameData, loop bool) {
	s.animations[name] = sprite.NewAnimation(frames, loop)

	// If this is the first animation, set it as current
	if s.currentAnim == "" {
		s.currentAnim = name
		s.animations[name].Reset()
	}
}

// PlayAnimation changes the current animation
func (s *SpriteComponent) PlayAnimation(name string) {
	if anim, exists := s.animations[name]; exists && anim != nil {
		s.currentAnim = name
		anim.Reset()
	}
}

// GetSprite returns the current sprite frame
func (s *SpriteComponent) GetSprite() *ebiten.Image {
	return s.sprite
}

// GetSecondarySprite returns the secondary sprite (e.g., eyes)
func (s *SpriteComponent) GetSecondarySprite() *ebiten.Image {
	if s.secondaryController != nil {
		return s.secondaryController.GetSprite()
	}
	return s.secondarySprite
}

// GetSecondaryOffset returns the secondary sprite offset
func (s *SpriteComponent) GetSecondaryOffset() common.Vector2 {
	return s.secondaryOffset
}

// SetScale sets the sprite scale
func (s *SpriteComponent) SetScale(scale float64) {
	s.scale = scale
}

// GetScale returns the sprite scale
func (s *SpriteComponent) GetScale() float64 {
	return s.scale
}

// SetFlipX sets whether the sprite should be flipped horizontally
func (s *SpriteComponent) SetFlipX(flip bool) {
	s.flipX = flip

	// Also set flip state on secondary controller if available
	if s.secondaryController != nil {
		s.secondaryController.SetFlipX(flip)
	}
}

// GetFlipX returns whether the sprite is flipped horizontally
func (s *SpriteComponent) GetFlipX() bool {
	return s.flipX
}

// GetCurrentAnimation returns the name of the current animation
func (s *SpriteComponent) GetCurrentAnimation() string {
	return s.currentAnim
}

func (s *SpriteComponent) GetCurrentFrame() int {
	if anim, exists := s.animations[s.currentAnim]; exists && anim != nil {
		return anim.GetCurrentFrameInt()
	}
	return 0
}

// ReverseAnimation reverses the current animation
func (s *SpriteComponent) ReverseAnimation() {
	if anim, exists := s.animations[s.currentAnim]; exists && anim != nil {
		anim.Reverse()
	}
}

// IsReversed returns whether the current animation is reversed
func (s *SpriteComponent) IsReversed() bool {
	if anim, exists := s.animations[s.currentAnim]; exists && anim != nil {
		return anim.IsReversed()
	}
	return false
}

// Draw draws the sprite
func (s *SpriteComponent) Draw(screen *ebiten.Image, renderer Renderer, entity *Entity) {
	if s.sprite == nil {
		return
	}

	// Draw main sprite and secondary sprite if available
	if s.secondaryController != nil || s.secondarySprite != nil {
		secondarySprite := s.GetSecondarySprite()
		renderer.DrawLayeredSprite(
			screen,
			s.sprite,
			secondarySprite,
			entity.Position,
			s.secondaryOffset,
			entity.Rotation,
			s.scale,
			s.flipX,
		)
	} else {
		renderer.DrawSprite(
			screen,
			s.sprite,
			entity.Position,
			entity.Rotation,
			s.scale,
			s.flipX)
	}
}
