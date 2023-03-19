package main

import (
	b2 "github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ppm = 1
)

func main2() {
	gravity := b2.MakeB2Vec2(0.0, 10.0)
	world := b2.MakeB2World(gravity)

	var groundBodyDef b2.B2BodyDef
	groundBodyDef.Position.Set(100.0, float64(screenHeight-50)/ppm)

	groundBodyDef.Type = b2.B2BodyType.B2_staticBody
	groundBody := world.CreateBody(&groundBodyDef)

	var groundBox b2.B2PolygonShape
	groundBox.SetAsBox(float64(screenWidth)/ppm, 25.0/ppm)

	groundBody.CreateFixture(&groundBox, 0.0)

	// Create the dynamic body
	bodyDef := b2.MakeB2BodyDef()
	bodyDef.Type = b2.B2BodyType.B2_dynamicBody
	bodyDef.Position.Set(100.0, 10.0)
	body := world.CreateBody(&bodyDef)

	bodyBox := b2.MakeB2PolygonShape()
	bodyBox.SetAsBox(5.0, 5.0)

	fd := b2.MakeB2FixtureDef()
	fd.Shape = &bodyBox
	fd.Density = 1.0
	fd.Friction = 0.3

	body.CreateFixtureFromDef(&fd)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {

		world.Step(1.0/60.0, 8, 3)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw the ground
		rl.DrawRectangle(0, screenHeight-50, screenWidth, 50, rl.LightGray)

		// Draw the dynamic body
		position := body.GetPosition()
		angle := body.GetAngle()

		rl.DrawRectanglePro(
			rl.NewRectangle(float32(position.X*ppm)-25, float32(position.Y*ppm)-25, 50, 50),
			rl.NewVector2(25, 25),
			float32(angle*180.0/b2.B2_pi),
			rl.SkyBlue,
		)

		rl.EndDrawing()

	}
}

func DrawB2PolygonShape(body *b2.B2Body, shape b2.B2PolygonShape) {
	pos := b2rlv(body.GetPosition())

	for i := 0; i < shape.M_count; i++ {
		if i == shape.M_count-1 {
			rl.DrawLineV(rl.Vector2Add(b2rlv(shape.M_vertices[i]), pos), rl.Vector2Add(b2rlv(shape.M_vertices[0]), pos), rl.Black)
		} else {
			rl.DrawLineV(rl.Vector2Add(b2rlv(shape.M_vertices[i]), pos), rl.Vector2Add(b2rlv(shape.M_vertices[i+1]), pos), rl.Black)
		}
	}
}

func b2rlv(vec b2.B2Vec2) rl.Vector2 {
	return rl.NewVector2(float32(vec.X), float32(vec.Y))
}
