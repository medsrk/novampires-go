package debug

import (
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

type Window interface {
	Draw()
}

type Manager struct {
	enabled bool
	windows []Window
}

func New() *Manager {
	return &Manager{
		enabled: true,
		windows: make([]Window, 0),
	}
}

func (m *Manager) Update() {
	if !m.enabled {
		return
	}

	ebimgui.Update(1.0 / 60.0) // Fixed update rate for ImGui
}

func (m *Manager) BeginFrame() {
	if !m.enabled {
		return
	}

	ebimgui.BeginFrame()
}

func (m *Manager) EndFrame() {
	if !m.enabled {
		return
	}

	ebimgui.EndFrame()
}

func (m *Manager) Draw(screen *ebiten.Image) {
	if !m.enabled {
		return
	}

	m.BeginFrame()

	for _, window := range m.windows {
		window.Draw()
	}

	m.EndFrame()
	ebimgui.Draw(screen)
}

func (m *Manager) AddWindow(w Window) {
	m.windows = append(m.windows, w)
}

func (m *Manager) Toggle() {
	m.enabled = !m.enabled
}

func (m *Manager) SetDisplaySize(width, height float32) {
	ebimgui.SetDisplaySize(width, height)
}
