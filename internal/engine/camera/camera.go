package camera

import (
	"math"
	"novampires-go/internal/common"
)

// Config holds camera configuration parameters
type Config struct {
	// How quickly the camera moves to its target position (0-1)
	Smoothing float64

	// Deadzone is the area around the target where the camera won't move
	Deadzone common.Rectangle

	// Bounds define the world boundaries the camera can't move beyond
	Bounds *common.Rectangle

	// The size of the viewport in screen coordinates
	ViewportSize common.Vector2
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Smoothing: 0.1,
		Deadzone: common.Rectangle{
			Size: common.Vector2{X: 10, Y: 10},
		},
		ViewportSize: common.Vector2{X: 1600, Y: 900},
	}
}

// Camera handles viewport management and coordinate transformations
type Camera struct {
	// Current position in world coordinates (top-left corner)
	pos common.Vector2

	// Current target position the camera is following
	target *common.Vector2

	// Configuration parameters
	config *Config

	// Current zoom level (1.0 = no zoom)
	zoom float64

	// Rotation in radians
	rotation float64
}

func New() *Camera {
	return NewWithConfig(DefaultConfig())
}

func NewWithConfig(config *Config) *Camera {
	return &Camera{
		config: config,
		zoom:   1.0,
	}
}

// Update handles camera movement and following behavior
func (c *Camera) Update() {
	if c.target == nil {
		return
	}

	targetCenter := *c.target

	// Calculate the desired camera position to center the target
	desiredPos := common.Vector2{
		X: targetCenter.X - c.config.ViewportSize.X/(2*c.zoom),
		Y: targetCenter.Y - c.config.ViewportSize.Y/(2*c.zoom),
	}

	// Apply smoothing
	c.pos.X += (desiredPos.X - c.pos.X) * c.config.Smoothing
	c.pos.Y += (desiredPos.Y - c.pos.Y) * c.config.Smoothing

	// Clamp to bounds if set
	if c.config.Bounds != nil {
		c.clampToBounds()
	}
}

// SetTarget sets the target position for the camera to follow
func (c *Camera) SetTarget(target *common.Vector2) {
	c.target = target
}

// GetTarget returns the current camera target position
func (c *Camera) GetTarget() *common.Vector2 {
	return c.target
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldPos common.Vector2) common.Vector2 {
	// Translate relative to camera
	x := worldPos.X - c.pos.X
	y := worldPos.Y - c.pos.Y

	// Apply rotation
	if c.rotation != 0 {
		cos := math.Cos(-c.rotation)
		sin := math.Sin(-c.rotation)
		oldX := x
		x = x*cos - y*sin
		y = oldX*sin + y*cos
	}

	// Apply zoom
	x *= c.zoom
	y *= c.zoom

	return common.Vector2{X: x, Y: y}
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenPos common.Vector2) common.Vector2 {
	// Unapply zoom
	x := screenPos.X / c.zoom
	y := screenPos.Y / c.zoom

	// Unapply rotation
	if c.rotation != 0 {
		cos := math.Cos(c.rotation)
		sin := math.Sin(c.rotation)
		oldX := x
		x = x*cos - y*sin
		y = oldX*sin + y*cos
	}

	// Translate back to world coordinates
	return common.Vector2{
		X: x + c.pos.X,
		Y: y + c.pos.Y,
	}
}

// GetViewport returns the current viewport rectangle in world coordinates
func (c *Camera) GetViewport() common.Rectangle {
	topLeft := c.pos
	size := common.Vector2{
		X: c.config.ViewportSize.X / c.zoom,
		Y: c.config.ViewportSize.Y / c.zoom,
	}
	return common.Rectangle{Pos: topLeft, Size: size}
}

// GetCenter returns the center position of the camera in world coordinates
func (c *Camera) GetCenter() common.Vector2 {
	return common.Vector2{
		X: c.pos.X + c.config.ViewportSize.X/(2*c.zoom),
		Y: c.pos.Y + c.config.ViewportSize.Y/(2*c.zoom),
	}
}

// SetCenter sets the camera position so that the given world coordinates are centered
func (c *Camera) SetCenter(worldPos common.Vector2) {
	c.pos = common.Vector2{
		X: worldPos.X - c.config.ViewportSize.X/(2*c.zoom),
		Y: worldPos.Y - c.config.ViewportSize.Y/(2*c.zoom),
	}
}

// SetZoom sets the camera zoom level
func (c *Camera) SetZoom(zoom float64) {
	// Store current center
	oldCenter := c.GetCenter()

	// Update zoom
	c.zoom = math.Max(0.1, zoom)

	// Restore center position
	c.SetCenter(oldCenter)
}

// SetRotation sets the camera rotation in radians
func (c *Camera) SetRotation(radians float64) {
	c.rotation = radians
}

// GetZoom returns the current zoom level
func (c *Camera) GetZoom() float64 {
	return c.zoom
}

// GetRotation returns the current rotation in radians
func (c *Camera) GetRotation() float64 {
	return c.rotation
}

// isInDeadzone checks if a point is within the camera's deadzone
func (c *Camera) isInDeadzone(pos common.Vector2) bool {
	// Create deadzone rectangle centered on camera
	center := c.GetCenter()
	deadzonePos := common.Vector2{
		X: center.X - c.config.Deadzone.Size.X/2,
		Y: center.Y - c.config.Deadzone.Size.Y/2,
	}
	deadzone := common.Rectangle{
		Pos:  deadzonePos,
		Size: c.config.Deadzone.Size,
	}
	return deadzone.Contains(pos)
}

// clampToBounds ensures the camera stays within the configured bounds
func (c *Camera) clampToBounds() {
	viewport := c.GetViewport()
	bounds := c.config.Bounds

	// Clamp horizontal position
	if viewport.Size.X < bounds.Size.X {
		if viewport.Pos.X < bounds.Pos.X {
			c.pos.X = bounds.Pos.X
		} else if viewport.Pos.X+viewport.Size.X > bounds.Pos.X+bounds.Size.X {
			c.pos.X = bounds.Pos.X + bounds.Size.X - viewport.Size.X
		}
	} else {
		// If viewport is larger than bounds, center it
		c.pos.X = bounds.Pos.X + (bounds.Size.X-viewport.Size.X)/2
	}

	// Clamp vertical position
	if viewport.Size.Y < bounds.Size.Y {
		if viewport.Pos.Y < bounds.Pos.Y {
			c.pos.Y = bounds.Pos.Y
		} else if viewport.Pos.Y+viewport.Size.Y > bounds.Pos.Y+bounds.Size.Y {
			c.pos.Y = bounds.Pos.Y + bounds.Size.Y - viewport.Size.Y
		}
	} else {
		// If viewport is larger than bounds, center it
		c.pos.Y = bounds.Pos.Y + (bounds.Size.Y-viewport.Size.Y)/2
	}
}
