package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

var (
	camera    rl.Camera2D
	gameState GameState
	dt        float32
)

type GameState struct {
	player  Player
	enemies []*Enemy
	rects   []rl.Rectangle
}

func initGameState() GameState {
	// Initialize your game state here
	player := NewPlayer()
	enemies := make([]*Enemy, 0)
	for i := 0; i < 100; i++ {
		enemies = append(enemies, NewEnemy())
	}
	gs := GameState{
		player:  player,
		enemies: enemies,
	}

	return gs
}

func gameLoop(gameStateCh chan<- GameState) {
	// Initialize your game state here
	gameState = initGameState()
	for {
		// Update your game state here
		var dx, dy float32
		if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
			dx = -1
		} else if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
			dx = 1
		}

		if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
			dy = -1
		} else if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
			dy = 1
		}

		// Restrict diagonal speed to be equal to horizontal or vertical speed
		if dx != 0 && dy != 0 {
			dx *= 0.7071 // approx. 1/sqrt(2)
			dy *= 0.7071 // approx. 1/sqrt(2)
		}

		gameState.player.pos.X += dx * gameState.player.speed.Value() * dt
		gameState.player.pos.Y += dy * gameState.player.speed.Value() * dt

		if rl.IsKeyDown(rl.KeyEqual) {
			// add 1% to the speed
			gameState.player.speed.AddBonus(1)
		}
		if rl.IsKeyDown(rl.KeyMinus) {
			if gameState.player.speed.bonus > 0 {
				// subtract 1% from the speed
				gameState.player.speed.AddBonus(-1)
			}
		}

		gameState.CheckCollisions()
		for _, e := range gameState.enemies {
			e.Update()
		}

		//gameState.player.pos.X = rl.Clamp(gameState.player.pos.X, 0, screenWidth-gameState.player.size.X)
		//gameState.player.pos.Y = rl.Clamp(gameState.player.pos.Y, 0, screenHeight-gameState.player.size.Y)

		// Send the updated game state through the channel
		gameStateCh <- gameState
	}
}

func (s *GameState) CheckCollisions() {
	// Check player collisions
	playerRect := s.player.Rect()
	for _, r := range s.rects {
		if rl.CheckCollisionRecs(playerRect, r) {
			fmt.Println("Collision!")
		}
	}

	// Check enemy collisions
	for _, e := range s.enemies {
		eRect := e.Rect()
		for _, r := range s.enemies {
			if rl.CheckCollisionRecs(eRect, r.Rect()) {
				if e != r {
					e.colour = rl.Red
					r.colour = rl.Red
					// get the direction of the collision and move the enemy away from the collision
					// by the amount of overlap
					overlap := rl.GetCollisionRec(eRect, r.Rect())
					if overlap.Width > overlap.Height {
						if e.pos.Y < r.Rect().Y {
							e.pos.Y -= overlap.Height
						} else {
							e.pos.Y += overlap.Height
						}
					} else {
						if e.pos.X < r.Rect().X {
							e.pos.X -= overlap.Width
						} else {
							e.pos.X += overlap.Width
						}
					}
				}
			}
		}
	}
}

func main() {
	// Initialize Raylib
	rl.InitWindow(screenWidth, screenHeight, "My Game")
	// Create a channel to communicate with the game logic goroutine
	gameStateCh := make(chan GameState)

	// Start the game logic goroutine
	go gameLoop(gameStateCh)

	camera = rl.NewCamera2D(rl.NewVector2(screenWidth/2, screenHeight/2), rl.NewVector2(0, 0), 0, 1)

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		dt = rl.GetFrameTime()
		// Poll the channel for updates from the game logic goroutine
		gameState = <-gameStateCh
		player := gameState.player
		// target the camera on the player
		target := rl.NewVector2(player.pos.X+(player.size.X/2), player.pos.Y+(player.size.Y/2))
		current := camera.Target
		lerpAmount := float32(0.2)
		x := current.X + (target.X-current.X)*lerpAmount
		y := current.Y + (target.Y-current.Y)*lerpAmount
		camera.Target = rl.NewVector2(x, y)

		// Draw the updated game state on the screen
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawLine(0, screenHeight/2, screenWidth, screenHeight/2, rl.Red)
		rl.DrawLine(screenWidth/2, 0, screenWidth/2, screenHeight, rl.Red)
		rl.BeginMode2D(camera)
		drawHelperGrid()

		//// Draw your game objects here using the updated game state
		rl.DrawRectangleV(gameState.player.pos, gameState.player.size, gameState.player.colour)

		for _, enemy := range gameState.enemies {
			rl.DrawRectangleV(enemy.pos, enemy.size, enemy.colour)
		}

		rl.EndMode2D()

		rl.DrawFPS(screenWidth-100, 10)
		rl.DrawText(fmt.Sprintf("bonusSpeed: %f, speed: %s, pos: %v", player.speed.bonus, player.speed.String(), player.pos), 10, 10, 20, rl.Black)

		rl.EndDrawing()
	}

	// Close Raylib
	rl.CloseWindow()
}

func drawHelperGrid() {
	// draw a grid to help with positioning
	for x := 0; x < screenWidth; x += 20 {
		rl.DrawLine(int32(x), 0, int32(x), screenHeight, rl.LightGray)
	}
	for y := 0; y < screenHeight; y += 20 {
		rl.DrawLine(0, int32(y), screenWidth, int32(y), rl.LightGray)
	}
}

func drawLineInDirection(pos rl.Vector2, dir rl.Vector2, length float32) {
	rl.DrawLineV(pos, rl.Vector2Add(pos, rl.Vector2Scale(dir, length)), rl.Red)
}
