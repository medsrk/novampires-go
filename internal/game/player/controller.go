package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/sprite"
	"time"
)

type Controller struct {
	inputManager common.InputProvider

	pos      common.Vector2
	vel      common.Vector2
	rotation float64

	config Config

	autoAim      bool
	usingGamepad bool
	lastMouseX   int
	lastMouseY   int
	lastAimDx    float64
	lastAimDy    float64

	// Sprite handling
	sprite      *ebiten.Image
	spriteSheet *ebiten.Image
	spriteScale float64
	animations  map[string]*sprite.Animation
	currentAnim string

	// eye sprite handling
	eyeController *EyeController

	lastUpdateTime time.Time
	flipX          bool
}

type Config struct {
	MaxSpeed      float64
	Acceleration  float64
	Deceleration  float64
	RotationSpeed float64

	Range float64
}

func DefaultConfig() Config {
	return Config{
		MaxSpeed:      5.0,
		Acceleration:  1.0,
		Deceleration:  0.5,
		RotationSpeed: 0.15,

		Range: 400.0,
	}
}

func NewController(inputManager common.InputProvider, initialPos common.Vector2, config Config) *Controller {
	return &Controller{
		inputManager:   inputManager,
		pos:            initialPos,
		vel:            common.Vector2{},
		rotation:       0,
		config:         config,
		autoAim:        true,
		animations:     make(map[string]*sprite.Animation),
		eyeController:  NewEyeController(),
		lastUpdateTime: time.Now(),
	}
}

// SetSpriteSheet sets the sprite sheet for the player
func (c *Controller) SetSpriteSheet(sheet *ebiten.Image) {
	c.spriteSheet = sheet

	// Safety check
	if sheet == nil {
		return
	}

	// Initialize animations map if it doesn't exist
	if c.animations == nil {
		c.animations = make(map[string]*sprite.Animation)
	}

	// Initialize animations if not already set up
	// Create idle animation (first row of spritesheet)
	idleFrames := []sprite.FrameData{
		{SrcX: 0, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
		{SrcX: 96, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
		{SrcX: 192, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
		{SrcX: 288, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 150},
	}
	c.animations["idle"] = sprite.NewAnimation(idleFrames, true)

	// Create walk animation (second row of spritesheet)
	walkFrames := []sprite.FrameData{
		{SrcX: 384, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 480, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 576, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 672, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 768, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
		{SrcX: 864, SrcY: 0, SrcWidth: 96, SrcHeight: 96, Duration: 100},
	}
	c.animations["walk"] = sprite.NewAnimation(walkFrames, true)

	// Start with idle animation
	c.currentAnim = "idle"
	if c.animations["idle"] != nil {
		c.animations["idle"].Reset()
	}

	// Initialize the sprite with the first frame
	if c.animations["idle"] != nil {
		frame := c.animations["idle"].GetCurrentFrame()
		rect := image.Rect(
			frame.SrcX,
			frame.SrcY,
			frame.SrcX+frame.SrcWidth,
			frame.SrcY+frame.SrcHeight,
		)
		c.sprite = c.spriteSheet.SubImage(rect).(*ebiten.Image)
	} else {
		// Fallback to full spritesheet if for some reason idle animation is nil
		c.sprite = c.spriteSheet
	}
}

func (c *Controller) SetEyeSpriteSheet(sheet *ebiten.Image) {
	c.eyeController.SetSpriteSheet(sheet)
}

func (c *Controller) SetSpriteSheetScale(scale float64) {
	c.spriteScale = scale
}

// Update method now handles sprite animation updates and passes mouse direction to eye controller
func (c *Controller) Update(targets []common.TargetInfo) {
	// Calculate delta time for animation updates
	currentTime := time.Now()
	deltaTime := currentTime.Sub(c.lastUpdateTime)
	c.lastUpdateTime = currentTime

	// Update movement and aiming
	c.updateMovement()
	c.updateAiming(targets)

	// Get aim direction for eye direction
	aimDirection := c.GetAimDirection()

	// Update eye look direction
	c.eyeController.UpdateLookDirection(aimDirection)

	// Update eye animations (blink)
	c.eyeController.Update(deltaTime)

	// Update character animations
	c.updateAnimation(deltaTime)

	// Update auto-aim state
	if c.inputManager.JustPressed(common.ActionAutoAttack) {
		c.autoAim = !c.autoAim
	}
	if c.inputManager.JustPressed(common.ActionUseAbility1) {
		c.eyeController.TriggerBlink()
	}
}

func (c *Controller) updateMovement() {
	dx, dy := c.inputManager.GetMovementVector()
	inputVec := common.Vector2{X: dx, Y: dy}
	inputMagnitude := inputVec.Magnitude()

	// Store previous velocity for animation state changes
	//prevVel := c.vel

	if inputMagnitude > 0 {
		// Scale max speed by input magnitude
		targetSpeed := c.config.MaxSpeed * inputMagnitude

		c.vel = c.vel.Add(inputVec.Scale(c.config.Acceleration))
		currentSpeed := c.vel.Magnitude()

		// Cap at the scaled max speed
		if currentSpeed > targetSpeed {
			c.vel = c.vel.Normalized().Scale(targetSpeed)
		}
	} else {
		currentSpeed := c.vel.Magnitude()
		if currentSpeed > 0 {
			newSpeed := math.Max(0, currentSpeed-c.config.Deceleration)
			if newSpeed > 0 {
				c.vel = c.vel.Normalized().Scale(newSpeed)
			} else {
				c.vel = common.Vector2{}
			}
		}
	}

	c.pos = c.pos.Add(c.vel)
}

func (c *Controller) GetAimVector() (float64, float64) {
	// Get current mouse position in WORLD coordinates
	mx, my := c.inputManager.GetMousePositionWorld()

	// Try gamepad first - this part remains unchanged
	if dx, dy, ok := c.inputManager.GetGamepadAim(); ok {
		c.usingGamepad = true
		c.lastAimDx, c.lastAimDy = dx, dy // Store gamepad aim vector

		return dx, dy
	}

	// Check for mouse movement
	if mx != c.lastMouseX || my != c.lastMouseY {
		c.usingGamepad = false
		c.lastMouseX, c.lastMouseY = mx, my

		// Calculate the direction vector from player to mouse
		mousePos := common.Vector2{X: float64(mx), Y: float64(my)}

		// Get the vector from player to mouse
		aimDirection := mousePos.Sub(c.pos)

		// Normalize the vector to get a direction
		if aimDirection.MagnitudeSquared() > 0 {
			aimDirection = aimDirection.Normalized()
		} else {
			// If mouse is exactly on player position, keep last direction
			aimDirection = common.Vector2{X: c.lastAimDx, Y: c.lastAimDy}
		}

		c.lastAimDx, c.lastAimDy = aimDirection.X, aimDirection.Y // Store NORMALIZED direction
		return c.lastAimDx, c.lastAimDy
	}

	// No new input, return last used aim (mouse or gamepad vector)
	return c.lastAimDx, c.lastAimDy
}

func (c *Controller) updateAiming(targets []common.TargetInfo) {
	if c.autoAim && len(targets) > 0 {
		// Auto-aim logic stays the same
		var closestTarget *common.TargetInfo
		closestDistSq := c.config.Range * c.config.Range

		for i, target := range targets {
			delta := target.Pos.Sub(c.pos)
			distSq := delta.MagnitudeSquared()

			if distSq < closestDistSq {
				closestDistSq = distSq
				closestTarget = &targets[i]
			}
		}

		if closestTarget != nil {
			aimDirection := closestTarget.Pos.Sub(c.pos)
			targetRotation := math.Atan2(aimDirection.Y, aimDirection.X)
			angleDiff := common.NormalizeAngle(targetRotation - c.rotation)
			rotationStep := angleDiff * c.config.RotationSpeed
			c.rotation = common.NormalizeAngle(c.rotation + rotationStep)
		}
	} else {
		// Manual aim using input manager's aim vector
		dx, dy := c.GetAimVector()

		if dx != 0 || dy != 0 {
			c.rotation = math.Atan2(dy, dx)
		}
	}
}

// updateAnimation handles sprite animation based on player state
func (c *Controller) updateAnimation(deltaTime time.Duration) {
	// Safety check - don't attempt to update animations if they're not initialized
	if c.spriteSheet == nil || len(c.animations) == 0 || c.animations["idle"] == nil {
		return // No animations to update
	}

	// Determine animation state based on player movement
	newAnim := c.currentAnim

	// If there's no current animation, default to idle
	if newAnim == "" {
		newAnim = "idle"
	}

	// If player is attacking, prioritize attack animation
	if c.inputManager.IsPressed(common.ActionAutoAttack) {
		attackAnim, hasAttack := c.animations["attack"]
		if hasAttack && (newAnim != "attack" || attackAnim.IsFinished()) {
			newAnim = "attack"
		}
	} else if c.vel.MagnitudeSquared() > 0.1 {
		// Player is moving
		if _, hasWalk := c.animations["walk"]; hasWalk {
			newAnim = "walk"
		}
	} else {
		// Player is idle
		if _, hasIdle := c.animations["idle"]; hasIdle {
			newAnim = "idle"
		}
	}

	// Change animation if state changed
	if newAnim != c.currentAnim {
		c.currentAnim = newAnim
		if anim, exists := c.animations[newAnim]; exists && anim != nil {
			anim.Reset()
		}
	}

	// Update the current animation
	if anim, exists := c.animations[c.currentAnim]; exists && anim != nil {
		anim.Update(deltaTime)
	}

	// Update the sprite frame based on current animation
	if anim, exists := c.animations[c.currentAnim]; exists && anim != nil {
		frame := anim.GetCurrentFrame()

		// Create a subimage from the sprite sheet
		rect := image.Rect(
			frame.SrcX,
			frame.SrcY,
			frame.SrcX+frame.SrcWidth,
			frame.SrcY+frame.SrcHeight,
		)

		c.sprite = c.spriteSheet.SubImage(rect).(*ebiten.Image)
	}
}

func (c *Controller) GetPosition() common.Vector2 {
	return c.pos
}

func (c *Controller) GetPositionPtr() *common.Vector2 {
	return &c.pos
}

func (c *Controller) GetRotation() float64 {
	return c.rotation
}

// GetVelocity returns the current velocity vector
func (c *Controller) GetVelocity() common.Vector2 {
	return c.vel
}

// GetAimDirection returns the direction the player is aiming
func (c *Controller) GetAimDirection() common.Vector2 {
	return common.Vector2{X: math.Cos(c.rotation), Y: math.Sin(c.rotation)}.Normalized()
}

func (c *Controller) SetPosition(pos common.Vector2) {
	c.pos = pos
}

// GetSprite returns the current sprite frame
func (c *Controller) GetSprite() *ebiten.Image {
	return c.sprite
}

// GetFlipX returns whether the sprite should be flipped horizontally
func (c *Controller) GetFlipX() bool {
	return c.flipX
}

// GetCurrentAnimation returns the name of the current animation
func (c *Controller) GetCurrentAnimation() string {
	return c.currentAnim
}

// PlayAnimation forces the player to play a specific animation
// PlayAnimation forces the player to play a specific animation
func (c *Controller) PlayAnimation(animName string) {
	// Safety check - make sure animation exists and is initialized
	if anim, exists := c.animations[animName]; exists && anim != nil {
		c.currentAnim = animName
		anim.Reset()
	}
}

// GetCurrentFrame returns the current frame index of the active animation
func (c *Controller) GetCurrentFrame() int {
	if anim, exists := c.animations[c.currentAnim]; exists && anim != nil {
		return anim.GetCurrentFrameInt()
	}
	return 0
}

// SetFrameRate changes the frame duration for all animations
func (c *Controller) SetFrameRate(frameRate float32) {
	// Convert frameRate (fps) to frame duration (ms)
	frameDuration := int(1000.0 / frameRate)

	// Update all animations
	for _, anim := range c.animations {
		for i := range anim.Frames {
			anim.Frames[i].Duration = frameDuration
		}
	}
}

// SetScale sets the player sprite scale
func (c *Controller) SetScale(scale float32) {
	c.spriteScale = float64(scale)
}

// SetFlipX explicitly sets the horizontal flip state
func (c *Controller) SetFlipX(flip bool) {
	c.flipX = flip
}

// GetAnimationDuration returns the total duration of the current animation in ms
func (c *Controller) GetAnimationDuration() float32 {
	if anim, exists := c.animations[c.currentAnim]; exists {
		totalDuration := 0
		for _, frame := range anim.Frames {
			totalDuration += frame.Duration
		}
		return float32(totalDuration)
	}
	return 0
}

// GetScale returns the current sprite scale
func (c *Controller) GetScale() float64 {
	return c.spriteScale
}

func (c *Controller) ReverseAnimation() {
	if c.currentAnim == "walk" {
		c.animations["walk"].Reverse()
	}
}

func (c *Controller) IsReversed() bool {
	if anim, exists := c.animations[c.currentAnim]; exists {
		return anim.IsReversed()
	}
	return false
}

func (c *Controller) GetEyeSprite() *ebiten.Image {
	return c.eyeController.GetSprite()
}

func (c *Controller) GetEyePosition() common.Vector2 {
	eyePos := c.eyeController.GetPosition()

	switch c.currentAnim {
	case "idle":
		if c.animations["idle"].GetCurrentFrameInt() == 2 {
			eyePos.Y += 4
		}
	case "walk":
		if c.animations["walk"].GetCurrentFrameInt() == 1 ||
			c.animations["walk"].GetCurrentFrameInt() == 4 {
			eyePos.Y += 4
		}
	}
	return eyePos
}

func (c *Controller) TriggerEyeBlink() {
	c.eyeController.TriggerBlink()
}

func (c *Controller) Direction() common.Vector2 {
	return common.Vector2{X: math.Cos(c.rotation), Y: math.Sin(c.rotation)}
}

func (c *Controller) DirectionString() string {

	dir := c.Direction()

	if dir.X > 0.5 {
		return "right"
	} else if dir.X < -0.5 {
		return "left"
	} else if dir.Y > 0.5 {
		return "down"
	} else if dir.Y < -0.5 {
		return "up"
	}
	return "idle"
}
