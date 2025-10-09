package lab2

import base "go-graphics/lab1"

func hasIntersection(y float32, edge Edge) bool {
	return (edge.v1.Y <= y && edge.v2.Y >= y) || (edge.v1.Y >= y && edge.v2.Y <= y)
}

func getIntersectionPoint(y float32, edge Edge) base.Vector {
	k := (y - edge.v1.Y) / (edge.v2.Y - edge.v1.Y)
	x := edge.v1.X + (edge.v2.X-edge.v1.X)*k
	return base.NewVector(x, y, 0)
}
