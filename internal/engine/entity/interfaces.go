package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"novampires-go/internal/common"
	"time"
)

// Renderer defines an interface for rendering operations
type Renderer interface {
	DrawSprite(screen *ebiten.Image, sprite *ebiten.Image, position common.Vector2, rotation float64, scale float64, flipX bool)
	DrawLayeredSprite(screen *ebiten.Image, baseSprite, overlaySprite *ebiten.Image, position, overlayOffset common.Vector2, rotation, scale float64, flipX bool)
	DrawAimLine(screen *ebiten.Image, start common.Vector2, direction common.Vector2, length float64)
	DrawCircle(screen *ebiten.Image, position common.Vector2, radius float64, fill color.RGBA)
	DrawLine(screen *ebiten.Image, start, end common.Vector2, lineWidth float64, stroke color.RGBA)
	DrawHealthBar(screen *ebiten.Image, position common.Vector2, width, height float64, percent float64)
	DrawGrid(screen *ebiten.Image)
}

// InputComponent defines an interface for processing input
type InputComponent interface {
	ProcessInput(entity *Entity)
	GetAimDirection() common.Vector2
}

// AnimationController defines an interface for controlling animations
type AnimationController interface {
	Update(dt time.Duration)
	GetSprite() *ebiten.Image
	PlayAnimation(name string)
	GetCurrentAnimation() string
	SetFlipX(flip bool)
	GetFlipX() bool
	IsReversed() bool
	ReverseAnimation()
}

// TargetProvider defines an interface for entities that can be targeted
type TargetProvider interface {
	GetTargetInfo() common.TargetInfo
}

// Updatable defines an interface for objects that can be updated
type Updatable interface {
	Update()
}

// Drawable defines an interface for objects that can be drawn
type Drawable interface {
	Draw(screen *ebiten.Image, renderer Renderer)
}

// Entity defines a common interface for all game entities
type EntityInterface interface {
	Updatable
	Drawable
	GetPosition() common.Vector2
	GetPositionPtr() *common.Vector2
	GetVelocity() common.Vector2
	SetVelocity(vel common.Vector2)
	GetRotation() float64
	SetRotation(rotation float64)
	GetID() uint64
}
