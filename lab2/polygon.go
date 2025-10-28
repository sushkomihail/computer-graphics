package lab2

import (
	base "go-graphics/lab1"
	"image/color"
)

type Vertex struct {
	Point  base.Vector
	Normal base.Vector
}

type Edge struct {
	A *Vertex
	B *Vertex
}

func NewEdge(a, b *Vertex) *Edge {
	return &Edge{
		A: a,
		B: b,
	}
}

type Polygon struct {
	vertices            []base.Vector
	transformedVertices []*Vertex
	edges               []Edge
	center              base.Vector
	planeEquation       PlaneEquation
	globalRotation      base.Vector
	color               color.Color
}

func NewPolygon(vertices []*Vertex, axesScale float32, color color.Color) *Polygon {
	baseVertices := make([]base.Vector, len(vertices))

	for i := range vertices {
		baseVertices[i] = vertices[i].Point
	}

	edges := generateEdges(vertices)

	return &Polygon{
		vertices:            baseVertices,
		transformedVertices: vertices,
		edges:               edges,
		planeEquation:       GetPlaneEquation(vertices[0].Point, vertices[1].Point, vertices[2].Point),
		globalRotation:      base.Vector{},
		color:               color,
	}
}

func (p Polygon) GetVertices() []*Vertex {
	return p.transformedVertices
}

func (p Polygon) GetEdges() []Edge {
	return p.edges
}

func (p Polygon) GetCenter() base.Vector {
	return p.center
}

func (p Polygon) GetPlaneEquation() PlaneEquation {
	return p.planeEquation
}

func CopyVertices(src []base.Vector, dst []*Vertex) {
	for i, v := range src {
		if dst[i] == nil {
			dst[i] = &Vertex{Point: v}
		} else {
			dst[i].Point = v
		}
	}
}

func (p *Polygon) ApplyTransformation(axes base.CoordinatesSystem, deltaRotation base.Vector) {
	CopyVertices(p.vertices, p.transformedVertices)

	p.globalRotation.Add(deltaRotation)

	xRotationMatrix := base.GetRotationXMatrix(p.globalRotation.X)
	yRotationMatrix := base.GetRotationYMatrix(p.globalRotation.Y)

	rotationMatrix := base.MultiplyMatrices(xRotationMatrix, yRotationMatrix)

	for i, v := range p.transformedVertices {
		v.Point.ApplyTransformationMatrix4x4(rotationMatrix)
		axes.ProjectVertex(&v.Point)
		p.transformedVertices[i] = v
	}

	p.planeEquation = GetPlaneEquation(
		p.transformedVertices[0].Point,
		p.transformedVertices[1].Point,
		p.transformedVertices[2].Point)
	p.center = calculateCenter(p.transformedVertices)
}

func (p Polygon) TryGetIntersections(y int) (bool, []Intersection) {
	intersections := make([]Intersection, 0, 4)

	for _, e := range p.edges {
		if !hasIntersection(float32(y), e) {
			continue
		}

		if e.A.Point.Y == e.B.Point.Y {
			intersections = append(intersections, *NewIntersection(e.A.Point, e), *NewIntersection(e.B.Point, e))
			break
		}

		intersection := getIntersection(float32(y), e)
		intersections = append(intersections, intersection)
	}

	return len(intersections) != 0, intersections
}

func generateEdges(vertices []*Vertex) []Edge {
	edges := make([]Edge, len(vertices))
	idx := 0

	for idx < len(vertices)-1 {
		edges[idx] = *NewEdge(vertices[idx], vertices[idx+1])
		idx++
	}

	edges[idx] = *NewEdge(vertices[idx], vertices[0])
	return edges
}

func calculateCenter(vertices []*Vertex) base.Vector {
	center := vertices[0].Point

	for i := 1; i < len(vertices); i++ {
		center.Add(vertices[i].Point)
	}

	center.Multiply(1 / float32(len(vertices)))
	return center
}
