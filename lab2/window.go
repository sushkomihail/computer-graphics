package lab2

import (
	base "go-graphics/lab1"
	"image/color"
	"log"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	windowWidth  = 640
	windowHeight = 480
)

type Window struct {
	mouse       base.Mouse
	axes        base.CoordinatesSystem
	polygons    []Polygon
	frameBuffer []int
	zBuffer     []float32
}

func NewWindow(axes base.CoordinatesSystem, polygons []Polygon) *Window {
	return &Window{
		mouse:       base.Mouse{},
		axes:        axes,
		polygons:    polygons,
		frameBuffer: make([]int, windowWidth),
		zBuffer:     make([]float32, windowWidth),
	}
}

func (w *Window) clearBuffers() {
	for i := range windowWidth {
		w.frameBuffer[i] = 0
		w.zBuffer[i] = 0
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

	for y := range windowHeight {
		w.clearBuffers()

		for i, p := range w.polygons {
			ok, intersections := p.tryGetIntersectionPoints(y)

			if !ok {
				continue
			}

			sort.Slice(intersections, func(i, j int) bool {
				return intersections[i].X < intersections[j].X
			})

			for x := int(intersections[0].X); x <= int(intersections[len(intersections)-1].X); x++ {
				if x < 0 && x > windowWidth-1 {
					continue
				}

				z := p.planeEquation.GetPlaneZ(float32(x), float32(y))

				if w.frameBuffer[x] == 0 {
					w.frameBuffer[x] = i + 1
					w.zBuffer[x] = z
				} else if z > w.zBuffer[x] {
					w.zBuffer[x] = z
					w.frameBuffer[x] = i + 1
				}
			}
		}

		for x := range w.frameBuffer {
			bufVal := w.frameBuffer[x]

			if bufVal != 0 {
				screen.Set(x, y, w.polygons[bufVal-1].color)
			}
		}
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
