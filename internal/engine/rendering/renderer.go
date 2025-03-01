// internal/game/rendering/renderer.go
package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/camera"
)

// ColorPalette defines a consistent set of colors for rendering
type ColorPalette struct {
	// UI colors
	UIBackground color.RGBA
	UIForeground color.RGBA
	UIAccent     color.RGBA
	UIHighlight  color.RGBA

	// Player colors
	PlayerBody    color.RGBA
	PlayerOutline color.RGBA
	PlayerAimLine color.RGBA

	// Projectile colors
	PlayerBullet color.RGBA
	EnemyBullet  color.RGBA

	// Health bar colors
	HealthBarBG   color.RGBA
	HealthBarFill color.RGBA
	ShieldBarFill color.RGBA

	// Enemy colors
	EnemyStandard color.RGBA
	EnemyElite    color.RGBA
	EnemyBoss     color.RGBA

	// Effect colors
	HitFlash      color.RGBA
	ExplosionBase color.RGBA
	DamageNumber  color.RGBA
}

// DefaultColorPalette returns a default color scheme
func DefaultColorPalette() ColorPalette {
	return ColorPalette{
		// UI colors
		UIBackground: color.RGBA{20, 20, 30, 255},
		UIForeground: color.RGBA{220, 220, 220, 255},
		UIAccent:     color.RGBA{86, 156, 214, 255},
		UIHighlight:  color.RGBA{156, 220, 254, 255},

		// Player colors
		PlayerBody:    color.RGBA{50, 205, 50, 255}, // Green
		PlayerOutline: color.RGBA{220, 220, 220, 255},
		PlayerAimLine: color.RGBA{200, 200, 200, 180},

		// Projectile colors
		PlayerBullet: color.RGBA{68, 221, 255, 255}, // Cyan
		EnemyBullet:  color.RGBA{255, 68, 68, 255},  // Red

		// Health bar colors
		HealthBarBG:   color.RGBA{40, 40, 40, 200},
		HealthBarFill: color.RGBA{0, 200, 0, 255},
		ShieldBarFill: color.RGBA{0, 128, 255, 255},

		// Enemy colors
		EnemyStandard: color.RGBA{204, 0, 0, 255},   // Red
		EnemyElite:    color.RGBA{204, 0, 204, 255}, // Purple
		EnemyBoss:     color.RGBA{102, 0, 0, 255},   // Dark red

		// Effect colors
		HitFlash:      color.RGBA{255, 255, 255, 200},
		ExplosionBase: color.RGBA{255, 165, 0, 255},   // Orange
		DamageNumber:  color.RGBA{255, 255, 100, 255}, // Yellow
	}
}

// RenderConfig defines rendering parameters
type RenderConfig struct {
	Scale         float64
	ColorPalette  ColorPalette
	ParticleLimit int
	LineThickness float64
	EnableBloom   bool
	MotionBlur    bool
	AntiAliasing  bool
}

// DefaultRenderConfig returns sensible rendering defaults
func DefaultRenderConfig() RenderConfig {
	return RenderConfig{
		Scale:         1.0,
		ColorPalette:  DefaultColorPalette(),
		ParticleLimit: 1000,
		LineThickness: 2.0,
		EnableBloom:   true,
		MotionBlur:    false,
		AntiAliasing:  true,
	}
}

// Renderer provides methods for drawing game elements
type Renderer struct {
	config RenderConfig
	camera *camera.Camera

	// UI buffer for screen-space rendering
	uiBuffer *ebiten.Image
}

// NewRenderer creates a new renderer with specified configuration
func NewRenderer(config RenderConfig, camera *camera.Camera) *Renderer {
	return &Renderer{
		config: config,
		camera: camera,
	}
}

func (r *Renderer) BeginFrame(screen *ebiten.Image) {
	// Initialize or resize UI buffer based on screen dimensions
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	if r.uiBuffer == nil ||
		r.uiBuffer.Bounds().Dx() != screenWidth ||
		r.uiBuffer.Bounds().Dy() != screenHeight {
		r.uiBuffer = ebiten.NewImage(screenWidth, screenHeight)
	}

	// Clear the screen and UI buffer
	screen.Fill(r.config.ColorPalette.UIBackground)
	r.uiBuffer.Clear()
}

func (r *Renderer) EndFrame(screen *ebiten.Image) {
	// Draw UI on top of everything
	screen.DrawImage(r.uiBuffer, nil)
}

// DrawCircle draws a filled circle in world coordinates
func (r *Renderer) DrawCircle(screen *ebiten.Image, position common.Vector2, radius float64, fill color.RGBA) {
	// Get the viewport - if object is outside, skip drawing
	viewport := r.camera.GetViewport()
	circle := common.Rectangle{
		Pos:  common.Vector2{X: position.X - radius, Y: position.Y - radius},
		Size: common.Vector2{X: radius * 2, Y: radius * 2},
	}

	if !viewport.Intersects(circle) {
		return // Circle is outside the viewport, skip drawing
	}

	// Apply camera transform to get screen coordinates
	screenPos := r.worldToScreen(position)

	// Draw the circle at screen position
	vector.DrawFilledCircle(
		screen,
		float32(screenPos.X),
		float32(screenPos.Y),
		float32(radius*r.camera.GetZoom()), // Scale radius by zoom
		fill,
		r.config.AntiAliasing,
	)
}

// DrawCircleOutline draws a circle outline in world coordinates
func (r *Renderer) DrawCircleOutline(screen *ebiten.Image, position common.Vector2, radius float64, lineWidth float64, stroke color.RGBA) {
	// Check if circle is in viewport
	viewport := r.camera.GetViewport()
	circle := common.Rectangle{
		Pos:  common.Vector2{X: position.X - radius, Y: position.Y - radius},
		Size: common.Vector2{X: radius * 2, Y: radius * 2},
	}

	if !viewport.Intersects(circle) {
		return // Circle is outside the viewport, skip drawing
	}

	// Get screen position
	screenPos := r.worldToScreen(position)

	// Scale radius and line width by zoom
	screenRadius := radius * r.camera.GetZoom()
	screenLineWidth := lineWidth * r.camera.GetZoom()

	// Draw the circle outline
	numSegments := int(math.Max(12, screenRadius/4)) // Adjust segments based on size

	for i := 0; i < numSegments; i++ {
		angle1 := float64(i) / float64(numSegments) * 2 * math.Pi
		angle2 := float64(i+1) / float64(numSegments) * 2 * math.Pi

		x1 := screenPos.X + math.Cos(angle1)*screenRadius
		y1 := screenPos.Y + math.Sin(angle1)*screenRadius
		x2 := screenPos.X + math.Cos(angle2)*screenRadius
		y2 := screenPos.Y + math.Sin(angle2)*screenRadius

		vector.StrokeLine(
			screen,
			float32(x1),
			float32(y1),
			float32(x2),
			float32(y2),
			float32(screenLineWidth),
			stroke,
			r.config.AntiAliasing,
		)
	}
}

// DrawRect draws a filled rectangle in world coordinates
func (r *Renderer) DrawRect(screen *ebiten.Image, rect common.Rectangle, fill color.RGBA) {
	// Check if rect is in viewport
	viewport := r.camera.GetViewport()
	if !viewport.Intersects(rect) {
		return // Rectangle is outside the viewport, skip drawing
	}

	// Convert to screen coordinates
	screenRect := r.worldRectToScreen(rect)

	vector.DrawFilledRect(
		screen,
		float32(screenRect.Pos.X),
		float32(screenRect.Pos.Y),
		float32(screenRect.Size.X),
		float32(screenRect.Size.Y),
		fill,
		r.config.AntiAliasing,
	)
}

// DrawRectOutline draws a rectangle outline in world coordinates
func (r *Renderer) DrawRectOutline(screen *ebiten.Image, rect common.Rectangle, lineWidth float64, stroke color.RGBA) {
	// Check if rect is in viewport
	viewport := r.camera.GetViewport()
	if !viewport.Intersects(rect) {
		return
	}

	// Convert to screen coordinates
	screenRect := r.worldRectToScreen(rect)
	screenLineWidth := lineWidth * r.camera.GetZoom()

	vector.StrokeRect(
		screen,
		float32(screenRect.Pos.X),
		float32(screenRect.Pos.Y),
		float32(screenRect.Size.X),
		float32(screenRect.Size.Y),
		float32(screenLineWidth),
		stroke,
		r.config.AntiAliasing,
	)
}

// DrawLine draws a line in world coordinates
func (r *Renderer) DrawLine(screen *ebiten.Image, start, end common.Vector2, lineWidth float64, stroke color.RGBA) {
	// Check if line intersects the viewport
	viewport := r.camera.GetViewport()

	// Simple check using bounding box of the line
	minX := math.Min(start.X, end.X)
	minY := math.Min(start.Y, end.Y)
	maxX := math.Max(start.X, end.X)
	maxY := math.Max(start.Y, end.Y)

	lineBounds := common.Rectangle{
		Pos:  common.Vector2{X: minX, Y: minY},
		Size: common.Vector2{X: maxX - minX, Y: maxY - minY},
	}

	if !viewport.Intersects(lineBounds) {
		return
	}

	// Convert to screen coordinates
	screenStart := r.worldToScreen(start)
	screenEnd := r.worldToScreen(end)

	// Scale line width by zoom
	screenLineWidth := lineWidth * r.camera.GetZoom()

	vector.StrokeLine(
		screen,
		float32(screenStart.X),
		float32(screenStart.Y),
		float32(screenEnd.X),
		float32(screenEnd.Y),
		float32(screenLineWidth),
		stroke,
		r.config.AntiAliasing,
	)
}

// DrawHealthBar draws a health bar in world coordinates
func (r *Renderer) DrawHealthBar(screen *ebiten.Image, position common.Vector2, width, height float64, percent float64) {
	// Create rectangle for health bar
	healthBarRect := common.Rectangle{
		Pos:  position,
		Size: common.Vector2{X: width, Y: height},
	}

	// Check if health bar is in viewport
	if !r.camera.IsRectVisible(healthBarRect) {
		return
	}

	// Get screen coordinates
	screenPos := r.worldToScreen(position)
	screenWidth := width * r.camera.GetZoom()
	screenHeight := height * r.camera.GetZoom()

	// Background
	vector.DrawFilledRect(
		screen,
		float32(screenPos.X),
		float32(screenPos.Y),
		float32(screenWidth),
		float32(screenHeight),
		r.config.ColorPalette.HealthBarBG,
		r.config.AntiAliasing,
	)

	// Health fill
	fillWidth := screenWidth * percent
	if fillWidth > 0 {
		vector.DrawFilledRect(
			screen,
			float32(screenPos.X),
			float32(screenPos.Y),
			float32(fillWidth),
			float32(screenHeight),
			r.config.ColorPalette.HealthBarFill,
			r.config.AntiAliasing,
		)
	}
}

// DrawPlayerCharacter draws the player with rotation in world coordinates
func (r *Renderer) DrawPlayerCharacter(screen, playerImg *ebiten.Image, position common.Vector2, rotation float64, radius float64) {
	// Check if player is in viewport
	viewport := r.camera.GetViewport()
	playerRect := common.Rectangle{
		Pos:  common.Vector2{X: position.X - radius, Y: position.Y - radius},
		Size: common.Vector2{X: radius * 2, Y: radius * 2},
	}

	if !viewport.Intersects(playerRect) {
		return
	}

	// Get screen position
	screenPos := r.worldToScreen(position)
	screenRadius := radius * r.camera.GetZoom()

	// Draw player body
	if playerImg == nil {
		// Fallback to circle if no image provided
		vector.DrawFilledCircle(
			screen,
			float32(screenPos.X),
			float32(screenPos.Y),
			float32(screenRadius),
			r.config.ColorPalette.PlayerBody,
			r.config.AntiAliasing,
		)
	} else {
		// Create image drawing options
		op := &ebiten.DrawImageOptions{}

		// Get image dimensions
		imgWidth, imgHeight := float64(playerImg.Bounds().Dx()), float64(playerImg.Bounds().Dy())

		// Center the image (translate to origin for proper rotation)
		op.GeoM.Translate(-imgWidth/2, -imgHeight/2)

		//// Apply rotation if needed
		//if rotation != 0 {
		//	op.GeoM.Rotate(rotation)
		//}

		// Apply base scale (adjust as needed for your sprite)
		baseScale := (radius * 3) / imgWidth // Scale sprite to match radius
		op.GeoM.Scale(baseScale, baseScale)

		// Apply camera zoom
		op.GeoM.Scale(r.camera.GetZoom(), r.camera.GetZoom())

		// Translate to screen position
		op.GeoM.Translate(screenPos.X, screenPos.Y)

		// Draw the image
		screen.DrawImage(playerImg, op)
	}

	// Draw direction indicator
	indicatorLength := screenRadius * 1.2
	dirX := screenPos.X + math.Cos(rotation)*indicatorLength
	dirY := screenPos.Y + math.Sin(rotation)*indicatorLength

	vector.StrokeLine(
		screen,
		float32(screenPos.X),
		float32(screenPos.Y),
		float32(dirX),
		float32(dirY),
		float32(r.config.LineThickness*r.camera.GetZoom()),
		r.config.ColorPalette.PlayerOutline,
		r.config.AntiAliasing,
	)
}

// DrawPlayerSprite draws the player sprite with proper camera transforms
func (r *Renderer) DrawPlayerSprite(
	screen *ebiten.Image,
	sprite *ebiten.Image,
	position common.Vector2,
	rotation float64,
	scale float64,
	flipX bool,
) {
	// Skip if no sprite provided
	if sprite == nil {
		// Fallback to circle rendering
		r.DrawPlayerCharacter(screen, nil, position, rotation, scale)
		return
	}

	// Check if player is in viewport
	viewport := r.camera.GetViewport()
	playerRect := common.Rectangle{
		Pos:  common.Vector2{X: position.X - scale, Y: position.Y - scale},
		Size: common.Vector2{X: scale * 2, Y: scale * 2},
	}

	if !viewport.Intersects(playerRect) {
		return
	}

	// Get screen position
	screenPos := r.worldToScreen(position)

	// Create drawing options
	op := &ebiten.DrawImageOptions{}

	// Get image dimensions
	imgWidth, imgHeight := float64(sprite.Bounds().Dx()), float64(sprite.Bounds().Dy())

	// Center the image (translate to origin for proper rotation)
	op.GeoM.Translate(-imgWidth/2, -imgHeight/2)

	// Apply flip if needed (before rotation)
	scaleX := 1.0
	if flipX {
		scaleX = -1.0
	}
	op.GeoM.Scale(scaleX, 1.0)

	// Apply rotation
	//if rotation != 0 {
	//	op.GeoM.Rotate(rotation)
	//}

	// Apply scale directly - using the scale parameter as a direct multiplier
	// The scale parameter should represent how many times larger than original size
	op.GeoM.Scale(scale, scale)

	// Apply camera zoom
	op.GeoM.Scale(r.camera.GetZoom(), r.camera.GetZoom())

	// Translate to screen position
	op.GeoM.Translate(screenPos.X, screenPos.Y)

	// Draw the image
	screen.DrawImage(sprite, op)

	// Draw direction indicator (optional, you can remove if not needed)
	indicatorLength := 20.0 * r.camera.GetZoom() * scale
	dirX := screenPos.X + math.Cos(rotation)*indicatorLength
	dirY := screenPos.Y + math.Sin(rotation)*indicatorLength

	vector.StrokeLine(
		screen,
		float32(screenPos.X),
		float32(screenPos.Y),
		float32(dirX),
		float32(dirY),
		float32(r.config.LineThickness*r.camera.GetZoom()),
		r.config.ColorPalette.PlayerOutline,
		r.config.AntiAliasing,
	)
}

// DrawAimLine draws the auto-aim targeting line in world coordinates
func (r *Renderer) DrawAimLine(screen *ebiten.Image, start common.Vector2, direction common.Vector2, length float64) {
	// Create end position
	end := common.Vector2{
		X: start.X + direction.X*length,
		Y: start.Y + direction.Y*length,
	}

	// Convert to screen coordinates and draw
	screenStart := r.worldToScreen(start)
	screenEnd := r.worldToScreen(end)

	vector.StrokeLine(
		screen,
		float32(screenStart.X),
		float32(screenStart.Y),
		float32(screenEnd.X),
		float32(screenEnd.Y),
		float32(r.config.LineThickness*0.5*r.camera.GetZoom()),
		r.config.ColorPalette.PlayerAimLine,
		r.config.AntiAliasing,
	)
}

// DrawGrid draws a reference grid that follows camera transformations
func (r *Renderer) DrawGrid(screen *ebiten.Image) {
	viewport := r.camera.GetViewport()
	gridSpacing := 100.0
	gridColor := color.RGBA{50, 50, 60, 80}

	// Calculate grid bounds
	startX := math.Floor(viewport.Pos.X/gridSpacing) * gridSpacing
	startY := math.Floor(viewport.Pos.Y/gridSpacing) * gridSpacing
	endX := viewport.Pos.X + viewport.Size.X
	endY := viewport.Pos.Y + viewport.Size.Y

	zoom := r.camera.GetZoom()

	// Draw vertical grid lines
	for x := startX; x <= endX; x += gridSpacing {
		// Convert world to screen
		screenX1, screenY1 := r.worldToScreen(common.Vector2{X: x, Y: startY}).X, r.worldToScreen(common.Vector2{X: x, Y: startY}).Y
		screenX2, screenY2 := r.worldToScreen(common.Vector2{X: x, Y: endY}).X, r.worldToScreen(common.Vector2{X: x, Y: endY}).Y

		vector.StrokeLine(
			screen,
			float32(screenX1),
			float32(screenY1),
			float32(screenX2),
			float32(screenY2),
			float32(1*zoom), // Scale line width by zoom
			gridColor,
			true,
		)
	}

	// Draw horizontal grid lines
	for y := startY; y <= endY; y += gridSpacing {
		// Convert world to screen
		screenX1, screenY1 := r.worldToScreen(common.Vector2{X: startX, Y: y}).X, r.worldToScreen(common.Vector2{X: startX, Y: y}).Y
		screenX2, screenY2 := r.worldToScreen(common.Vector2{X: endX, Y: y}).X, r.worldToScreen(common.Vector2{X: endX, Y: y}).Y

		vector.StrokeLine(
			screen,
			float32(screenX1),
			float32(screenY1),
			float32(screenX2),
			float32(screenY2),
			float32(1*zoom), // Scale line width by zoom
			gridColor,
			true,
		)
	}
}

// DrawUIText draws text directly to the UI buffer (no transformation)
func (r *Renderer) DrawUIText(text string, pos common.Vector2, col color.RGBA) {
	ebitenutil.DebugPrintAt(r.uiBuffer, text, int(pos.X), int(pos.Y))
}

// Helper functions for coordinate transformations

// worldToScreen converts world coordinates to screen coordinates
func (r *Renderer) worldToScreen(worldPos common.Vector2) common.Vector2 {
	return r.camera.WorldToScreen(worldPos)
}

// worldRectToScreen converts a world rectangle to screen coordinates
func (r *Renderer) worldRectToScreen(worldRect common.Rectangle) common.Rectangle {
	// Get the top-left corner in screen coordinates
	screenPos := r.worldToScreen(worldRect.Pos)

	// Scale the size by the zoom factor
	zoom := r.camera.GetZoom()
	screenWidth := worldRect.Size.X * zoom
	screenHeight := worldRect.Size.Y * zoom

	return common.Rectangle{
		Pos:  screenPos,
		Size: common.Vector2{X: screenWidth, Y: screenHeight},
	}
}

// screenToWorld converts screen coordinates to world coordinates
func (r *Renderer) screenToWorld(screenPos common.Vector2) common.Vector2 {
	return r.camera.ScreenToWorld(screenPos)
}
