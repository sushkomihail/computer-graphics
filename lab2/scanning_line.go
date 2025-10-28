package lab2

import (
	base "go-graphics/lab1"
)

type Intersection struct {
	point base.Vector
	edge  Edge
}

func NewIntersection(point base.Vector, edge Edge) *Intersection {
	return &Intersection{
		point: point,
		edge:  edge,
	}
}

func (i Intersection) GetPoint() base.Vector {
	return i.point
}

func (i Intersection) GetEdge() Edge {
	return i.edge
}

func hasIntersection(y float32, edge Edge) bool {
	return (edge.A.Point.Y <= y && edge.B.Point.Y >= y) || (edge.A.Point.Y >= y && edge.B.Point.Y <= y)
}

func getIntersection(y float32, edge Edge) Intersection {
	k := (y - edge.A.Point.Y) / (edge.B.Point.Y - edge.A.Point.Y)
	x := edge.A.Point.X + (edge.B.Point.X-edge.A.Point.X)*k
	return *NewIntersection(base.NewVector(x, y, 0), edge)
}
