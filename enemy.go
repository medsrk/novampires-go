package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

const (
	top = iota
	right
	bottom
	left
)

type Enemy struct {
	pos         rl.Vector2
	size        rl.Vector2
	bb          rl.BoundingBox
	speed       Speed
	dir         rl.Vector2
	isColliding bool
	colour      rl.Color
}

func NewEnemy() *Enemy {
	e := &Enemy{
		pos:    randomPositionOffScreen(),
		size:   rl.NewVector2(20, 20),
		bb:     rl.NewBoundingBox(rl.NewVector3(0, 0, 0), rl.NewVector3(20, 20, 0)),
		speed:  NewSpeed(100, 0),
		colour: rl.Green,
	}
	gameState.rects = append(gameState.rects, e.Rect())

	return e
}

func randomPositionOffScreen() rl.Vector2 {
	var pos rl.Vector2
	// pick a random side of the screen
	side := rand.Intn(3)
	switch side {
	case top:
		pos.X = rand.Float32() * screenWidth
		pos.Y = 0
	case right:
		pos.X = screenWidth
		pos.Y = rand.Float32() * screenHeight
	case bottom:
		pos.X = rand.Float32() * screenWidth
		pos.Y = screenHeight
	case left:
		pos.X = 0
		pos.Y = rand.Float32() * screenHeight
	}

	return pos
}

func (e *Enemy) Rect() rl.Rectangle {
	return rl.NewRectangle(e.pos.X, e.pos.Y, e.size.X, e.size.Y)
}

func (e *Enemy) Move(target rl.Vector2) {
	dir := rl.Vector2Subtract(target, e.pos)
	length := rl.Vector2Length(dir)
	if length > 0 {
		// Scale the direction vector to the desired length
		scale := e.speed.Value() * dt / length
		dir = rl.Vector2Scale(dir, scale)
	}

	e.dir = dir
	e.pos = rl.Vector2Add(e.pos, dir)
}

func (e *Enemy) Update() {
	// Move towards the player
	e.Move(gameState.player.pos)
}

func (e *Enemy) Center() rl.Vector2 {
	return rl.NewVector2(e.pos.X+(e.size.X/2), e.pos.Y+(e.size.Y/2))
}
