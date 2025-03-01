// internal/engine/camera/camera.go
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

	// Visible world area
	visibleArea common.Rectangle
}

func New() *Camera {
	return NewWithConfig(DefaultConfig())
}

func NewWithConfig(config *Config) *Camera {
	cam := &Camera{
		config:   config,
		zoom:     1.0,
		rotation: 0.0,
	}

	// Initialize the visible area
	cam.updateVisibleArea()
	return cam
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

	// Update the visible area after position changes
	c.updateVisibleArea()
}

// updateVisibleArea calculates the world rectangle that's currently visible
func (c *Camera) updateVisibleArea() {
	// Calculate the half-sizes of the viewport in world coordinates
	halfWidth := c.config.ViewportSize.X / (2 * c.zoom)
	halfHeight := c.config.ViewportSize.Y / (2 * c.zoom)

	// Set the visible area based on the current camera position
	c.visibleArea = common.Rectangle{
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

// SetTarget sets the target position for the camera to follow
func (c *Camera) SetTarget(target *common.Vector2) {
	c.target = target
}

// GetTarget returns the current camera target position
func (c *Camera) GetTarget() *common.Vector2 {
	return c.target
}

// GetTransform returns the transformation matrix for rendering
func (c *Camera) GetTransform() ebiten.GeoM {
	m := ebiten.GeoM{}

	// 1. Translate to center the camera position
	m.Translate(-c.pos.X, -c.pos.Y)

	// 2. Scale according to zoom
	m.Scale(c.zoom, c.zoom)

	// 3. Rotate if needed
	if c.rotation != 0 {
		// For rotation, we need to:
		// - Move to origin
		// - Rotate around origin
		// - Move back

		//m.Translate(screenWidth/2/c.zoom, screenHeight/2/c.zoom)
		m.Rotate(c.rotation)
		//m.Translate(-screenWidth/2/c.zoom, -screenHeight/2/c.zoom)
	}

	// 4. Center on screen
	screenWidth := c.config.ViewportSize.X
	screenHeight := c.config.ViewportSize.Y
	m.Translate(screenWidth/2, screenHeight/2)

	return m
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenPos common.Vector2) common.Vector2 {
	// Create inverse transform to convert screen to world
	screenWidth := c.config.ViewportSize.X
	screenHeight := c.config.ViewportSize.Y

	// Start with identity matrix
	m := ebiten.GeoM{}

	// 1. Translate screen position to be relative to screen center
	m.Translate(-screenWidth/2, -screenHeight/2)

	// 2. Apply inverse zoom
	m.Scale(1/c.zoom, 1/c.zoom)

	// 3. Apply inverse rotation if needed
	if c.rotation != 0 {
		m.Rotate(-c.rotation)
	}

	// 4. Translate to world position
	m.Translate(c.pos.X, c.pos.Y)

	// Apply transformation
	worldX, worldY := m.Apply(screenPos.X, screenPos.Y)
	return common.Vector2{X: worldX, Y: worldY}
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldPos common.Vector2) common.Vector2 {
	// Apply the same transform used for rendering
	m := c.GetTransform()
	screenX, screenY := m.Apply(worldPos.X, worldPos.Y)
	return common.Vector2{X: screenX, Y: screenY}
}

// GetViewport returns the current viewport rectangle in world coordinates
func (c *Camera) GetViewport() common.Rectangle {
	return c.visibleArea
}

// GetVisibleTiles returns the tile coordinates that are visible
// minTileX, minTileY, maxTileX, maxTileY are in tile coordinates
func (c *Camera) GetVisibleTiles(tileSize int) (minTileX, minTileY, maxTileX, maxTileY int) {
	// Convert viewport bounds to tile coordinates
	minTileX = int(math.Floor(c.visibleArea.Pos.X / float64(tileSize)))
	minTileY = int(math.Floor(c.visibleArea.Pos.Y / float64(tileSize)))

	maxTileX = int(math.Ceil((c.visibleArea.Pos.X + c.visibleArea.Size.X) / float64(tileSize)))
	maxTileY = int(math.Ceil((c.visibleArea.Pos.Y + c.visibleArea.Size.Y) / float64(tileSize)))

	return
}

// IsRectVisible checks if a rectangle in world coordinates is visible on screen
func (c *Camera) IsRectVisible(rect common.Rectangle) bool {
	return c.visibleArea.Intersects(rect)
}

// SetCenter explicitly sets the camera's position
func (c *Camera) SetCenter(pos common.Vector2) {
	c.pos = pos

	// Clamp to bounds if necessary
	if c.config.Bounds != nil {
		c.clampToBounds()
	}

	// Update the visible area after position changes
	c.updateVisibleArea()
}

// GetCenter returns the camera's center position
func (c *Camera) GetCenter() common.Vector2 {
	return c.pos
}

// SetZoom sets the camera zoom level
func (c *Camera) SetZoom(zoom float64) {
	c.zoom = math.Max(0.1, zoom) // Prevent negative or zero zoom
	c.updateVisibleArea()
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
