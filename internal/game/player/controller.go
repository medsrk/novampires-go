package player

import (
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/input"
)

type Controller struct {
	inputManager input.InputProvider

	pos      common.Vector2
	vel      common.Vector2
	rotation float64

	config Config

	autoAim bool
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
		MaxSpeed:      3.0,
		Acceleration:  0.2,
		Deceleration:  0.1,
		RotationSpeed: 0.15,

		Range: 400.0,
	}
}

func NewController(inputManager input.InputProvider, initialPos common.Vector2, config Config) *Controller {
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
	if c.inputManager.JustPressed(input.ActionAutoAttack) {
		c.autoAim = !c.autoAim
	}
}

func (c *Controller) updateMovement() {
	dx, dy := c.inputManager.GetMovementVector()
	inputVec := common.Vector2{X: dx, Y: dy}

	if inputVec.Magnitude() > 0 {
		c.vel = c.vel.Add(inputVec.Scale(c.config.Acceleration))

		if c.vel.Magnitude() > c.config.MaxSpeed {
			c.vel = c.vel.Normalized().Scale(c.config.MaxSpeed)
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

func (c *Controller) updateAiming(targets []common.TargetInfo) {
	if len(targets) == 0 {
		return
	}
	if !c.autoAim {
		dx, dy := c.inputManager.GetAimVector()
		aimVec := common.Vector2{X: dx, Y: dy}
		if aimVec.Magnitude() > 0 {
			c.rotation = math.Atan2(aimVec.Y, aimVec.X)
		}

	} else {

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

		// Aim at closest target with aim assist
		if closestTarget != nil {
			aimDirection := closestTarget.Pos.Sub(c.pos)
			targetRotation := math.Atan2(aimDirection.Y, aimDirection.X)

			// Gradual rotation towards target (aim assist)
			angleDiff := common.NormalizeAngle(targetRotation - c.rotation)
			rotationStep := angleDiff * c.config.RotationSpeed
			c.rotation = common.NormalizeAngle(c.rotation + rotationStep)
		}
	}
}

func (c *Controller) GetPosition() common.Vector2 {
	return c.pos
}

func (c *Controller) GetRotation() float64 {
	return c.rotation
}

// GetAimDirection returns the direction the player is aiming
func (c *Controller) GetAimDirection() common.Vector2 {
	return common.Vector2{X: math.Cos(c.rotation), Y: math.Sin(c.rotation)}
}

func (c *Controller) SetPosition(pos common.Vector2) {
	c.pos = pos
}
