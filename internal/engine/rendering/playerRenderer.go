package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/common"
)

func (r *Renderer) DrawLayeredPlayerSprite(
	screen *ebiten.Image,
	baseSprite *ebiten.Image,
	eyeSprite *ebiten.Image,
	position common.Vector2,
	eyePosition common.Vector2, // Relative position from character center
	rotation float64,
	scale float64,
	flipX bool,
) {
	// Skip if no base sprite provided
	if baseSprite == nil {
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

	// STEP 1: Draw the base character sprite
	baseOp := &ebiten.DrawImageOptions{}

	// Get image dimensions
	imgWidth, imgHeight := float64(baseSprite.Bounds().Dx()), float64(baseSprite.Bounds().Dy())

	// Center the image
	baseOp.GeoM.Translate(-imgWidth/2, -imgHeight/2)

	// Apply flip if needed
	scaleX := 1.0
	if flipX {
		scaleX = -1.0
	}
	baseOp.GeoM.Scale(scaleX, 1.0)

	// Apply scale
	baseOp.GeoM.Scale(scale, scale)

	// Apply camera zoom
	baseOp.GeoM.Scale(r.camera.GetZoom(), r.camera.GetZoom())

	// Translate to screen position
	baseOp.GeoM.Translate(screenPos.X, screenPos.Y)

	// Draw the base character sprite
	screen.DrawImage(baseSprite, baseOp)

	// STEP 2: Draw the eye sprite as a separate layer
	if eyeSprite != nil {
		eyeOp := &ebiten.DrawImageOptions{}

		// Get eye image dimensions
		eyeWidth, eyeHeight := float64(eyeSprite.Bounds().Dx()), float64(eyeSprite.Bounds().Dy())

		// Center the eye image
		eyeOp.GeoM.Translate(-eyeWidth/2, -eyeHeight/2)

		// Apply flip if needed
		if flipX {
			eyeOp.GeoM.Scale(-1.0, 1.0)
		}

		// Apply scale
		eyeOp.GeoM.Scale(scale, scale)

		// Apply camera zoom
		eyeOp.GeoM.Scale(r.camera.GetZoom(), r.camera.GetZoom())

		// Calculate eye position - adjusted for relative positioning
		eyeOffsetX, eyeOffsetY := eyePosition.X, eyePosition.Y

		// If flipped, invert the X offset
		if flipX {
			eyeOffsetX = -eyeOffsetX
		}

		// Scale the offset by the character scale
		eyeOffsetX *= scale
		eyeOffsetY *= scale

		// Scale the offset by camera zoom
		eyeOffsetX *= r.camera.GetZoom()
		eyeOffsetY *= r.camera.GetZoom()

		// Translate to screen position with offset
		eyeOp.GeoM.Translate(screenPos.X+eyeOffsetX, screenPos.Y+eyeOffsetY)

		// Draw the eye sprite
		screen.DrawImage(eyeSprite, eyeOp)
	}
}
