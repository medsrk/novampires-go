// internal/engine/camera/debug.go
package camera

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	"novampires-go/internal/common"
	"novampires-go/internal/engine/debug"
	"unsafe"
)

type DebugWindow struct {
	camera  *Camera
	open    bool
	openPtr unsafe.Pointer

	// Values for sliders
	zoom      float32
	rotation  float32
	smoothing float32
	deadzoneX float32
	deadzoneY float32

	// Pointers for sliders
	zoomPtr      unsafe.Pointer
	rotationPtr  unsafe.Pointer
	smoothingPtr unsafe.Pointer
	deadzoneXPtr unsafe.Pointer
	deadzoneYPtr unsafe.Pointer
}

func NewDebugWindow(camera *Camera) *DebugWindow {
	w := &DebugWindow{
		camera:    camera,
		open:      true,
		zoom:      float32(camera.GetZoom()),
		rotation:  float32(camera.GetRotation()),
		smoothing: float32(camera.config.Smoothing),
		deadzoneX: float32(camera.config.Deadzone.Size.X),
		deadzoneY: float32(camera.config.Deadzone.Size.Y),
	}

	w.openPtr = unsafe.Pointer(&w.open)
	w.zoomPtr = unsafe.Pointer(&w.zoom)
	w.rotationPtr = unsafe.Pointer(&w.rotation)
	w.smoothingPtr = unsafe.Pointer(&w.smoothing)
	w.deadzoneXPtr = unsafe.Pointer(&w.deadzoneX)
	w.deadzoneYPtr = unsafe.Pointer(&w.deadzoneY)

	return w
}

func (w *DebugWindow) Draw() {
	if !w.open {
		return
	}

	debug.FixedWindow("Camera Debug", 300, 400, func() {
		// Camera position and target
		debug.CollapsingSection("Position & Target", func() {
			center := w.camera.GetCenter()
			debug.LabeledValue("Center:", fmt.Sprintf("(%.2f, %.2f)", center.X, center.Y), nil)

			target := w.camera.GetTarget()
			if target != nil {
				debug.LabeledValue("Target:", fmt.Sprintf("(%.2f, %.2f)", target.X, target.Y), nil)
			} else {
				debug.LabeledValue("Target:", "None", nil)
			}
		})

		// Transform info
		debug.CollapsingSection("Transform", func() {
			// Zoom control
			if imgui.SliderFloat("Zoom", (*float32)(w.zoomPtr), 0.1, 5.0) {
				w.camera.SetZoom(float64(w.zoom))
			}

			// Rotation control
			if imgui.SliderFloat("Rotation", (*float32)(w.rotationPtr), -3.14, 3.14) {
				w.camera.SetRotation(float64(w.rotation))
			}
		})

		// Viewport info
		debug.CollapsingSection("Viewport", func() {
			viewport := w.camera.GetViewport()
			debug.LabeledValue("Position:", fmt.Sprintf("(%.2f, %.2f)", viewport.Pos.X, viewport.Pos.Y), nil)
			debug.LabeledValue("Size:", fmt.Sprintf("(%.2f, %.2f)", viewport.Size.X, viewport.Size.Y), nil)
		})

		// Camera config
		debug.CollapsingSection("Configuration", func() {
			if imgui.SliderFloat("Smoothing", (*float32)(w.smoothingPtr), 0.01, 1.0) {
				w.camera.config.Smoothing = float64(w.smoothing)
			}

			if imgui.SliderFloat("Deadzone X", (*float32)(w.deadzoneXPtr), 0, 100) {
				w.camera.config.Deadzone.Size.X = float64(w.deadzoneX)
			}

			if imgui.SliderFloat("Deadzone Y", (*float32)(w.deadzoneYPtr), 0, 100) {
				w.camera.config.Deadzone.Size.Y = float64(w.deadzoneY)
			}
		})
	})
}

func (w *DebugWindow) Name() string {
	return common.WindowCameraDebug
}

func (w *DebugWindow) Toggle() {
	w.open = !w.open
}

func (w *DebugWindow) IsOpen() bool {
	return w.open
}

func (w *DebugWindow) Close() {
	w.open = false
}

func (c *Camera) CreateDebugWindow() *DebugWindow {
	return NewDebugWindow(c)
}
