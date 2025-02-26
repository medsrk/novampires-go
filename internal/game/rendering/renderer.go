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
	config      RenderConfig
	camera      *camera.Camera
	worldBuffer *ebiten.Image
	uiBuffer    *ebiten.Image

	// World origin in buffer coordinates
	worldOriginX float64
	worldOriginY float64
}

// NewRenderer creates a new renderer with specified configuration
func NewRenderer(config RenderConfig, camera *camera.Camera) *Renderer {
	return &Renderer{
		config: config,
		camera: camera,
	}
}

func (r *Renderer) BeginFrame(screen *ebiten.Image) {
	// Get screen dimensions
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate the required buffer size based on camera zoom
	zoom := r.camera.GetZoom()
	// Add a small margin to prevent edge artifacts during rotation
	bufferWidth := int(float64(screenWidth) / zoom * 1.2)
	bufferHeight := int(float64(screenHeight) / zoom * 1.2)

	// Create or resize world buffer as needed
	if r.worldBuffer == nil ||
		r.worldBuffer.Bounds().Dx() != bufferWidth ||
		r.worldBuffer.Bounds().Dy() != bufferHeight {
		r.worldBuffer = ebiten.NewImage(bufferWidth, bufferHeight)
	}

	// Create or resize UI buffer as needed
	if r.uiBuffer == nil ||
		r.uiBuffer.Bounds().Dx() != screenWidth ||
		r.uiBuffer.Bounds().Dy() != screenHeight {
		r.uiBuffer = ebiten.NewImage(screenWidth, screenHeight)
	}

	// Clear buffers
	r.worldBuffer.Clear()
	r.worldBuffer.Fill(r.config.ColorPalette.UIBackground)
	r.uiBuffer.Clear()
}

func (r *Renderer) EndFrame(screen *ebiten.Image) {
	// Get camera information
	camZoom := r.camera.GetZoom()
	camRotation := r.camera.GetRotation()

	// Get buffer dimensions
	bufferWidth := float64(r.worldBuffer.Bounds().Dx())
	bufferHeight := float64(r.worldBuffer.Bounds().Dy())

	// Draw world buffer with camera transform
	op := &ebiten.DrawImageOptions{}

	// 1. Center the buffer around camera position
	op.GeoM.Translate(-bufferWidth/2, -bufferHeight/2)

	// 2. Apply rotation if needed
	if camRotation != 0 {
		op.GeoM.Rotate(camRotation)
	}

	// 3. Apply zoom
	op.GeoM.Scale(camZoom, camZoom)

	// 4. Move to screen center
	screenWidth := float64(screen.Bounds().Dx())
	screenHeight := float64(screen.Bounds().Dy())
	op.GeoM.Translate(screenWidth/2, screenHeight/2)

	// Draw the world buffer
	screen.DrawImage(r.worldBuffer, op)

	// Draw UI directly (no transform)
	screen.DrawImage(r.uiBuffer, nil)
}

// worldToBuffer converts world coordinates to buffer coordinates
func (r *Renderer) worldToBuffer(x, y float64) (float64, float64) {
	camPos := r.camera.GetCenter()
	bufferWidth := float64(r.worldBuffer.Bounds().Dx())
	bufferHeight := float64(r.worldBuffer.Bounds().Dy())

	// Offset coordinates relative to camera position
	return x - camPos.X + bufferWidth/2, y - camPos.Y + bufferHeight/2
}

// DrawCircle draws a filled circle
func (r *Renderer) DrawCircle(screen *ebiten.Image, position common.Vector2, radius float64, fill color.RGBA) {
	bufferX, bufferY := r.worldToBuffer(position.X, position.Y)

	vector.DrawFilledCircle(
		r.worldBuffer,
		float32(bufferX),
		float32(bufferY),
		float32(radius),
		fill,
		r.config.AntiAliasing,
	)
}

// DrawCircleOutline draws a circle outline
func (r *Renderer) DrawCircleOutline(screen *ebiten.Image, position common.Vector2, radius float64, lineWidth float64, stroke color.RGBA) {
	bufferX, bufferY := r.worldToBuffer(position.X, position.Y)

	// No built-in circle outline in Ebiten, use strokedCircle
	numSegments := 24 // Adjust based on radius for better performance
	for i := 0; i < numSegments; i++ {
		angle1 := float64(i) / float64(numSegments) * 2 * math.Pi
		angle2 := float64(i+1) / float64(numSegments) * 2 * math.Pi

		x1 := bufferX + math.Cos(angle1)*radius
		y1 := bufferY + math.Sin(angle1)*radius
		x2 := bufferX + math.Cos(angle2)*radius
		y2 := bufferY + math.Sin(angle2)*radius

		vector.StrokeLine(
			r.worldBuffer,
			float32(x1),
			float32(y1),
			float32(x2),
			float32(y2),
			float32(lineWidth),
			stroke,
			r.config.AntiAliasing,
		)
	}
}

// DrawRect draws a filled rectangle
func (r *Renderer) DrawRect(screen *ebiten.Image, rect common.Rectangle, fill color.RGBA) {
	bufferX, bufferY := r.worldToBuffer(rect.Pos.X, rect.Pos.Y)

	vector.DrawFilledRect(
		r.worldBuffer,
		float32(bufferX),
		float32(bufferY),
		float32(rect.Size.X),
		float32(rect.Size.Y),
		fill,
		r.config.AntiAliasing,
	)
}

// DrawRectOutline draws a rectangle outline
func (r *Renderer) DrawRectOutline(screen *ebiten.Image, rect common.Rectangle, lineWidth float64, stroke color.RGBA) {
	bufferX, bufferY := r.worldToBuffer(rect.Pos.X, rect.Pos.Y)

	vector.StrokeRect(
		r.worldBuffer,
		float32(bufferX),
		float32(bufferY),
		float32(rect.Size.X),
		float32(rect.Size.Y),
		float32(lineWidth),
		stroke,
		r.config.AntiAliasing,
	)
}

// DrawLine draws a line
func (r *Renderer) DrawLine(screen *ebiten.Image, start, end common.Vector2, lineWidth float64, stroke color.RGBA) {
	startX, startY := r.worldToBuffer(start.X, start.Y)
	endX, endY := r.worldToBuffer(end.X, end.Y)

	vector.StrokeLine(
		r.worldBuffer,
		float32(startX),
		float32(startY),
		float32(endX),
		float32(endY),
		float32(lineWidth),
		stroke,
		r.config.AntiAliasing,
	)
}

// DrawHealthBar draws a health bar
func (r *Renderer) DrawHealthBar(screen *ebiten.Image, position common.Vector2, width, height float64, percent float64) {
	bufferX, bufferY := r.worldToBuffer(position.X, position.Y)
	colors := r.config.ColorPalette

	// Background
	vector.DrawFilledRect(
		r.worldBuffer,
		float32(bufferX),
		float32(bufferY),
		float32(width),
		float32(height),
		colors.HealthBarBG,
		r.config.AntiAliasing,
	)

	// Health fill
	fillWidth := width * percent
	if fillWidth > 0 {
		vector.DrawFilledRect(
			r.worldBuffer,
			float32(bufferX),
			float32(bufferY),
			float32(fillWidth),
			float32(height),
			colors.HealthBarFill,
			r.config.AntiAliasing,
		)
	}
}

// DrawPlayerCharacter draws the player with rotation
func (r *Renderer) DrawPlayerCharacter(screen *ebiten.Image, position common.Vector2, rotation float64, radius float64) {
	bufferX, bufferY := r.worldToBuffer(position.X, position.Y)
	colors := r.config.ColorPalette

	// Draw player body
	vector.DrawFilledCircle(
		r.worldBuffer,
		float32(bufferX),
		float32(bufferY),
		float32(radius),
		colors.PlayerBody,
		r.config.AntiAliasing,
	)

	// Draw direction indicator
	indicatorLength := radius * 1.2
	dirX := bufferX + math.Cos(rotation)*indicatorLength
	dirY := bufferY + math.Sin(rotation)*indicatorLength

	vector.StrokeLine(
		r.worldBuffer,
		float32(bufferX),
		float32(bufferY),
		float32(dirX),
		float32(dirY),
		float32(r.config.LineThickness),
		colors.PlayerOutline,
		r.config.AntiAliasing,
	)
}

// DrawAimLine draws the auto-aim targeting line
func (r *Renderer) DrawAimLine(screen *ebiten.Image, start common.Vector2, direction common.Vector2, length float64) {
	startX, startY := r.worldToBuffer(start.X, start.Y)
	endX := startX + direction.X*length
	endY := startY + direction.Y*length

	vector.StrokeLine(
		r.worldBuffer,
		float32(startX),
		float32(startY),
		float32(endX),
		float32(endY),
		float32(r.config.LineThickness*0.5),
		r.config.ColorPalette.PlayerAimLine,
		r.config.AntiAliasing,
	)
}

// DrawGrid draws a reference grid that follows camera transformations
func (r *Renderer) DrawGrid(screen *ebiten.Image) {
	viewport := r.camera.GetViewport()
	gridSpacing := 100.0
	gridColor := color.RGBA{50, 50, 60, 80}

	// Calculate grid lines needed for current view
	startX := math.Floor(viewport.Pos.X/gridSpacing) * gridSpacing
	startY := math.Floor(viewport.Pos.Y/gridSpacing) * gridSpacing
	endX := viewport.Pos.X + viewport.Size.X + gridSpacing
	endY := viewport.Pos.Y + viewport.Size.Y + gridSpacing

	// Draw vertical lines
	for x := startX; x <= endX; x += gridSpacing {
		bufferX, _ := r.worldToBuffer(x, 0)
		bufferStartY, _ := r.worldToBuffer(0, startY)
		_, bufferEndY := r.worldToBuffer(0, endY)

		vector.StrokeLine(
			r.worldBuffer,
			float32(bufferX),
			float32(bufferStartY),
			float32(bufferX),
			float32(bufferEndY),
			1,
			gridColor,
			true,
		)
	}

	// Draw horizontal lines
	for y := startY; y <= endY; y += gridSpacing {
		_, bufferY := r.worldToBuffer(0, y)
		bufferStartX, _ := r.worldToBuffer(startX, 0)
		bufferEndX, _ := r.worldToBuffer(endX, 0)

		vector.StrokeLine(
			r.worldBuffer,
			float32(bufferStartX),
			float32(bufferY),
			float32(bufferEndX),
			float32(bufferY),
			1,
			gridColor,
			true,
		)
	}
}

// DrawUIText draws text directly to the UI buffer (no transformation)
func (r *Renderer) DrawUIText(text string, pos common.Vector2, col color.RGBA) {
	ebitenutil.DebugPrintAt(r.uiBuffer, text, int(pos.X), int(pos.Y))
}
