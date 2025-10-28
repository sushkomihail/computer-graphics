package lab3

import (
	"fmt"
	base "go-graphics/lab1"
	graphics "go-graphics/lab2"
	"image/color"
	"log"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	windowWidth  = 640
	windowHeight = 480

	animationDuration = 5
)

type Window struct {
	mouse base.Mouse

	camera base.CoordinatesSystem

	object RotationObject

	frameBuffer []float64
	zBuffer     []float32

	animation          base.Animation
	cameraFromPosition base.Vector
	cameraToPosition   base.Vector

	isSimpleLightingModel bool
}

func NewWindow(axes base.CoordinatesSystem, object RotationObject, cameraFromPosition, cameraToPosition base.Vector) *Window {
	return &Window{
		mouse:                 base.Mouse{},
		camera:                axes,
		object:                object,
		frameBuffer:           make([]float64, windowWidth),
		zBuffer:               make([]float32, windowWidth),
		animation:             *base.NewAnimation(animationDuration),
		cameraFromPosition:    cameraFromPosition,
		cameraToPosition:      cameraToPosition,
		isSimpleLightingModel: true,
	}
}

func (w *Window) clearBuffers() {
	for i := range windowWidth {
		w.frameBuffer[i] = -1
		w.zBuffer[i] = 0
	}
}

func (w *Window) fillBuffers(y int) {
	distanceToLightSource := calculateMinDistanceToLightSource(w.object.vertices)

	for _, p := range w.object.polygons {
		ok, intersections := p.TryGetIntersections(y)

		if !ok {
			continue
		}

		sort.Slice(intersections, func(i, j int) bool {
			return intersections[i].GetPoint().X < intersections[j].GetPoint().X
		})

		leftIntersection := intersections[0]
		rightIntersection := intersections[len(intersections)-1]

		xStart := int(math.Max(0, float64(leftIntersection.GetPoint().X)))
		xEnd := int(math.Min(windowWidth-1, float64(rightIntersection.GetPoint().X)))

		var intensity float64
		var leftIntensity float64
		var rightIntensity float64

		if w.isSimpleLightingModel {
			lightDirection := lightSourcePosition
			lightDirection.Add(p.GetCenter().Reverse())

			normal := base.NewVector(
				p.GetPlaneEquation().A,
				p.GetPlaneEquation().B,
				p.GetPlaneEquation().C,
			)

			observeDirection := observerPosition
			observeDirection.Add(p.GetCenter().Reverse())

			intensity = calculateIntensity(distanceToLightSource, lightDirection, normal, observeDirection)
		} else {
			leftIntersectionPoint := leftIntersection.GetPoint()
			leftIntersectionPoint.Z = p.GetPlaneEquation().GetPlaneZ(leftIntersectionPoint.X, leftIntersectionPoint.Y)

			rightIntersectionPoint := rightIntersection.GetPoint()
			rightIntersectionPoint.Z = p.GetPlaneEquation().GetPlaneZ(rightIntersectionPoint.X, rightIntersectionPoint.Y)

			leftIntensity = calculateIntersectionIntensity(leftIntersection, leftIntersectionPoint, distanceToLightSource)
			rightIntensity = calculateIntersectionIntensity(rightIntersection, rightIntersectionPoint, distanceToLightSource)
		}

		for x := xStart; x <= xEnd; x++ {
			if x < 0 && x > windowWidth-1 {
				continue
			}

			z := p.GetPlaneEquation().GetPlaneZ(float32(x), float32(y))

			if !w.isSimpleLightingModel {
				var k float64

				if xStart != xEnd {
					k = float64(x-xStart) / float64(xEnd-xStart)
				}

				intensity = leftIntensity + (rightIntensity-leftIntensity)*k
			}

			if (w.frameBuffer)[x] == -1 {
				(w.frameBuffer)[x] = intensity
				(w.zBuffer)[x] = z
			} else if z > (w.zBuffer)[x] {
				(w.zBuffer)[x] = z
				(w.frameBuffer)[x] = intensity
			}
		}
	}
}

func calculateMinDistanceToLightSource(vertices []*graphics.Vertex) float64 {
	minDistance := base.CalculateDistance(vertices[0].Point, lightSourcePosition)

	for i := 1; i < len(vertices); i++ {
		distance := base.CalculateDistance(vertices[i].Point, lightSourcePosition)

		if distance < minDistance {
			minDistance = distance
		}
	}

	return minDistance
}

func (w *Window) drawHelpText(screen *ebiten.Image) {
	x, y := windowWidth-160, 10
	space := 20

	lighting := []string{
		fmt.Sprintf("intensity: %.1f", lightIntensity),
		fmt.Sprintf("specular reflection: %.1f", specReflCoef),
		fmt.Sprintf("diffuse reflection: %.1f", diffReflCoef),
	}

	for _, s := range lighting {
		ebitenutil.DebugPrintAt(screen, s, x, y)
		y += space
	}

	controls := []string{
		"controls:",
		"   rotation:",
		"      y-axis: x/mouse dragging",
		"      x-axis: y/mouse dragging",
		"   lighting:",
		"      intensity: i/shift + i",
		"      specular reflection: s/shift + s",
		"      diffuse reflection: d/shift + d",
		"   camera animation: p",
	}

	x, y = 10, 10

	for _, s := range controls {
		ebitenutil.DebugPrintAt(screen, s, x, y)
		y += space
	}
}

func (w *Window) Update() error {
	deltaRotation := base.HandleRotationInput(&w.mouse)
	base.HandleAnimationInput(&w.animation)
	base.HandleLightingInput(&w.isSimpleLightingModel, &lightIntensity, &specReflCoef, &diffReflCoef)

	w.object.ApplyTransformation(w.camera, deltaRotation)

	if !w.animation.IsFinished() {
		w.animation.Update(1 / ebiten.ActualFPS())
		w.camera.SetPosition(animateCameraMoving(w.cameraFromPosition, w.cameraToPosition, &w.animation))

		if w.animation.IsFinished() {
			w.cameraFromPosition, w.cameraToPosition = w.cameraToPosition, w.cameraFromPosition
		}
	}

	return nil
}

func (w *Window) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{133, 108, 65, 255})
	w.drawHelpText(screen)

	for y := range windowHeight {
		w.clearBuffers()
		w.fillBuffers(y)

		for x := range w.frameBuffer {
			intensity := w.frameBuffer[x]

			if intensity != -1 {
				channel := uint8(255 * intensity)
				color := color.RGBA{channel, channel, channel, 255}
				screen.Set(x, y, color)
			}
		}
	}
}

func (w *Window) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func Run() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Lab 3")

	padding := 120
	cameraFromPosition := base.NewVector(float32(padding), windowHeight/2+70, 0)
	cameraToPosition := base.NewVector(float32(windowWidth-padding), windowHeight/2+70, 0)
	camera := *base.NewCoordinatesSystem(cameraFromPosition, 20)

	generatingCurve := getSemicircle(20, 5)
	object := *NewRotationObject(generatingCurve, 20, camera.GetScale())

	window := NewWindow(camera, object, cameraFromPosition, cameraToPosition)

	if err := ebiten.RunGame(window); err != nil {
		log.Fatal(err)
	}
}
