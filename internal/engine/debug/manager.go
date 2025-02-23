package debug

import (
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"novampires-go/internal/common"
)

type Window interface {
	Name() string
	Draw()
	IsOpen() bool
	Toggle()
	Close()
}

type Manager struct {
	enabled bool
	windows map[string]Window
	im      common.InputProvider
}

type Deps struct {
	InputManager common.InputProvider
}

func New(deps Deps) *Manager {
	return &Manager{
		enabled: true,
		windows: make(map[string]Window),
		im:      deps.InputManager,
	}
}

func (m *Manager) Update() {
	if !m.enabled {
		return
	}

	// In debug manager's Update()
	if m.im.JustPressed(common.ActionTogglePlayerDebug) {
		m.GetWindow(common.WindowPlayerDebug).Toggle()
	}
	if m.im.JustPressed(common.ActionToggleBindingEditor) {
		m.GetWindow(common.WindowBindingEdit).Toggle()
	}
	if m.im.JustPressed(common.ActionToggleInputDebug) {
		m.GetWindow(common.WindowInputDebug).Toggle()
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
	m.windows[w.Name()] = w
}

func (m *Manager) RemoveWindow(name string) {
	delete(m.windows, name)
}

func (m *Manager) GetWindow(name string) Window {
	return m.windows[name]
}

func (m *Manager) IsOpen() bool {
	return m.enabled
}

func (m *Manager) Close() {
	m.enabled = false
}

func (m *Manager) Toggle() {
	m.enabled = !m.enabled
}

func (m *Manager) SetDisplaySize(width, height float32) {
	ebimgui.SetDisplaySize(width, height)
}
