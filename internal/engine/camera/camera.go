package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	// Current position in world coordinates (center of view)
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

	// Apply smoothing to move toward the target
	c.pos.X += (targetCenter.X - c.pos.X) * c.config.Smoothing
	c.pos.Y += (targetCenter.Y - c.pos.Y) * c.config.Smoothing

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

// GetTransform returns the transformation matrix for rendering
// This method is kept for compatibility but the actual transform
// calculation is now done in the renderer's EndFrame method
func (c *Camera) GetTransform() ebiten.GeoM {
	m := ebiten.GeoM{}

	// The specific transform is now applied in the renderer
	// to ensure proper coordination with the world buffer

	return m
}

// ScreenToWorld converts screen coordinates to world coordinates
//func (c *Camera) ScreenToWorld(screenPos common.Vector2) common.Vector2 {
//	// Get a copy of our transform and invert it
//	inv := c.GetTransform()
//	if !inv.IsInvertible() {
//		return common.Vector2{} // Return origin if not invertible
//	}
//
//	inv.Invert()
//
//	// Apply inverse transform to screen coordinates
//	worldX, worldY := inv.Apply(screenPos.X, screenPos.Y)
//	return common.Vector2{X: worldX, Y: worldY}
//}

func (c *Camera) ScreenToWorld(screenPos common.Vector2) common.Vector2 {
	// Important: This method needs to reverse exactly what happens in the renderer's EndFrame

	// 1. Get screen dimensions
	screenWidth, screenHeight := c.config.ViewportSize.X, c.config.ViewportSize.Y

	// 2. Create a transformation matrix that converts screen coords to world coords
	m := ebiten.GeoM{}

	// 3. First, we need to move from screen position to the center
	m.Translate(-screenWidth/2, -screenHeight/2)

	// 4. Apply the inverse of zoom and rotation
	if c.rotation != 0 {
		m.Rotate(-c.rotation) // Inverse rotation
	}
	m.Scale(1/c.zoom, 1/c.zoom) // Inverse zoom

	// 5. Move the origin to camera position (translating from screen center to camera pos)
	m.Translate(c.pos.X, c.pos.Y)

	// 6. Apply the transformation
	worldX, worldY := m.Apply(screenPos.X, screenPos.Y)

	return common.Vector2{X: worldX, Y: worldY}
}

// GetViewport returns the current viewport rectangle in world coordinates
func (c *Camera) GetViewport() common.Rectangle {
	halfWidth := c.config.ViewportSize.X / (2 * c.zoom)
	halfHeight := c.config.ViewportSize.Y / (2 * c.zoom)

	return common.Rectangle{
		Pos: common.Vector2{
			X: c.pos.X - halfWidth,
			Y: c.pos.Y - halfHeight,
		},
		Size: common.Vector2{
			X: halfWidth * 2,
			Y: halfHeight * 2,
		},
	}
}

// SetCenter explicitly sets the camera's position
func (c *Camera) SetCenter(pos common.Vector2) {
	c.pos = pos
}

// GetCenter returns the camera's center position
func (c *Camera) GetCenter() common.Vector2 {
	return c.pos
}

// SetZoom sets the camera zoom level
func (c *Camera) SetZoom(zoom float64) {
	c.zoom = math.Max(0.1, zoom) // Prevent negative or zero zoom
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

// clampToBounds ensures the camera stays within the configured bounds
func (c *Camera) clampToBounds() {
	bounds := c.config.Bounds
	if bounds == nil {
		return
	}

	// Calculate visible area at current zoom
	halfWidth := c.config.ViewportSize.X / (2 * c.zoom)
	halfHeight := c.config.ViewportSize.Y / (2 * c.zoom)

	// Calculate allowed camera position range
	minX := bounds.Pos.X + halfWidth
	maxX := bounds.Pos.X + bounds.Size.X - halfWidth
	minY := bounds.Pos.Y + halfHeight
	maxY := bounds.Pos.Y + bounds.Size.Y - halfHeight

	// Handle case where viewport is larger than bounds
	if halfWidth*2 > bounds.Size.X {
		// Center horizontally
		c.pos.X = bounds.Pos.X + bounds.Size.X/2
	} else {
		c.pos.X = math.Max(minX, math.Min(maxX, c.pos.X))
	}

	if halfHeight*2 > bounds.Size.Y {
		// Center vertically
		c.pos.Y = bounds.Pos.Y + bounds.Size.Y/2
	} else {
		c.pos.Y = math.Max(minY, math.Min(maxY, c.pos.Y))
	}
}
