// internal/entity/renderer_adapter.go
package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/rendering"
)

// RendererAdapter adapts the rendering.Renderer to the entity.Renderer interface
type RendererAdapter struct {
	renderer *rendering.Renderer
}

// NewRendererAdapter creates a new renderer adapter
func NewRendererAdapter(renderer *rendering.Renderer) *RendererAdapter {
	return &RendererAdapter{
		renderer: renderer,
	}
}

// DrawSprite draws a sprite with the wrapped renderer
func (r *RendererAdapter) DrawSprite(
	screen *ebiten.Image,
	sprite *ebiten.Image,
	position common.Vector2,
	rotation float64,
	scale float64,
	flipX bool,
) {
	r.renderer.DrawPlayerSprite(screen, sprite, position, rotation, scale, flipX)
}

// DrawLayeredSprite draws a sprite with an overlay (like eyes) with the wrapped renderer
func (r *RendererAdapter) DrawLayeredSprite(
	screen *ebiten.Image,
	baseSprite, overlaySprite *ebiten.Image,
	position, overlayOffset common.Vector2,
	rotation, scale float64,
	flipX bool,
) {
	r.renderer.DrawLayeredPlayerSprite(
		screen,
		baseSprite,
		overlaySprite,
		position,
		overlayOffset,
		rotation,
		scale,
		flipX,
	)
}

// DrawAimLine draws an aim line with the wrapped renderer
func (r *RendererAdapter) DrawAimLine(
	screen *ebiten.Image,
	start common.Vector2,
	direction common.Vector2,
	length float64,
) {
	r.renderer.DrawAimLine(screen, start, direction, length)
}

// DrawCircle draws a filled circle in world coordinates
func (r *RendererAdapter) DrawCircle(
	screen *ebiten.Image,
	position common.Vector2,
	radius float64,
	fill color.RGBA,
) {
	r.renderer.DrawCircle(screen, position, radius, fill)
}

// DrawLine draws a line in world coordinates
func (r *RendererAdapter) DrawLine(
	screen *ebiten.Image,
	start, end common.Vector2,
	lineWidth float64,
	stroke color.RGBA,
) {
	r.renderer.DrawLine(screen, start, end, lineWidth, stroke)
}

// DrawHealthBar draws a health bar in world coordinates
func (r *RendererAdapter) DrawHealthBar(
	screen *ebiten.Image,
	position common.Vector2,
	width, height float64,
	percent float64,
) {
	r.renderer.DrawHealthBar(screen, position, width, height, percent)
}

// DrawGrid draws a reference grid
func (r *RendererAdapter) DrawGrid(screen *ebiten.Image) {
	r.renderer.DrawGrid(screen)
}
