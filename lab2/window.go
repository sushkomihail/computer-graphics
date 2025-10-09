package lab2

import (
	base "go-graphics/lab1"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	windowWidth  = 640
	windowHeight = 480
)

type Window struct {
	mouse        base.Mouse
	axes         base.CoordinatesSystem
	polygons     []Polygon
	screenBuffer []int
	zBuffer      []int
}

func NewWindow(axes base.CoordinatesSystem, polygons []Polygon) *Window {
	return &Window{
		mouse:        base.Mouse{},
		axes:         axes,
		polygons:     polygons,
		screenBuffer: make([]int, windowWidth*windowHeight),
		zBuffer:      make([]int, windowWidth*windowHeight),
	}
}

func (w *Window) drawHelpText(screen *ebiten.Image) {
	x, y := windowWidth-250, 10
	space := 20

	controls := []string{
		"controls:",
		"   rotation:",
		"      x-axis: x/mouse dragging",
		"      z-axis: y/mouse dragging",
	}

	x, y = 10, 10

	for _, s := range controls {
		ebitenutil.DebugPrintAt(screen, s, x, y)
		y += space
	}
}

func (w *Window) Update() error {
	deltaRotation := base.HandleRotationInput(&w.mouse)

	for i := range w.polygons {
		w.polygons[i].ApplyTransformation(w.axes, deltaRotation)
	}

	return nil
}

func (w *Window) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{133, 108, 65, 255})
	w.drawHelpText(screen)

	for i := range w.zBuffer {
		w.zBuffer[i] = int(math.Inf(-1))
	}

	for _, v := range w.polygons {
		v.Fill(screen, windowWidth, windowHeight, w.zBuffer)
	}
}

func (w *Window) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func Run() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Lab 2")

	axes := *base.NewCoordinatesSystem(base.NewVector(windowWidth/2, windowHeight/2, 0), 20)

	triangle := NewPolygon([]base.Vector{
		base.NewVector(0, 1, 2),
		base.NewVector(-1, -1, -3),
		base.NewVector(1, -1, -3),
	}, axes.GetScale(), color.RGBA{255, 255, 255, 255})

	square := NewPolygon([]base.Vector{
		base.NewVector(-1, -1, 0.5),
		base.NewVector(1, -1, 0.5),
		base.NewVector(1, 1, 0.5),
		base.NewVector(-1, 1, 0.5),
	}, axes.GetScale(), color.RGBA{150, 150, 150, 255})

	pentagon := NewPolygon([]base.Vector{
		base.NewVector(0, -3, -1.5),
		base.NewVector(-3, -1, -1.5),
		base.NewVector(-2, 3, -1.5),
		base.NewVector(2, 3, -1.5),
		base.NewVector(3, -1, -1.5),
	}, axes.GetScale(), color.RGBA{99, 99, 99, 255})

	hexagon := NewPolygon([]base.Vector{
		base.NewVector(6, 0, -4),
		base.NewVector(3, -5, -4),
		base.NewVector(-3, -5, -4),
		base.NewVector(-6, 0, -4),
		base.NewVector(-3, 5, -4),
		base.NewVector(3, 5, -4),
	}, axes.GetScale(), color.RGBA{66, 66, 66, 255})

	window := NewWindow(axes, []Polygon{
		*triangle,
		*square,
		*pentagon,
		*hexagon,
	})

	if err := ebiten.RunGame(window); err != nil {
		log.Fatal(err)
	}
}
