package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	pos    rl.Vector2
	size   rl.Vector2
	speed  Speed
	colour rl.Color
}

func NewPlayer() Player {
	p := Player{
		pos:    rl.NewVector2(screenWidth/2, screenHeight/2),
		size:   rl.NewVector2(20, 20),
		speed:  NewSpeed(200, 0),
		colour: rl.Red,
	}
	gameState.rects = append(gameState.rects, p.Rect())

	return p
}

func (p *Player) Rect() rl.Rectangle {
	return rl.NewRectangle(p.pos.X, p.pos.Y, p.size.X, p.size.Y)
}

func (p *Player) Update() {

}

func (p *Player) Move(dx, dy float32) {
	p.pos.X += dx * p.speed.Value()
	p.pos.Y += dy * p.speed.Value()
}

func (p *Player) BoundingBox() rl.BoundingBox {
	return rl.NewBoundingBox(rl.NewVector3(p.pos.X, p.pos.Y, 0), rl.NewVector3(p.size.X, p.size.Y, 0))
}
