package lab1

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	windowWidth  int     = 640
	windowHeight int     = 480
	axesScale    float32 = 20
	strokeWidth  float32 = 2.5

	focusDistance float32 = 200

	aimationDuration float64 = 0.5
)

type Window struct {
	axes             CoordinatesSystem
	mouse            Mouse
	letter           Object
	projectionMatrix [4][4]float32

	animation      Animation
	fromProjection [4][4]float32
	toProjection   [4][4]float32
}

func NewWindow(axes *CoordinatesSystem, letter *Object) Window {
	return Window{
		axes:             *axes,
		mouse:            Mouse{},
		letter:           *letter,
		projectionMatrix: GetIdentityMatrix(),
		animation:        *NewAnimation(aimationDuration),
		fromProjection:   GetIdentityMatrix(),
		toProjection:     GetPerspectiveProjectionMatrix(focusDistance),
	}
}

func (w *Window) drawHelpText(screen *ebiten.Image) {
	x, y := windowWidth-250, 10
	space := 20

	transform := []string{
		"global position: " + w.letter.position.ToString(),
		"local rotation: " + w.letter.rotation.ToString(),
		"scale: " + w.letter.scale.ToString(),
	}

	for _, s := range transform {
		ebitenutil.DebugPrintAt(screen, s, x, y)
		y += space
	}

	controls := []string{
		"controls:",
		"   translation:",
		"      x-axis: a/d",
		"      z-axis: w/s",
		"      y-axis: arrow up/arrow down",
		"   rotation:",
		"      x-axis: x/mouse dragging",
		"      z-axis: y/mouse dragging",
		"      y-axis: z",
		"   scaling: mouse wheel",
		"   switch projection: p",
	}

	x, y = 10, 10

	for _, s := range controls {
		ebitenutil.DebugPrintAt(screen, s, x, y)
		y += space
	}
}

func (w *Window) Update() error {
	deltaPosition := HandleTranslationInput()
	deltaRotation := HandleRotationInput(&w.mouse)
	scaleFactor := HandleScaleInput()
	HandleAnimationInput(&w.animation)

	w.letter.ApplyTransformation(deltaPosition, deltaRotation, scaleFactor, w.projectionMatrix, w.axes)

	if !w.animation.IsFinished() {
		w.animation.Update(1 / ebiten.ActualFPS())
		w.projectionMatrix = AnimatePerspectiveChange(w.fromProjection, w.toProjection, &w.animation)

		if w.animation.IsFinished() {
			w.fromProjection, w.toProjection = w.toProjection, w.fromProjection
		}
	}

	return nil
}

func (w *Window) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{133, 108, 65, 255})
	w.axes.DrawAxes(screen, float32(windowWidth/2), float32(windowHeight/2))
	w.letter.DrawObject(screen)
	w.drawHelpText(screen)
}

func (w *Window) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func Run() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Lab 1")

	axesPosition := NewVector(float32(windowWidth/2), float32(windowHeight/2), 0)
	axes := NewCoordinatesSystem(axesPosition, axesScale)

	letterVertices := []Vector{
		NewVector(3*axes.scale, -3*axes.scale, -0.5*axes.scale),
		NewVector(0.5*axes.scale, -4*axes.scale, -0.5*axes.scale),
		NewVector(-2*axes.scale, -3*axes.scale, -0.5*axes.scale),
		NewVector(-3*axes.scale, 0, -0.5*axes.scale),
		NewVector(-2*axes.scale, 3*axes.scale, -0.5*axes.scale),
		NewVector(0.5*axes.scale, 4*axes.scale, -0.5*axes.scale),
		NewVector(3*axes.scale, 3*axes.scale, -0.5*axes.scale),

		NewVector(3*axes.scale, 3*axes.scale, 0.5*axes.scale),
		NewVector(0.5*axes.scale, 4*axes.scale, 0.5*axes.scale),
		NewVector(-2*axes.scale, 3*axes.scale, 0.5*axes.scale),
		NewVector(-3*axes.scale, 0, 0.5*axes.scale),
		NewVector(-2*axes.scale, -3*axes.scale, 0.5*axes.scale),
		NewVector(0.5*axes.scale, -4*axes.scale, 0.5*axes.scale),
		NewVector(3*axes.scale, -3*axes.scale, 0.5*axes.scale),

		NewVector(0, 0, 0),
		NewVector(axes.scale, 0, 0),
		NewVector(0, -axes.scale, 0),
		NewVector(0, 0, axes.scale),
	}

	window := NewWindow(axes, NewObject(letterVertices))

	if err := ebiten.RunGame(&window); err != nil {
		log.Fatal(err)
	}
}
