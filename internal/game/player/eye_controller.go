package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"math/rand"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/sprite"
	"time"
)

var blinkIntervalMin = 180
var blinkIntervalRange = 120

// Direction represents where the eyes are looking
type Direction int

const (
	LookingCenter Direction = iota
	LookingRight
	LookingUp
	LookingDown
	LookingUpRight
	LookingDownRight
)

type EyeController struct {
	// The eye spritesheet
	spriteSheet *ebiten.Image

	// Current eye sprite
	currentSprite *ebiten.Image

	// All eye direction frames in a single row
	allFrames []sprite.FrameData

	// The blink animation
	animation *sprite.Animation

	// Current direction
	direction Direction

	// Should flip horizontally (for left directions)
	flipX bool

	// Is currently blinking
	isBlinking bool

	// Blink timer
	blinkTimer    int
	blinkInterval int

	// Position relative to character center
	position common.Vector2
}

func NewEyeController() *EyeController {
	return &EyeController{
		isBlinking:    false,
		blinkTimer:    0,
		blinkInterval: 180 + rand.Intn(120), // 3-5 seconds
		position:      common.Vector2{X: 0, Y: 0},
		direction:     LookingCenter,
		flipX:         false,
	}
}

func (c *EyeController) SetSpriteSheet(sheet *ebiten.Image) {
	c.spriteSheet = sheet
	if c.spriteSheet == nil {
		return
	}

	frameWidth := 96
	frameDuration := 100

	// Load all frames from a single row
	// Frame order: Center (3), Right (3), Left (3), Up (3), Down (3), UpRight (3), DownRight (3), UpLeft (3), DownLeft (3)
	c.allFrames = []sprite.FrameData{
		// Center (0, 1, 2)
		{0, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 2, 0, frameWidth, frameWidth, frameDuration},

		// Right (3, 4, 5)
		{frameWidth * 3, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 4, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 5, 0, frameWidth, frameWidth, frameDuration},

		// Left (6, 7, 8)
		{frameWidth * 6, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 7, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 8, 0, frameWidth, frameWidth, frameDuration},

		// Up (9, 10, 11)
		{frameWidth * 9, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 10, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 11, 0, frameWidth, frameWidth, frameDuration},

		// Down (12, 13, 14)
		{frameWidth * 12, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 13, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 14, 0, frameWidth, frameWidth, frameDuration},

		// UpRight (15, 16, 17)
		{frameWidth * 15, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 16, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 17, 0, frameWidth, frameWidth, frameDuration},

		// DownRight (18, 19, 20)
		{frameWidth * 18, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 19, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 20, 0, frameWidth, frameWidth, frameDuration},

		// UpLeft (21, 22, 23)
		{frameWidth * 21, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 22, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 23, 0, frameWidth, frameWidth, frameDuration},

		// DownLeft (24, 25, 26)
		{frameWidth * 24, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 25, 0, frameWidth, frameWidth, frameDuration},
		{frameWidth * 26, 0, frameWidth, frameWidth, frameDuration},
	}

	// Create initial blink animation for center eyes
	blinkSequence := []int{0, 1, 2, 1}
	c.animation = sprite.NewAnimationWithSequence(c.allFrames[:3], blinkSequence, false)

	// Initialize with the first frame
	frame := c.allFrames[0]
	rect := image.Rect(
		frame.SrcX,
		frame.SrcY,
		frame.SrcX+frame.SrcWidth,
		frame.SrcY+frame.SrcHeight,
	)
	c.currentSprite = sheet.SubImage(rect).(*ebiten.Image)
}

func (c *EyeController) Update(dt time.Duration) {
	if c.spriteSheet == nil {
		return
	}

	c.updateBlinking(dt)
	c.updateSprite()
}

func (c *EyeController) updateBlinking(dt time.Duration) {
	if c.isBlinking {
		c.animation.Update(dt)

		if c.animation.IsFinished() {
			c.isBlinking = false
			c.blinkTimer = 0
			c.blinkInterval = blinkIntervalMin + rand.Intn(blinkIntervalRange)
		}
		return
	}

	c.blinkTimer++

	if c.blinkTimer >= c.blinkInterval {
		c.isBlinking = true

		// Create a new animation for the current eye direction
		var startIdx int
		switch c.direction {
		case LookingCenter:
			startIdx = 0
		case LookingRight:
			if c.flipX {
				// Use Left frames (6-8)
				startIdx = 6
			} else {
				// Use Right frames (3-5)
				startIdx = 3
			}
		case LookingUp:
			startIdx = 9
		case LookingDown:
			startIdx = 12
		case LookingUpRight:
			if c.flipX {
				// Use UpLeft frames (21-23)
				startIdx = 21
			} else {
				// Use UpRight frames (15-17)
				startIdx = 15
			}
		case LookingDownRight:
			if c.flipX {
				// Use DownLeft frames (24-26)
				startIdx = 24
			} else {
				// Use DownRight frames (18-20)
				startIdx = 18
			}
		}

		blinkSequence := []int{0, 1, 2, 1}
		directionFrames := c.allFrames[startIdx : startIdx+3]
		c.animation = sprite.NewAnimationWithSequence(directionFrames, blinkSequence, false)
		c.animation.Reset()
	}
}

// updateSprite updates the current sprite based on state
func (c *EyeController) updateSprite() {
	if c.isBlinking {
		// Get current frame from the blinking animation
		frame := c.animation.GetCurrentFrame()

		// Update the current sprite from the frame data
		rect := image.Rect(
			frame.SrcX,
			frame.SrcY,
			frame.SrcX+frame.SrcWidth,
			frame.SrcY+frame.SrcHeight,
		)
		c.currentSprite = c.spriteSheet.SubImage(rect).(*ebiten.Image)
	} else {
		// When not blinking, use the first frame for current direction
		var frameIdx int
		switch c.direction {
		case LookingCenter:
			frameIdx = 0
		case LookingRight:
			if c.flipX {
				// Use Left frames
				frameIdx = 6
			} else {
				// Use Right frames
				frameIdx = 3
			}
		case LookingUp:
			frameIdx = 9
		case LookingDown:
			frameIdx = 12
		case LookingUpRight:
			if c.flipX {
				// Use UpLeft frames
				frameIdx = 21
			} else {
				// Use UpRight frames
				frameIdx = 15
			}
		case LookingDownRight:
			if c.flipX {
				// Use DownLeft frames
				frameIdx = 24
			} else {
				// Use DownRight frames
				frameIdx = 18
			}
		}

		frame := c.allFrames[frameIdx]
		rect := image.Rect(
			frame.SrcX,
			frame.SrcY,
			frame.SrcX+frame.SrcWidth,
			frame.SrcY+frame.SrcHeight,
		)
		c.currentSprite = c.spriteSheet.SubImage(rect).(*ebiten.Image)
	}
}

// UpdateLookDirection updates the eye direction based on the aim vector
func (c *EyeController) UpdateLookDirection(aimDirection common.Vector2) {
	if c.spriteSheet == nil {
		return
	}

	oldDirection := c.direction
	oldFlipX := c.flipX

	// Skip if vector is too small
	magnitude := aimDirection.Magnitude()
	if magnitude < 0.3 {
		c.direction = LookingCenter
		c.flipX = false

		// Only update sprite if direction changed
		if oldDirection != c.direction || oldFlipX != c.flipX {
			c.updateSprite()
		}
		return
	}

	// Normalize angle to 0-2π
	angle := math.Atan2(aimDirection.Y, aimDirection.X)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Default to not flipping
	c.flipX = false

	// Determine direction based on angle
	if angle >= 7*math.Pi/4 || angle < math.Pi/4 {
		// Right
		c.direction = LookingRight
	} else if angle >= math.Pi/4 && angle < 3*math.Pi/4 {
		// Down
		c.direction = LookingDown
	} else if angle >= 3*math.Pi/4 && angle < 5*math.Pi/4 {
		// Left (flip right)
		c.direction = LookingRight
		//c.flipX = true
	} else if angle >= 5*math.Pi/4 && angle < 7*math.Pi/4 {
		// Up
		c.direction = LookingUp
	}

	// Handle diagonals
	if angle >= math.Pi/8 && angle < 3*math.Pi/8 {
		// Down-Right
		c.direction = LookingUpRight
	} else if angle >= 5*math.Pi/8 && angle < 7*math.Pi/8 {
		// Down-Left (flip Down-Right)
		c.direction = LookingUpRight
		//c.flipX = true
	} else if angle >= 9*math.Pi/8 && angle < 11*math.Pi/8 {
		// Up-Left (flip Up-Right)
		c.direction = LookingDownRight
		//c.flipX = true
	} else if angle >= 13*math.Pi/8 && angle < 15*math.Pi/8 {
		// Up-Right
		c.direction = LookingDownRight
	}

	// Only update sprite if direction changed
	if oldDirection != c.direction || oldFlipX != c.flipX {
		c.updateSprite()
	}
}

// GetSprite returns the current eye sprite
func (c *EyeController) GetSprite() *ebiten.Image {
	return c.currentSprite
}

// GetPosition returns the eye position relative to character center
func (c *EyeController) GetPosition() common.Vector2 {
	return c.position
}

// GetFlipX returns whether the sprite should be flipped horizontally
func (c *EyeController) GetFlipX() bool {
	return c.flipX
}

// SetPosition sets the eye position relative to character center
func (c *EyeController) SetPosition(pos common.Vector2) {
	c.position = pos
}

// TriggerBlink manually triggers a blink animation
func (c *EyeController) TriggerBlink() {
	if !c.isBlinking {
		c.isBlinking = true

		// Create a new animation for the current eye direction
		var startIdx int
		switch c.direction {
		case LookingCenter:
			startIdx = 0
		case LookingRight:
			if c.flipX {
				startIdx = 6 // Left
			} else {
				startIdx = 3 // Right
			}
		case LookingUp:
			startIdx = 9
		case LookingDown:
			startIdx = 12
		case LookingUpRight:
			if c.flipX {
				startIdx = 21 // UpLeft
			} else {
				startIdx = 15 // UpRight
			}
		case LookingDownRight:
			if c.flipX {
				startIdx = 24 // DownLeft
			} else {
				startIdx = 18 // DownRight
			}
		}

		blinkSequence := []int{0, 1, 2, 1}
		directionFrames := c.allFrames[startIdx : startIdx+3]
		c.animation = sprite.NewAnimationWithSequence(directionFrames, blinkSequence, false)
		c.animation.Reset()
	}
}
