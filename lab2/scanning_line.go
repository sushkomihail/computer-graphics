package lab2

import base "go-graphics/lab1"

func hasIntersection(y float32, edge Edge) bool {
	return (edge.v1.Y <= y && edge.v2.Y >= y) || (edge.v1.Y >= y && edge.v2.Y <= y)
}

func tryGetIntersectionPoint(y float32, edge Edge) (bool, base.Vector) {
	if hasIntersection(y, edge) {
		k := (y - edge.v1.Y) / (edge.v2.Y - edge.v1.Y)
		x := edge.v1.X + (edge.v2.X-edge.v1.X)*k
		return true, base.NewVector(x, y, 0)
	}

	return false, base.NewVector(0, 0, 0)
}
