package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

type Character struct {
	Pos   rl.Vector2
	Speed float32
}

func NewCharacter(x, y float32) *Character {
	return &Character{
		Pos:   rl.NewVector2(x, y),
		Speed: 6.0,
	}
}

func (c *Character) Move() {
	deltaX := 0.0
	deltaY := 0.0

	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		deltaX += 1.0
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		deltaX -= 1.0
	}
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		deltaY -= 1.0
	}
	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		deltaY += 1.0
	}

	if deltaX != 0 || deltaY != 0 {
		// Normalize the direction vector
		length := math.Sqrt(deltaX*deltaX + deltaY*deltaY)
		deltaX /= length
		deltaY /= length

		// Apply the normalized direction to the character's position
		c.Pos.X += float32(deltaX) * c.Speed
		c.Pos.Y += float32(deltaY) * c.Speed
	}
}
func (c *Character) Draw() {
	// Draw Character
	rl.DrawCircleV(c.Pos, 20, rl.Red)
}

// DrawGrid draws a grid extending beyond the current camera view
func DrawGrid(camera rl.Camera2D, spacing int) {
	// Calculate the camera's view bounds plus a buffer zone
	offsetX := int(camera.Offset.X) + 100 // Extend 100 pixels beyond the screen edge
	offsetY := int(camera.Offset.Y) + 100 // Extend 100 pixels beyond the screen edge

	// Calculate starting points adjusted to the nearest grid lines
	startX := int(camera.Target.X) - offsetX
	startY := int(camera.Target.Y) - offsetY
	endX := int(camera.Target.X) + offsetX
	endY := int(camera.Target.Y) + offsetY

	// Adjust to the nearest grid lines
	startX -= startX % spacing
	startY -= startY % spacing

	// Draw vertical lines
	for x := startX; x < endX; x += spacing {
		rl.DrawLine(int32(x), int32(startY), int32(x), int32(endY), rl.Fade(rl.LightGray, 0.7))
	}
	// Draw horizontal lines
	for y := startY; y < endY; y += spacing {
		rl.DrawLine(int32(startX), int32(y), int32(endX), int32(y), rl.Fade(rl.LightGray, 0.7))
	}
}

func updateCameraCenter(camera *rl.Camera2D, character *Character, deadZoneRadius float32) {
	// Calculate the distance from the character to the camera target
	distance := rl.Vector2Distance(character.Pos, camera.Target)

	// If the character is outside the dead zone radius, adjust the camera
	if distance > deadZoneRadius {
		// Calculate direction from camera to character
		direction := rl.Vector2Subtract(character.Pos, camera.Target)
		direction = rl.Vector2Normalize(direction)

		// Move camera target towards the character
		camera.Target = rl.Vector2Add(camera.Target, rl.Vector2Scale(direction, distance-deadZoneRadius))
	}
}

// Enemy represents an enemy character
type Enemy struct {
	Position          rl.Vector2
	MaxHealth, Health int
	Speed             float32
	Update            func(*Enemy, rl.Vector2) // Function to update enemy movement
	Active            bool
}

// NewEnemy initializes a new enemy with a specific behavior
func NewEnemy(posX, posY, speed float32, behavior func(*Enemy, rl.Vector2)) *Enemy {
	return &Enemy{
		MaxHealth: 10,
		Health:    10,
		Position:  rl.Vector2{X: posX, Y: posY},
		Speed:     speed,
		Update:    behavior,
		Active:    true,
	}
}

func (e *Enemy) Draw() {
	if e.Active {
		rl.DrawCircleV(e.Position, 20, rl.Blue)
	}
}

func (e *Enemy) TakeDamage(damage int) {
	e.Health -= damage
	if e.Health <= 0 {
		e.Health = 0
		e.Active = false
	}
}

func (e *Enemy) IsDead() bool {
	return e.Health <= 0
}

// Function type for generating enemy behavior functions
type BehaviorGenerator func() func(*Enemy, rl.Vector2)

// Generator for zigzag behavior
func ZigzagBehavior() func(*Enemy, rl.Vector2) {
	angle := 0.0 // Initialize angle state for the oscillation
	return func(enemy *Enemy, playerPosition rl.Vector2) {
		angle += 0.05 // Increment angle to change the sine value over time

		// Calculate the direct movement vector towards the player
		direction := rl.Vector2Subtract(playerPosition, enemy.Position)
		direction = rl.Vector2Normalize(direction)

		// Calculate a perpendicular vector to create the zigzag effect
		perpVector := rl.Vector2{X: -direction.Y, Y: direction.X}
		// Apply the sine function to oscillate the side-to-side movement
		perpVector = rl.Vector2Scale(perpVector, float32(math.Sin(angle)*2.0))

		// Combine the direct movement with the oscillating perpendicular movement
		finalDirection := rl.Vector2Add(direction, perpVector)
		finalDirection = rl.Vector2Normalize(finalDirection)

		// Move the enemy along the final direction vector scaled by speed
		enemy.Position = rl.Vector2Add(enemy.Position, rl.Vector2Scale(finalDirection, enemy.Speed))
	}
}

// DirectChase behavior: enemy moves directly towards the player
func DirectChase(enemy *Enemy, playerPosition rl.Vector2) {
	direction := rl.Vector2Subtract(playerPosition, enemy.Position)
	if rl.Vector2Length(direction) > 1 { // Avoid division by zero
		direction = rl.Vector2Normalize(direction)
		enemy.Position = rl.Vector2Add(enemy.Position, rl.Vector2Scale(direction, enemy.Speed))
	}
}

// MaintainDistance behavior: tries to keep a distance from the player
func MaintainDistance(enemy *Enemy, playerPosition rl.Vector2) {
	const idealDistance = 100.0 // Ideal distance from the player
	direction := rl.Vector2Subtract(playerPosition, enemy.Position)
	distance := rl.Vector2Length(direction)

	if distance > idealDistance {
		direction = rl.Vector2Normalize(direction)
		enemy.Position = rl.Vector2Add(enemy.Position, rl.Vector2Scale(direction, enemy.Speed))
	} else if distance < idealDistance {
		direction = rl.Vector2Normalize(direction)
		enemy.Position = rl.Vector2Subtract(enemy.Position, rl.Vector2Scale(direction, enemy.Speed))
	}
}

// Check and resolve collisions among all enemies
func resolveCollisions(enemies []*Enemy) {
	for i := 0; i < len(enemies); i++ {
		for j := i + 1; j < len(enemies); j++ {
			enemy1 := enemies[i]
			enemy2 := enemies[j]
			// Calculate distance between enemies
			distance := rl.Vector2Distance(enemy1.Position, enemy2.Position)
			// Assume both enemies have a radius of 20 for collision detection
			if distance < 40 { // 20 radius each, adjust as per your enemy size
				// Collision detected, push both enemies away from each other
				overlap := 40 - distance
				direction := rl.Vector2Subtract(enemy1.Position, enemy2.Position)
				if rl.Vector2Length(direction) > 0 {
					direction = rl.Vector2Normalize(direction)
				}
				// Adjust positions
				enemy1.Position = rl.Vector2Add(enemy1.Position, rl.Vector2Scale(direction, overlap/2))
				enemy2.Position = rl.Vector2Subtract(enemy2.Position, rl.Vector2Scale(direction, overlap/2))
			}
		}
	}
}

func checkCollisionWithPlayer(player *Character, enemies []*Enemy) {
	for _, enemy := range enemies {
		distance := rl.Vector2Distance(player.Pos, enemy.Position)
		if distance < 40 { // 20 radius each, adjust as per your enemy size
			enemy.TakeDamage(1)
			fmt.Println("Player collided with enemy, health:", enemy.Health)
		}
	}
}

var bg_dark = rl.NewColor(34, 32, 52, 255)

func main() {
	rl.InitWindow(800, 600, "Grid and Camera Control")
	defer rl.CloseWindow()

	var enemies []*Enemy
	character := NewCharacter(400, 300)
	for i := 0; i < 100; i++ {
		e1 := NewEnemy(rand.Float32()*800-0, rand.Float32()*450-0, 2.0, DirectChase)
		e2 := NewEnemy(rand.Float32()*800-0, rand.Float32()*450-0, 1.0, ZigzagBehavior())
		e3 := NewEnemy(rand.Float32()*800-0, rand.Float32()*450-0, 3.0, MaintainDistance)

		enemies = append(enemies, e1, e2, e3)
	}

	camera := rl.Camera2D{
		Target:   character.Pos,
		Offset:   rl.NewVector2(400, 300),
		Rotation: 0.0,
		Zoom:     1.0,
	}

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		character.Move()
		updateCameraCenter(&camera, character, 200)
		resolveCollisions(enemies)
		checkCollisionWithPlayer(character, enemies)
		rl.BeginDrawing()
		rl.ClearBackground(bg_dark)

		rl.BeginMode2D(camera)
		DrawGrid(camera, 50) // Draw grid lines every 50 pixels
		character.Draw()
		for _, enemy := range enemies {
			if enemy.Active {
				enemy.Update(enemy, character.Pos)
				enemy.Draw()
			}
		}
		rl.EndMode2D()

		rl.EndDrawing()
	}
}
