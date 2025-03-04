// internal/entity/player_input.go
package entity

import (
	"math"
	"novampires-go/internal/common"
)

// PlayerInput processes player input and updates entity state
type PlayerInput struct {
	inputManager common.InputProvider

	config PlayerInputConfig

	// State tracking
	usingGamepad bool
	lastMouseX   int
	lastMouseY   int
	lastAimDx    float64
	lastAimDy    float64

	// Target tracking for auto-aim
	currentTargets []common.TargetInfo
	autoAim        bool

	entity *Entity
}

// PlayerInputConfig contains configuration for player input
type PlayerInputConfig struct {
	MaxSpeed      float64
	Acceleration  float64
	Deceleration  float64
	RotationSpeed float64
	AutoAimRange  float64
}

// DefaultPlayerInputConfig returns default player input configuration
func DefaultPlayerInputConfig() PlayerInputConfig {
	return PlayerInputConfig{
		MaxSpeed:      5.0,
		Acceleration:  1.0,
		Deceleration:  0.5,
		RotationSpeed: 0.15,
		AutoAimRange:  400.0,
	}
}

// NewPlayerInput creates a new player input component
func NewPlayerInput(inputManager common.InputProvider, config PlayerInputConfig, entity *Entity) *PlayerInput {
	return &PlayerInput{
		inputManager: inputManager,
		config:       config,
		autoAim:      true,
		entity:       entity,
	}
}

// ProcessInput processes player input and updates entity state
func (p *PlayerInput) ProcessInput(entity *Entity) {
	p.updateMovement(entity)
	p.updateAiming(entity)

	// Update auto-aim state
	if p.inputManager.JustPressed(common.ActionAutoAttack) {
		p.autoAim = !p.autoAim
	}

	// Check if sprite component exists
	sprite := entity.GetSprite()
	if sprite == nil {
		return
	}

	// Update animation based on movement
	p.updateAnimation(entity, sprite)
}

// UpdateTargets updates the list of potential targets for auto-aim
func (p *PlayerInput) UpdateTargets(targets []common.TargetInfo) {
	p.currentTargets = targets
}

// updateMovement handles player movement input
func (p *PlayerInput) updateMovement(entity *Entity) {
	dx, dy := p.inputManager.GetMovementVector()
	inputVec := common.Vector2{X: dx, Y: dy}
	inputMagnitude := inputVec.Magnitude()

	velocity := entity.GetVelocity()

	if inputMagnitude > 0 {
		// Scale max speed by input magnitude
		targetSpeed := p.config.MaxSpeed * inputMagnitude

		velocity = velocity.Add(inputVec.Scale(p.config.Acceleration))
		currentSpeed := velocity.Magnitude()

		// Cap at the scaled max speed
		if currentSpeed > targetSpeed {
			velocity = velocity.Normalized().Scale(targetSpeed)
		}
	} else {
		currentSpeed := velocity.Magnitude()
		if currentSpeed > 0 {
			newSpeed := math.Max(0, currentSpeed-p.config.Deceleration)
			if newSpeed > 0 {
				velocity = velocity.Normalized().Scale(newSpeed)
			} else {
				velocity = common.Vector2{}
			}
		}
	}

	entity.SetVelocity(velocity)
}

// GetAimVector returns the normalized aim vector
func (p *PlayerInput) GetAimVector() (float64, float64) {
	// Get current mouse position in WORLD coordinates
	mx, my := p.inputManager.GetMousePositionWorld()

	// Try gamepad first
	if dx, dy, ok := p.inputManager.GetGamepadAim(); ok {
		p.usingGamepad = true
		p.lastAimDx, p.lastAimDy = dx, dy // Store gamepad aim vector

		return dx, dy
	}

	// Check for mouse movement
	if mx != p.lastMouseX || my != p.lastMouseY {
		p.usingGamepad = false
		p.lastMouseX, p.lastMouseY = mx, my

		// Get the entity's current position
		mousePos := common.Vector2{X: float64(mx), Y: float64(my)}

		entityPos := p.entity.GetPosition()

		// Calculate the direction vector from player to mouse
		aimDirection := mousePos.Sub(entityPos)

		// Normalize the vector to get a direction
		if aimDirection.MagnitudeSquared() > 0 {
			aimDirection = aimDirection.Normalized()
		} else {
			// If mouse is exactly on player position, keep last direction
			aimDirection = common.Vector2{X: p.lastAimDx, Y: p.lastAimDy}
		}

		p.lastAimDx, p.lastAimDy = aimDirection.X, aimDirection.Y // Store NORMALIZED direction
		return p.lastAimDx, p.lastAimDy
	}

	// No new input, return last used aim (mouse or gamepad vector)
	return p.lastAimDx, p.lastAimDy
}

// GetAimDirection returns the current aiming direction
func (p *PlayerInput) GetAimDirection() common.Vector2 {
	dx, dy := p.GetAimVector()
	return common.Vector2{X: dx, Y: dy}
}

// updateAiming handles player aiming input
func (p *PlayerInput) updateAiming(entity *Entity) {
	if p.autoAim && len(p.currentTargets) > 0 {
		// Auto-aim logic
		var closestTarget *common.TargetInfo
		closestDistSq := p.config.AutoAimRange * p.config.AutoAimRange
		entityPos := entity.GetPosition()

		for i, target := range p.currentTargets {
			delta := target.Pos.Sub(entityPos)
			distSq := delta.MagnitudeSquared()

			if distSq < closestDistSq {
				closestDistSq = distSq
				closestTarget = &p.currentTargets[i]
			}
		}

		if closestTarget != nil {
			aimDirection := closestTarget.Pos.Sub(entityPos)
			targetRotation := math.Atan2(aimDirection.Y, aimDirection.X)
			currentRotation := entity.GetRotation()
			angleDiff := common.NormalizeAngle(targetRotation - currentRotation)
			rotationStep := angleDiff * p.config.RotationSpeed
			entity.SetRotation(common.NormalizeAngle(currentRotation + rotationStep))
		}
	} else {
		// Manual aim using input manager's aim vector
		dx, dy := p.GetAimVector()

		if dx != 0 || dy != 0 {
			entity.SetRotation(math.Atan2(dy, dx))
		}
	}
}

// updateAnimation updates the entity's animation based on its movement
func (p *PlayerInput) updateAnimation(entity *Entity, sprite *SpriteComponent) {
	velocity := entity.GetVelocity()
	aimDirection := p.GetAimDirection()

	// Update sprite direction based on aim
	if aimDirection.X < 0 {
		sprite.SetFlipX(true)
	} else if aimDirection.X > 0 {
		sprite.SetFlipX(false)
	}

	// Manage animation reversal when needed
	if sprite.GetFlipX() && velocity.X > 0 {
		sprite.ReverseAnimation()
	}

	// Set animation based on movement
	currentAnim := sprite.GetCurrentAnimation()

	// Handle animation state changes
	if velocity.MagnitudeSquared() > 0.1 {
		// Player is moving
		if currentAnim != "walk" {
			sprite.PlayAnimation("walk")
		}
	} else {
		// Player is idle
		if currentAnim != "idle" {
			sprite.PlayAnimation("idle")
		}
	}
}

// IsUsingGamepad returns whether the player is using a gamepad
func (p *PlayerInput) IsUsingGamepad() bool {
	return p.usingGamepad
}

// IsAutoAimEnabled returns whether auto-aim is enabled
func (p *PlayerInput) IsAutoAimEnabled() bool {
	return p.autoAim
}

// SetAutoAim sets whether auto-aim is enabled
func (p *PlayerInput) SetAutoAim(enabled bool) {
	p.autoAim = enabled
}
