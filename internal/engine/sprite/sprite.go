package sprite

import (
	"sync"
	"time"
)

// FrameData represents a single frame in an animation
type FrameData struct {
	// Source rectangle in the spritesheet
	SrcX, SrcY, SrcWidth, SrcHeight int

	// Duration in milliseconds
	Duration int
}

// Animation represents a sequence of frames
type Animation struct {
	// Frames in the animation
	Frames []FrameData

	// Whether the animation should loop
	Loop bool

	// Current state
	currentFrame int
	elapsed      time.Duration
	reversed     bool
	finished     bool

	// Mutex for concurrent access
	mu sync.RWMutex
}

// NewAnimation creates a new animation with the given frames
func NewAnimation(frames []FrameData, loop bool) *Animation {
	// Safety check - don't create animations with no frames
	if len(frames) == 0 {
		return nil
	}

	return &Animation{
		Frames:       frames,
		Loop:         loop,
		currentFrame: 0,
		elapsed:      0,
		reversed:     false,
		finished:     false,
	}
}

// CreateAnimationFromStrip creates a single row animation with safety checks
func CreateAnimationFromStrip(
	frameWidth, frameHeight int,
	startX, startY int,
	frameCount, frameDuration int,
	loop bool,
) *Animation {
	// Safety checks
	if frameCount <= 0 || frameWidth <= 0 || frameHeight <= 0 || frameDuration <= 0 {
		return nil
	}

	frames := make([]FrameData, frameCount)

	for i := 0; i < frameCount; i++ {
		frames[i] = FrameData{
			SrcX:      startX + i*frameWidth,
			SrcY:      startY,
			SrcWidth:  frameWidth,
			SrcHeight: frameHeight,
			Duration:  frameDuration,
		}
	}

	return NewAnimation(frames, loop)
}

func NewAnimationWithSequence(baseFrames []FrameData, sequence []int, loop bool) *Animation {
	// Safety check - don't create animations with no frames or sequence
	if len(baseFrames) == 0 || len(sequence) == 0 {
		return nil
	}

	// Create a new array of frames based on the sequence
	sequencedFrames := make([]FrameData, len(sequence))
	for i, idx := range sequence {
		// Make sure we don't access an index out of bounds
		if idx >= 0 && idx < len(baseFrames) {
			sequencedFrames[i] = baseFrames[idx]
		} else {
			// If index is invalid, use the first frame as a fallback
			sequencedFrames[i] = baseFrames[0]
		}
	}

	// Create a new animation with the sequenced frames
	return NewAnimation(sequencedFrames, loop)
}

// Update advances the animation based on elapsed time
func (a *Animation) Update(dt time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.finished {
		return
	}

	a.elapsed += dt

	frameDuration := time.Duration(a.Frames[a.currentFrame].Duration) * time.Millisecond
	if a.elapsed >= frameDuration {
		a.elapsed -= frameDuration
		if a.reversed {
			a.currentFrame--
			if a.currentFrame < 0 {
				if a.Loop {
					a.currentFrame = len(a.Frames) - 1
				} else {
					a.currentFrame = 0
					a.finished = true
				}
			}
		} else {
			a.currentFrame++
			if a.currentFrame >= len(a.Frames) {
				if a.Loop {
					a.currentFrame = 0
				} else {
					a.currentFrame = len(a.Frames) - 1
					a.finished = true
				}
			}
		}
	}
}

// Reset resets the animation to the first frame
func (a *Animation) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.currentFrame = 0
	a.elapsed = 0
	a.finished = false
	a.reversed = false
}

// IsFinished returns whether the animation has finished
func (a *Animation) IsFinished() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.finished
}

// GetCurrentFrame returns the current frame data
func (a *Animation) GetCurrentFrame() FrameData {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.Frames[a.currentFrame]
}

func (a *Animation) GetCurrentFrameInt() int {
	return a.currentFrame
}

func (a *Animation) Reverse() {
	a.reversed = true
}

func (a *Animation) IsReversed() bool {
	return a.reversed
}

func (a *Animation) SetDuration() {
	
}