package player

import (
	"math"
	"novampires-go/internal/common"
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
		inputManager: inputManager,
		pos:          initialPos,
		vel:          common.Vector2{},
		rotation:     0,
		config:       config,
		autoAim:      true,
	}
}

func (c *Controller) Update(targets []common.TargetInfo) {
	c.updateMovement()
	c.updateAiming(targets)

	// Update auto-aim state
	if c.inputManager.JustPressed(common.ActionAutoAttack) {
		c.autoAim = !c.autoAim
	}
}

func (c *Controller) updateMovement() {
	dx, dy := c.inputManager.GetMovementVector()
	inputVec := common.Vector2{X: dx, Y: dy}
	inputMagnitude := inputVec.Magnitude()

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

func (c *Controller) GetPosition() common.Vector2 {
	return c.pos
}

func (c *Controller) GetPositionPtr() *common.Vector2 {
	return &c.pos
}

func (c *Controller) GetRotation() float64 {
	return c.rotation
}

// GetAimDirection returns the direction the player is aiming
func (c *Controller) GetAimDirection() common.Vector2 {
	return common.Vector2{X: math.Cos(c.rotation), Y: math.Sin(c.rotation)}.Normalized()
}

func (c *Controller) SetPosition(pos common.Vector2) {
	c.pos = pos
}
