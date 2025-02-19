package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"novampires-go/internal/common"
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
}

// NewRenderer creates a new renderer with specified configuration
func NewRenderer(config RenderConfig) *Renderer {
	return &Renderer{
		config: config,
	}
}

// Draw renders the current game state
func (r *Renderer) Draw(screen *ebiten.Image) {
	// Clear the screen with background color
	screen.Fill(r.config.ColorPalette.UIBackground)

	// Draw a simple grid background
	gridColor := color.RGBA{50, 50, 60, 80}
	gridSpacing := 100

	// Vertical grid lines
	for x := gridSpacing; x < screen.Bounds().Dx(); x += gridSpacing {
		vector.StrokeLine(
			screen,
			float32(x),
			0,
			float32(x),
			float32(screen.Bounds().Dy()),
			1,
			gridColor,
			false,
		)
	}

	// Horizontal grid lines
	for y := gridSpacing; y < screen.Bounds().Dy(); y += gridSpacing {
		vector.StrokeLine(
			screen,
			0,
			float32(y),
			float32(screen.Bounds().Dx()),
			float32(y),
			1,
			gridColor,
			false,
		)
	}

	// Draw a test player character in the center
	playerPos := common.Vector2{
		X: float64(screen.Bounds().Dx()) / 2,
		Y: float64(screen.Bounds().Dy()) / 2,
	}
	playerRotation := math.Pi / 4 // 45 degrees
	r.DrawPlayerCharacter(screen, playerPos, playerRotation, 20.0)

	// Draw aim line for player
	dirVec := common.Vector2{X: math.Cos(playerRotation), Y: math.Sin(playerRotation)}
	r.DrawAimLine(screen, playerPos, dirVec, 200)

	// Draw some test enemies
	enemyPositions := []common.Vector2{
		{X: float64(screen.Bounds().Dx()) / 4, Y: float64(screen.Bounds().Dy()) / 4},
		{X: float64(screen.Bounds().Dx()) * 3 / 4, Y: float64(screen.Bounds().Dy()) / 4},
		{X: float64(screen.Bounds().Dx()) / 4, Y: float64(screen.Bounds().Dy()) * 3 / 4},
		{X: float64(screen.Bounds().Dx()) * 3 / 4, Y: float64(screen.Bounds().Dy()) * 3 / 4},
	}

	for _, pos := range enemyPositions {
		angle := pos.Angle()
		r.DrawCircle(screen, pos, 15, r.config.ColorPalette.EnemyStandard)

		// Draw direction indicator
		dirX := pos.X + math.Cos(angle)*18
		dirY := pos.Y + math.Sin(angle)*18
		r.DrawLine(
			screen,
			pos,
			common.Vector2{X: dirX, Y: dirY},
			r.config.LineThickness,
			r.config.ColorPalette.UIForeground,
		)

		// Draw health bar
		healthBarWidth := 30.0
		healthBarHeight := 4.0
		healthBarPos := common.Vector2{
			X: pos.X - healthBarWidth/2,
			Y: pos.Y - 25,
		}
		r.DrawHealthBar(screen, healthBarPos, healthBarWidth, healthBarHeight, 0.7)
	}

	// Draw some projectiles
	projectilePositions := []common.Vector2{
		{X: playerPos.X + 100, Y: playerPos.Y},
		{X: playerPos.X - 50, Y: playerPos.Y + 120},
		{X: playerPos.X + 70, Y: playerPos.Y - 80},
	}

	for _, pos := range projectilePositions {
		r.DrawCircle(screen, pos, 5, r.config.ColorPalette.PlayerBullet)
	}
}

// DrawScene renders a complete game scene with background, entities and UI
func (r *Renderer) DrawScene(screen *ebiten.Image, entities []interface{}, ui interface{}) {
	// Clear the screen with background color
	screen.Fill(r.config.ColorPalette.UIBackground)

	// Draw background elements (if any)
	r.drawBackground(screen)

	// Draw all entities
	//for _, entity := range entities {
	//	// Type switch based on entity type
	//	// This would handle different entity types when implemented
	//}

	// Draw UI elements
	r.drawUI(screen, ui)
}

// drawBackground renders background elements
func (r *Renderer) drawBackground(screen *ebiten.Image) {
	// Draw grid or environment elements
	// For now, we'll leave it empty
}

// drawUI renders user interface elements
func (r *Renderer) drawUI(screen *ebiten.Image, ui interface{}) {
	// Render UI components when implemented
}

// DrawBullet draws a projectile
func (r *Renderer) DrawBullet(screen *ebiten.Image, position common.Vector2, direction common.Vector2,
	radius float64, isPlayerBullet bool) {
	// Choose color based on bullet owner
	var bulletColor color.RGBA
	if isPlayerBullet {
		bulletColor = r.config.ColorPalette.PlayerBullet
	} else {
		bulletColor = r.config.ColorPalette.EnemyBullet
	}

	// Draw the bullet
	r.DrawCircle(screen, position, radius, bulletColor)

	// Draw trail effect (optional)
	trailLength := radius * 2
	trailEnd := common.Vector2{
		X: position.X - direction.X*trailLength,
		Y: position.Y - direction.Y*trailLength,
	}

	// Use a semi-transparent version of bullet color for trail
	trailColor := bulletColor
	trailColor.A = 128

	r.DrawLine(screen, position, trailEnd, radius*0.8, trailColor)
}

// DrawEnemy renders an enemy with appropriate indicators
func (r *Renderer) DrawEnemy(screen *ebiten.Image, position common.Vector2,
	radius float64, rotation float64, health float64,
	enemyType int) {
	// Choose color based on enemy type
	var enemyColor color.RGBA
	switch enemyType {
	case 1: // Standard
		enemyColor = r.config.ColorPalette.EnemyStandard
	case 2: // Elite
		enemyColor = r.config.ColorPalette.EnemyElite
	case 3: // Boss
		enemyColor = r.config.ColorPalette.EnemyBoss
	default:
		enemyColor = r.config.ColorPalette.EnemyStandard
	}

	// Draw enemy body
	r.DrawCircle(screen, position, radius, enemyColor)

	// Draw direction indicator
	indicatorLength := radius * 1.2
	dirX := position.X + math.Cos(rotation)*indicatorLength
	dirY := position.Y + math.Sin(rotation)*indicatorLength

	r.DrawLine(
		screen,
		position,
		common.Vector2{X: dirX, Y: dirY},
		r.config.LineThickness,
		r.config.ColorPalette.UIForeground,
	)

	// Draw health bar
	healthBarWidth := radius * 2
	healthBarHeight := 4.0
	healthBarPos := common.Vector2{
		X: position.X - healthBarWidth/2,
		Y: position.Y - radius - 10,
	}
	r.DrawHealthBar(screen, healthBarPos, healthBarWidth, healthBarHeight, health)
}

// DrawCircle draws a filled circle
func (r *Renderer) DrawCircle(screen *ebiten.Image, position common.Vector2, radius float64, fill color.RGBA) {
	vector.DrawFilledCircle(
		screen,
		float32(position.X),
		float32(position.Y),
		float32(radius),
		fill,
		r.config.AntiAliasing,
	)
}

// DrawCircleOutline draws a circle outline
func (r *Renderer) DrawCircleOutline(screen *ebiten.Image, position common.Vector2, radius float64, lineWidth float64, stroke color.RGBA) {
	// No built-in circle outline in Ebiten, use strokedCircle
	numSegments := 24 // Adjust based on radius for better performance
	for i := 0; i < numSegments; i++ {
		angle1 := float64(i) / float64(numSegments) * 2 * 3.14159
		angle2 := float64(i+1) / float64(numSegments) * 2 * 3.14159

		x1 := position.X + math.Cos(angle1)*radius
		y1 := position.Y + math.Sin(angle1)*radius
		x2 := position.X + math.Cos(angle2)*radius
		y2 := position.Y + math.Sin(angle2)*radius

		vector.StrokeLine(
			screen,
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
	vector.DrawFilledRect(
		screen,
		float32(rect.Pos.X),
		float32(rect.Pos.Y),
		float32(rect.Size.X),
		float32(rect.Size.Y),
		fill,
		r.config.AntiAliasing,
	)
}

// DrawRectOutline draws a rectangle outline
func (r *Renderer) DrawRectOutline(screen *ebiten.Image, rect common.Rectangle, lineWidth float64, stroke color.RGBA) {
	x, y := rect.Pos.X, rect.Pos.Y
	w, h := rect.Size.X, rect.Size.Y

	vector.StrokeRect(
		screen,
		float32(x),
		float32(y),
		float32(w),
		float32(h),
		float32(lineWidth),
		stroke,
		r.config.AntiAliasing,
	)
}

// DrawLine draws a line
func (r *Renderer) DrawLine(screen *ebiten.Image, start, end common.Vector2, lineWidth float64, stroke color.RGBA) {
	vector.StrokeLine(
		screen,
		float32(start.X),
		float32(start.Y),
		float32(end.X),
		float32(end.Y),
		float32(lineWidth),
		stroke,
		r.config.AntiAliasing,
	)
}

// DrawHealthBar draws a health bar
func (r *Renderer) DrawHealthBar(screen *ebiten.Image, position common.Vector2, width, height float64, percent float64) {
	colors := r.config.ColorPalette

	// Background
	barRect := common.Rectangle{
		Pos:  common.Vector2{X: position.X, Y: position.Y},
		Size: common.Vector2{X: width, Y: height},
	}
	r.DrawRect(screen, barRect, colors.HealthBarBG)

	// Health fill
	fillWidth := width * percent
	if fillWidth > 0 {
		fillRect := common.Rectangle{
			Pos:  common.Vector2{X: position.X, Y: position.Y},
			Size: common.Vector2{X: fillWidth, Y: height},
		}
		r.DrawRect(screen, fillRect, colors.HealthBarFill)
	}
}

// DrawPlayerCharacter draws the player with rotation
func (r *Renderer) DrawPlayerCharacter(screen *ebiten.Image, position common.Vector2, rotation float64, radius float64) {
	colors := r.config.ColorPalette

	// Draw player body
	r.DrawCircle(screen, position, radius, colors.PlayerBody)

	// Draw direction indicator
	indicatorLength := radius * 1.2
	dirX := position.X + math.Cos(rotation)*indicatorLength
	dirY := position.Y + math.Sin(rotation)*indicatorLength

	r.DrawLine(
		screen,
		position,
		common.Vector2{X: dirX, Y: dirY},
		r.config.LineThickness,
		colors.PlayerOutline,
	)
}

// DrawAimLine draws the auto-aim targeting line
func (r *Renderer) DrawAimLine(screen *ebiten.Image, start common.Vector2, direction common.Vector2, length float64) {
	end := common.Vector2{
		X: start.X + direction.X*length,
		Y: start.Y + direction.Y*length,
	}

	r.DrawLine(
		screen,
		start,
		end,
		r.config.LineThickness*0.5,
		r.config.ColorPalette.PlayerAimLine,
	)
}

// DrawDebugInfo draws useful debug information when enabled
func (r *Renderer) DrawDebugInfo(screen *ebiten.Image, debugInfo map[string]string) {
	// Implement debug rendering if needed
}
