package input

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestInputMap(t *testing.T) {
	t.Run("new input map has default bindings", func(t *testing.T) {
		im := NewInputMap()

		action, exists := im.GetAction(ebiten.KeyW)
		assert.True(t, exists)
		assert.Equal(t, MoveUp, action)
	})

	t.Run("can bind and unbind keys", func(t *testing.T) {
		im := NewInputMap()

		// Bind a new key
		im.Bind(ebiten.KeyT, Attack)
		action, exists := im.GetAction(ebiten.KeyT)
		assert.True(t, exists)
		assert.Equal(t, Attack, action)

		// Unbind it
		im.UnBind(ebiten.KeyT)
		_, exists = im.GetAction(ebiten.KeyT)
		assert.False(t, exists)
	})
}

// mockKeyStateProvider simulates ebiten.IsKeyPressed for testing
type mockKeyStateProvider struct {
	keyStates map[ebiten.Key]bool
}

func newMockKeyStateProvider() *mockKeyStateProvider {
	return &mockKeyStateProvider{
		keyStates: make(map[ebiten.Key]bool),
	}
}

func (m *mockKeyStateProvider) IsKeyPressed(key ebiten.Key) bool {
	pressed, exists := m.keyStates[key]
	return exists && pressed
}

func (m *mockKeyStateProvider) setKeyPressed(key ebiten.Key, pressed bool) {
	m.keyStates[key] = pressed
}

func TestInputManager(t *testing.T) {
	t.Run("movement vector is normalized for diagonal movement", func(t *testing.T) {
		mock := newMockKeyStateProvider()
		m := NewManagerWithProvider(mock)

		// Mock diagonal movement (up+right)
		mock.setKeyPressed(ebiten.KeyW, true)
		mock.setKeyPressed(ebiten.KeyD, true)

		// Need one update to establish state
		m.Update()

		dx, dy := m.GetMovementVector()

		// Vector should be normalized (length = 1)
		length := math.Sqrt(dx*dx + dy*dy)
		assert.InDelta(t, 1.0, length, 0.0001)

		// Verify actual values
		expectedDx := 0.707107
		expectedDy := -0.707107
		assert.InDelta(t, expectedDx, dx, 0.0001)
		assert.InDelta(t, expectedDy, dy, 0.0001)
	})

	t.Run("detects just pressed state", func(t *testing.T) {
		mock := newMockKeyStateProvider()
		m := NewManagerWithProvider(mock)

		// First frame: key not pressed
		mock.setKeyPressed(ebiten.KeyW, false)
		m.Update()
		assert.False(t, m.JustPressed(MoveUp), "Should not be just pressed initially")

		// Second frame: key pressed
		mock.setKeyPressed(ebiten.KeyW, true)
		m.Update()
		assert.True(t, m.JustPressed(MoveUp), "Should be just pressed when key first pressed")

		// Third frame: key still pressed
		m.Update()
		assert.False(t, m.JustPressed(MoveUp), "Should not be just pressed when key held")
	})

	t.Run("detects just released state", func(t *testing.T) {
		mock := newMockKeyStateProvider()
		m := NewManagerWithProvider(mock)

		// First frame: key pressed
		mock.setKeyPressed(ebiten.KeyW, true)
		m.Update()
		assert.False(t, m.JustReleased(MoveUp), "Should not be just released initially")

		// Second frame: key released
		mock.setKeyPressed(ebiten.KeyW, false)
		m.Update()
		assert.True(t, m.JustReleased(MoveUp), "Should be just released when key first released")

		// Third frame: key still released
		m.Update()
		assert.False(t, m.JustReleased(MoveUp), "Should not be just released when key remains released")
	})

	t.Run("rebinding keys works", func(t *testing.T) {
		mock := newMockKeyStateProvider()
		m := NewManagerWithProvider(mock)

		// First verify W works for MoveUp
		mock.setKeyPressed(ebiten.KeyW, true)
		m.Update()
		assert.True(t, m.IsActive(MoveUp), "W should trigger MoveUp initially")

		// Rebind MoveUp from W to T
		m.RebindKey(ebiten.KeyW, ebiten.KeyT)

		// W should no longer work
		m.Update()
		assert.False(t, m.IsActive(MoveUp), "W should no longer trigger MoveUp")

		// But T should
		mock.setKeyPressed(ebiten.KeyW, false)
		mock.setKeyPressed(ebiten.KeyT, true)
		m.Update()
		assert.True(t, m.IsActive(MoveUp), "T should now trigger MoveUp")
	})
}
