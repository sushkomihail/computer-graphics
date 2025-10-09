package lab2

import (
	base "go-graphics/lab1"
	"image/color"
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

func (p Polygon) tryGetIntersectionPoints(y int) (bool, []base.Vector) {
	points := make([]base.Vector, 0, 4)

	for _, e := range p.edges {
		if !hasIntersection(float32(y), e) {
			continue
		}

		if e.v1.Y == e.v2.Y {
			points = append(points, e.v1, e.v2)
			break
		}

		point := getIntersectionPoint(float32(y), e)
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
