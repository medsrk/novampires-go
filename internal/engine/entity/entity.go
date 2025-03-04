package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/common"
)

// Entity represents a base game entity with core functionality
type Entity struct {
	// Core positioning and physics
	Position common.Vector2
	Velocity common.Vector2
	Rotation float64

	// Core identity
	ID uint64

	// Optional components
	sprite *SpriteComponent
	input  InputComponent
}

// NewEntity creates a new entity with the given parameters
func NewEntity(id uint64, position common.Vector2) *Entity {
	return &Entity{
		ID:       id,
		Position: position,
		Velocity: common.Vector2{},
		Rotation: 0,
	}
}

// Update updates the entity state
func (e *Entity) Update() {
	// Update position based on velocity
	e.Position = e.Position.Add(e.Velocity)

	// Update sprite if available
	if e.sprite != nil {
		e.sprite.Update(e)
	}

	// Process input if available
	if e.input != nil {
		e.input.ProcessInput(e)
	}
}

// Draw draws the entity
func (e *Entity) Draw(screen *ebiten.Image, renderer Renderer) {
	if e.sprite != nil {
		e.sprite.Draw(screen, renderer, e)
	}
}

// SetSprite assigns a sprite component to the entity
func (e *Entity) SetSprite(sprite *SpriteComponent) {
	e.sprite = sprite
}

// GetSprite returns the entity's sprite component
func (e *Entity) GetSprite() *SpriteComponent {
	return e.sprite
}

// SetInput assigns an input component to the entity
func (e *Entity) SetInput(input InputComponent) {
	e.input = input
}

// GetInput returns the entity's input component
func (e *Entity) GetInput() InputComponent {
	return e.input
}

// GetPosition returns the entity position
func (e *Entity) GetPosition() common.Vector2 {
	return e.Position
}

// GetPositionPtr returns a pointer to the entity's position
func (e *Entity) GetPositionPtr() *common.Vector2 {
	return &e.Position
}

// GetVelocity returns the entity velocity
func (e *Entity) GetVelocity() common.Vector2 {
	return e.Velocity
}

// SetVelocity sets the entity velocity
func (e *Entity) SetVelocity(vel common.Vector2) {
	e.Velocity = vel
}

// GetRotation returns the entity rotation
func (e *Entity) GetRotation() float64 {
	return e.Rotation
}

// SetRotation sets the entity rotation
func (e *Entity) SetRotation(rotation float64) {
	e.Rotation = rotation
}

// GetID returns the entity ID
func (e *Entity) GetID() uint64 {
	return e.ID
}
