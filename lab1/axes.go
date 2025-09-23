package lab1

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CoordinatesSystem struct {
	position Vector
	scale    float32
}

func NewCoordinatesSystem(position Vector, scale float32) *CoordinatesSystem {
	return &CoordinatesSystem{
		position: position,
		scale:    scale,
	}
}

func (cs *CoordinatesSystem) DrawAxes(screen *ebiten.Image, xLength, yLength float32) {
	xColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	yColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	x0, y0 := cs.position.X, cs.position.Y
	x1, y1 := cs.position.X+xLength, cs.position.Y
	vector.StrokeLine(screen, x0, y0, x1, y1, float32(strokeWidth), xColor, false)

	x1, y1 = cs.position.X, -cs.position.Y+yLength
	vector.StrokeLine(screen, x0, y0, x1, y1, float32(strokeWidth), yColor, false)
}

func (cs *CoordinatesSystem) ProjectVertex(v *Vector) {
	mat := GetTranslationMatrix(cs.position)
	v.ApplyTransformationMatrix4x4(mat)
}
