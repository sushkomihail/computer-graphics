package lab2

import (
	base "go-graphics/lab1"
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Edge struct {
	v1 base.Vector
	v2 base.Vector
}

func NewEdge(v1, v2 base.Vector) *Edge {
	return &Edge{
		v1: v1,
		v2: v2,
	}
}

type Polygon struct {
	vertices       []base.Vector
	edges          []Edge
	planeEquation  PlaneEquation
	globalRotation base.Vector
	color          color.Color
}

func NewPolygon(vertices []base.Vector, axesScale float32, color color.Color) *Polygon {
	for i := range vertices {
		vertices[i].Multiply(axesScale)
	}

	return &Polygon{
		vertices:       vertices,
		edges:          generateEdges(vertices),
		planeEquation:  GetPlaneEquation(vertices[0], vertices[1], vertices[2]),
		globalRotation: base.Vector{},
		color:          color,
	}
}

func (p Polygon) GetEdges() []Edge {
	return p.edges
}

func (p *Polygon) ApplyTransformation(axes base.CoordinatesSystem, deltaRotation base.Vector) {
	verticesCopy := make([]base.Vector, len(p.vertices))
	copy(verticesCopy, p.vertices)

	p.globalRotation.Add(deltaRotation)

	xRotationMatrix := base.GetRotationXMatrix(p.globalRotation.X)
	yRotationMatrix := base.GetRotationYMatrix(p.globalRotation.Y)

	rotationMatrix := base.MultiplyMatrices(xRotationMatrix, yRotationMatrix)

	for i, v := range verticesCopy {
		v.ApplyTransformationMatrix4x4(rotationMatrix)
		axes.ProjectVertex(&v)
		verticesCopy[i] = v
	}

	p.edges = generateEdges(verticesCopy)
	p.planeEquation = GetPlaneEquation(verticesCopy[0], verticesCopy[1], verticesCopy[2])
}

func (p Polygon) Fill(screen *ebiten.Image, width, height int, zBuffer []int) {
	for y := range height {
		ok, points := p.tryGetIntersectionPoints(y + 1)

		if !ok {
			continue
		}

		sort.Slice(points, func(i, j int) bool {
			return points[i].X < points[j].X
		})

		for x := int(points[0].X); x <= int(points[len(points)-1].X); x++ {
			if x < 0 && x > width {
				continue
			}

			z := int(math.Round(float64(p.planeEquation.GetPlaneZ(float32(x), float32(y)))))
			bufferOffset := y*width + x

			if z > zBuffer[bufferOffset] {
				screen.Set(x, y, p.color)
				zBuffer[bufferOffset] = z
			}
		}
	}
}

func (p Polygon) tryGetIntersectionPoints(line int) (bool, []base.Vector) {
	points := make([]base.Vector, 0, 4)

	for _, e := range p.edges {
		if hasIntersection(float32(line), e) && e.v1.Y == e.v2.Y {
			points = append(points, e.v1, e.v2)
			break
		}

		ok, point := tryGetIntersectionPoint(float32(line), e)

		if !ok {
			continue
		}

		points = append(points, point)
	}

	return len(points) != 0, points
}

func generateEdges(vertices []base.Vector) []Edge {
	edges := make([]Edge, len(vertices))
	idx := 0

	for idx < len(vertices)-1 {
		edges[idx] = *NewEdge(vertices[idx], vertices[idx+1])
		idx++
	}

	edges[idx] = *NewEdge(vertices[idx], vertices[0])
	return edges
}
